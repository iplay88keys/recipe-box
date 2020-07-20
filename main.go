package main

import (
    "database/sql"
    "fmt"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
    "time"

    "github.com/go-redis/redis"

    "github.com/iplay88keys/my-recipe-library/pkg/api"
    "github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
    "github.com/iplay88keys/my-recipe-library/pkg/api/users"
    "github.com/iplay88keys/my-recipe-library/pkg/repositories"
    "github.com/iplay88keys/my-recipe-library/pkg/token"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    var (
        databaseURL   string
        redisURL      string
        accessSecret  string
        refreshSecret string
    )
    static := "ui/build"
    port := "8080"

    envPort := os.Getenv("PORT")
    if envPort != "" {
        port = envPort
    }

    envDatabaseURL := os.Getenv("DATABASE_URL")
    if envDatabaseURL == "" {
        panic("DATABASE_URL is required, formatted as: mysql://$USER:$PASSWORD@tcp($HOST:$PORT)/$DATABASE")
    }
    databaseURL = envDatabaseURL

    envRedisURL := os.Getenv("REDIS_URL")
    if envRedisURL == "" {
        panic("REDIS_URL is required")
    }
    redisURL = envRedisURL

    envAccessSecret := os.Getenv("ACCESS_SECRET")
    if envAccessSecret == "" {
        panic("ACCESS_SECRET is required")
    }
    accessSecret = envAccessSecret

    envRefreshSecret := os.Getenv("REFRESH_SECRET")
    if envRefreshSecret == "" {
        panic("REFRESH_SECRET is required")
    }
    refreshSecret = envRefreshSecret

    db, err := connectToMySQL(databaseURL)
    if err != nil {
        panic(err)
    }

    redisClient, err := connectToRedis(redisURL)
    if err != nil {
        panic(err)
    }

    recipesRepo := repositories.NewRecipesRepository(db)
    ingredientsRepo := repositories.NewIngredientsRepository(db)
    stepsRepo := repositories.NewStepsRepository(db)
    usersRepo := repositories.NewUsersRepository(db)

    redisRepo := repositories.NewRedisRepository(redisClient)
    tokenService := token.NewService(accessSecret, refreshSecret)

    a := api.New(&api.Config{
        Port:                  port,
        StaticDir:             static,
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

    fmt.Printf("Serving at http://localhost:%s\n", port)
    fmt.Println("ctrl-c to quit")
    stopApi := a.Start()

    defer stopApi()
    defer disconnectFromMySQL(db)
    defer disconnectFromRedis(redisClient)

    blockUntilSigterm()
}

func connectToMySQL(databaseURL string) (*sql.DB, error) {
    unquotedURL, err := strconv.Unquote(databaseURL)
    if err == nil {
        databaseURL = unquotedURL
    }

    db, err := sql.Open("mysql", strings.TrimSpace(strings.TrimPrefix(databaseURL, "mysql://")))
    if err != nil {
        return nil, err
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
