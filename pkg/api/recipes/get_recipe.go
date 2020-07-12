package recipes

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sort"
    "strconv"

    "github.com/gorilla/mux"

    "github.com/iplay88keys/recipe-box/pkg/api"
    "github.com/iplay88keys/recipe-box/pkg/repositories"
)

type RecipeResponse struct {
    repositories.Recipe
    Ingredients []*repositories.Ingredient `json:"ingredients"`
    Steps       []*repositories.Step       `json:"steps"`
}

type ByIngredientNumber []*repositories.Ingredient

func (a ByIngredientNumber) Len() int           { return len(a) }
func (a ByIngredientNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIngredientNumber) Less(i, j int) bool { return *a[i].IngredientNumber < *a[j].IngredientNumber }

type ByStepNumber []*repositories.Step

func (a ByStepNumber) Len() int           { return len(a) }
func (a ByStepNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStepNumber) Less(i, j int) bool { return *a[i].StepNumber < *a[j].StepNumber }

type getRecipe func(recipeID int) (*repositories.Recipe, error)
type getIngredientsForRecipe func(recipeID int) ([]*repositories.Ingredient, error)
type getStepsForRecipe func(recipeID int) ([]*repositories.Step, error)

func GetRecipe(getRecipe getRecipe, getIngredientsForRecipe getIngredientsForRecipe, getStepsForRecipe getStepsForRecipe) api.Endpoint {
    return api.Endpoint{
        Path:   "recipes/{id:[0-9]+}",
        Method: http.MethodGet,
        Handler: func(w http.ResponseWriter, r *http.Request) {
            vars := mux.Vars(r)
            recipeID, err := strconv.Atoi(vars["id"])
            if err != nil {
                fmt.Printf("Recipe endpoint missing id: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            recipe, err := getRecipe(recipeID)
            if err != nil {
                fmt.Printf("FormError getting recipe: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            recipeIngredients, err := getIngredientsForRecipe(recipeID)
            if err != nil {
                fmt.Printf("FormError getting ingredients for recipe: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            recipeSteps, err := getStepsForRecipe(recipeID)
            if err != nil {
                fmt.Printf("FormError getting steps for recipe: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            sort.Sort(ByIngredientNumber(recipeIngredients))
            sort.Sort(ByStepNumber(recipeSteps))

            recipeBytes, err := json.Marshal(&RecipeResponse{
                Recipe:      *recipe,
                Ingredients: recipeIngredients,
                Steps:       recipeSteps,
            })
            if err != nil {
                fmt.Printf("FormError marshaling recipe: %s\n", err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            api.LogWriteErr(w.Write(recipeBytes))
        },
    }
}
