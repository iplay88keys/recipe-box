package main

import (
    "database/sql"
    "flag"
    "fmt"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"

    "github.com/gorilla/mux"

    "github.com/iplay88keys/recipe-box/pkg/api"
    "github.com/iplay88keys/recipe-box/pkg/api/recipes"
    "github.com/iplay88keys/recipe-box/pkg/api/users"
    "github.com/iplay88keys/recipe-box/pkg/repositories"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    var static string
    var port string
    var databaseURL string

    flag.StringVar(&static, "static", "ui/build", "the directory to serve static files from")
    flag.StringVar(&port, "port", "8080", "the port to listen on")
    flag.StringVar(&databaseURL, "databaseURL", "", "the url for the database formatted as: mysql://$USER:$PASSWORD@tcp($HOST:$PORT)/$DATABASE")
    flag.Parse()

    envPort := os.Getenv("PORT")
    if envPort != "" {
        port = envPort
    }

    envDatabaseURL := os.Getenv("DATABASE_URL")
    if envDatabaseURL != "" {
        databaseURL = envDatabaseURL
    }

    unquotedURL, err := strconv.Unquote(databaseURL)
    if err == nil {
        databaseURL = unquotedURL
    }

    db, err := sql.Open("mysql", strings.TrimSpace(strings.TrimPrefix(databaseURL, "mysql://")))
    if err != nil {
        panic(err)
    }

    recipesRepo := repositories.NewRecipesRepository(db)
    ingredientsRepo := repositories.NewIngredientsRepository(db)
    stepsRepo := repositories.NewStepsRepository(db)
    usersRepo := repositories.NewUsersRepository(db)

    mux.NewRouter()
    a := api.New(&api.Config{
        Port:      port,
        StaticDir: "ui/build",
        Endpoints: []api.Endpoint{
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
        },
    })

    fmt.Printf("Serving at http://localhost:%s\n", port)
    fmt.Println("ctrl-c to quit")
    stopApi := a.Start()

    defer stopApi()

    blockUntilSigterm()
}

func blockUntilSigterm() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

    <-sigs
}
