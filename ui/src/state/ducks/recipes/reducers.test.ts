import { fetchRecipeAsync, fetchRecipesAsync } from "./actions";
import { recipeReducer } from "./reducers";
import { Ingredient, RecipeListResponse, RecipeResponse, Step } from "./types";

describe("reducer", () => {
    it("should handle FETCH_RECIPES_REQUEST", () => {
        const updatedState = recipeReducer(undefined, fetchRecipesAsync.request());

        expect(updatedState).toEqual({
            recipes: [],
            recipe: {} as RecipeResponse,
            loading: true,
            error: ""
        });
    });

    it("should handle FETCH_RECIPES_SUCCESS", () => {
        let recipes = {
            recipes: [{
                id: 0,
                name: "First",
                description: "One"
            }] as RecipeResponse[]
        } as RecipeListResponse;

        const updatedState = recipeReducer(undefined, fetchRecipesAsync.success(recipes));

        expect(updatedState).toEqual({
            recipes: recipes.recipes,
            recipe: {} as RecipeResponse,
            loading: false,
            error: ""
        });
    });

    it("should handle FETCH_RECIPES_FAILURE", () => {
        let err = {
            message: "some error"
        } as Error;

        const updatedState = recipeReducer(undefined, fetchRecipesAsync.failure(err));

        expect(updatedState).toEqual({
            recipes: [],
            recipe: {} as RecipeResponse,
            loading: false,
            error: "some error"
        });
    });

    it("should handle FETCH_RECIPE_REQUEST", () => {
        const updatedState = recipeReducer(undefined, fetchRecipeAsync.request(1));

        expect(updatedState).toEqual({
            recipes: [],
            recipe: {} as RecipeResponse,
            loading: true,
            error: ""
        });
    });

    it("should handle FETCH_RECIPE_SUCCESS", () => {
        let recipe = {
            id: 0,
            name: "Root Beer Float",
            description: "Delicious",
            creator: "User1",
            servings: 1,
            prep_time: "5 m",
            total_time: "5 m",
            source: "Some Website",
            ingredients: [{
                ingredient: "Vanilla Ice Cream",
                ingredient_number: 0,
                amount: 1,
                measurement: "Scoop",
                preparation: "Frozen"
            }] as Ingredient[],
            steps: [{
                step_number: 1,
                instructions: "Place ice cream in glass."
            }] as Step[]
        } as RecipeResponse;

        const updatedState = recipeReducer(undefined, fetchRecipeAsync.success(recipe));

        expect(updatedState).toEqual({
            recipes: [],
            recipe: recipe,
            loading: false,
            error: ""
        });
    });

    it("should handle FETCH_RECIPE_FAILURE", () => {
        let err = {
            message: "some error"
        } as Error;

        const updatedState = recipeReducer(undefined, fetchRecipeAsync.failure(err));

        expect(updatedState).toEqual({
            recipes: [],
            recipe: {} as RecipeResponse,
            loading: false,
            error: "some error"
        });
    });
});
