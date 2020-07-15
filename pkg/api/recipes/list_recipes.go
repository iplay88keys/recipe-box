package recipes

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/iplay88keys/my-recipe-library/pkg/api/auth"

    "github.com/iplay88keys/my-recipe-library/pkg/api"
    "github.com/iplay88keys/my-recipe-library/pkg/repositories"
)

type RecipeListResponse struct {
    Recipes []*repositories.Recipe `json:"recipes"`
}

type listRecipes func(userID int64) ([]*repositories.Recipe, error)

func ListRecipes(listRecipes listRecipes) api.Endpoint {
    return api.Endpoint{
        Path:   "recipes",
        Method: http.MethodGet,
        Auth:   true,
        Handler: func(w http.ResponseWriter, r *http.Request) {
            id := r.Context().Value(auth.ContextUserKey)
            userID, ok := id.(int64)
            if !ok {
                fmt.Printf("Failed to cast userID to int64 for list recipes endpoint: '%s'\n", id)
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            recipes, err := listRecipes(userID)
            if err != nil {
                if err == sql.ErrNoRows {
                    w.WriteHeader(http.StatusNoContent)
                    return
                }

                fmt.Printf("FormError listing recipes: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            recipeBytes, err := json.Marshal(&RecipeListResponse{
                Recipes: recipes,
            })
            if err != nil {
                fmt.Printf("FormError marshaling recipe list: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            api.LogWriteErr(w.Write(recipeBytes))
        },
    }
}
