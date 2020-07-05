package recipes

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/iplay88keys/recipe-box/pkg/api"
    "github.com/iplay88keys/recipe-box/pkg/repositories"
)

type RecipeListResponse struct {
    Recipes []*repositories.Recipe `json:"recipes"`
}

type listRecipes func() ([]*repositories.Recipe, error)

func ListRecipes(listRecipes listRecipes) api.Endpoint {
    return api.Endpoint{
        Path:   "recipes",
        Method: http.MethodGet,
        Handler: func(w http.ResponseWriter, r *http.Request) {
            recipes, err := listRecipes()
            if err != nil {
                fmt.Printf("Error listing recipes: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            recipeBytes, err := json.Marshal(&RecipeListResponse{
                Recipes: recipes,
            })
            if err != nil {
                fmt.Printf("Error marshaling recipe list: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            api.LogWriteErr(w.Write(recipeBytes))
        },
    }
}
