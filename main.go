package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
    "time"

    "code.cloudfoundry.org/go-envstruct"
    "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
    "github.com/go-redis/redis"

    "github.com/iplay88keys/my-recipe-library/pkg/api"
    "github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
    "github.com/iplay88keys/my-recipe-library/pkg/api/users"
    "github.com/iplay88keys/my-recipe-library/pkg/repositories"
    "github.com/iplay88keys/my-recipe-library/pkg/token"

    _ "github.com/go-sql-driver/mysql"
)

type Config struct {
    DatabaseCreds DBCreds `env:"DATABASE_CREDS, required"`
    RedisURL      string  `env:"REDIS_URL,      required"`
    AccessSecret  string  `env:"ACCESS_SECRET,  required"`
    RefreshSecret string  `env:"REFRESH_SECRET, required"`
    Port          string  `env:"PORT"`
    Static        string  `env:"STATIC_DIR"`
}

type DBCreds struct {
    URL          string `json:"url"`
    InstanceName string `json:"gcloud_instance_name"`
    DBName       string `json:"gcloud_db_name"`
    User         string `json:"gcloud_user"`
    Password     string `json:"gcloud_password"`
}

func (d *DBCreds) UnmarshalEnv(data string) error {
    return json.Unmarshal([]byte(data), d)
}

func main() {
    config := Config{
        Port:   "8080",
        Static: "ui/build",
    }

    err := envstruct.Load(&config)
    if err != nil {
        panic(err)
    }

    db, err := connectToMySQL(&config)
    if err != nil {
        panic(err)
    }

    redisClient, err := connectToRedis(config.RedisURL)
    if err != nil {
        panic(err)
    }

    recipesRepo := repositories.NewRecipesRepository(db)
    ingredientsRepo := repositories.NewIngredientsRepository(db)
    stepsRepo := repositories.NewStepsRepository(db)
    usersRepo := repositories.NewUsersRepository(db)

    redisRepo := repositories.NewRedisRepository(redisClient)
    tokenService := token.NewService(config.AccessSecret, config.RefreshSecret)

    a := api.New(&api.Config{
        Port:                  config.Port,
        StaticDir:             config.Static,
        Validate:              tokenService.ValidateToken,
        RetrieveAccessDetails: redisRepo.RetrieveTokenDetails,
        Endpoints: []*api.Endpoint{
            recipes.CreateRecipe(recipesRepo.Insert),
            recipes.ListRecipes(recipesRepo.List),
            recipes.GetRecipe(
                recipesRepo.Get,
                ingredientsRepo.GetForRecipe,
                stepsRepo.GetForRecipe,
            ),
            users.Register(
                usersRepo.ExistsByUsername,
                usersRepo.ExistsByEmail,
                usersRepo.Insert,
            ),
            users.Login(
                usersRepo.Verify,
                tokenService.CreateToken,
                redisRepo.StoreTokenDetails,
            ),
            users.Logout(
                tokenService.ValidateToken,
                redisRepo.DeleteTokenDetails,
            ),
        },
    })

    fmt.Printf("Serving at http://localhost:%s\n", config.Port)
    fmt.Println("ctrl-c to quit")
    stopApi := a.Start()

    defer stopApi()
    defer disconnectFromMySQL(db)
    defer disconnectFromRedis(redisClient)

    blockUntilSigterm()
}

func connectToMySQL(config *Config) (db *sql.DB, err error) {
    if config.DatabaseCreds.URL != "" {
        var unquotedURL string
        url := config.DatabaseCreds.URL

        unquotedURL, err = strconv.Unquote(url)
        if err == nil {
            url = unquotedURL
        }

        db, err = sql.Open("mysql", strings.TrimSpace(strings.TrimPrefix(url, "mysql://")))
        if err != nil {
            return nil, err
        }

    } else {
        cfg := mysql.Cfg(config.DatabaseCreds.InstanceName, config.DatabaseCreds.User, config.DatabaseCreds.Password)
        cfg.DBName = config.DatabaseCreds.DBName

        db, err = mysql.DialCfg(cfg)
        if err != nil {
            return nil, err
        }
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}

func connectToRedis(redisURL string) (redis.Cmdable, error) {
    options, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, err
    }

    redisClient := redis.NewClient(options)

    _, err = redisClient.Ping().Result()
    if err != nil {
        return nil, err
    }

    return redisClient, nil
}

func disconnectFromMySQL(db *sql.DB) {
    var stats sql.DBStats
    stats = db.Stats()

    var maxCount int
    for stats.InUse != 0 {
        if maxCount == 10 {
            break
        }

        stats = db.Stats()

        fmt.Printf("Waiting on open mySQL connections: %d in use\n", stats.InUse)

        maxCount += 1
        time.Sleep(100 * time.Millisecond)
    }

    err := db.Close()
    if err != nil {
        panic(err)
    }
}

func disconnectFromRedis(client redis.Cmdable) {
    err := client.(*redis.Client).Close()
    if err != nil {
        panic(err)
    }
}

func blockUntilSigterm() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

    <-sigs
}
