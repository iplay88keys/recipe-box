import { recipesReducer } from "./reducers";
import { Ingredient, RecipeResponse, RecipeListResponse, RecipeActionTypes, Step } from "./types";
import { Error } from "../types";

describe("reducer", () => {
    it("should return the initial state", () => {
        const updatedState = recipesReducer(undefined, {type: "UNKNOWN"});

        expect(updatedState).toEqual({
            recipes: [],
            recipe: {} as RecipeResponse,
            loading: true,
            error: ""
        });
    });

    it("should handle FETCH_RECIPES_REQUEST", () => {
        const action = {
            type: RecipeActionTypes.FETCH_RECIPES_REQUEST
        };

        const updatedState = recipesReducer(undefined, action);

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
        const action = {
            type: RecipeActionTypes.FETCH_RECIPES_SUCCESS,
            payload: recipes
        };

        const updatedState = recipesReducer(undefined, action);

        expect(updatedState).toEqual({
            recipes: recipes.recipes,
            recipe: {} as RecipeResponse,
            loading: false,
            error: ""
        });
    });

    it("should handle FETCH_RECIPES_FAILURE", () => {
        let err = {
            error: "some error"
        } as Error;
        const action = {
            type: RecipeActionTypes.FETCH_RECIPES_FAILURE,
            payload: err
        };

        const updatedState = recipesReducer(undefined, action);

        expect(updatedState).toEqual({
            recipes: [],
            recipe: {} as RecipeResponse,
            loading: false,
            error: "some error"
        });
    });

    it("should handle FETCH_RECIPE_REQUEST", () => {
        const action = {
            type: RecipeActionTypes.FETCH_RECIPE_REQUEST
        };

        const updatedState = recipesReducer(undefined, action);

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
        const action = {
            type: RecipeActionTypes.FETCH_RECIPE_SUCCESS,
            payload: recipe
        };

        const updatedState = recipesReducer(undefined, action);

        expect(updatedState).toEqual({
            recipes: [],
            recipe: recipe,
            loading: false,
            error: ""
        });
    });

    it("should handle FETCH_RECIPE_FAILURE", () => {
        let err = {
            error: "some error"
        } as Error;
        const action = {
            type: RecipeActionTypes.FETCH_RECIPE_FAILURE,
            payload: err
        };

        const updatedState = recipesReducer(undefined, action);

        expect(updatedState).toEqual({
            recipes: [],
            recipe: {} as RecipeResponse,
            loading: false,
            error: "some error"
        });
    });
});
