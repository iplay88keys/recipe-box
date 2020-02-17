import {
    fetchRecipeAsync,
    fetchRecipesAsync

} from "./actions";
import { Ingredient, RecipeResponse, RecipeActionTypes, Step } from "./types";

describe("actions", () => {
    describe("recipes", () => {
        it("should create an action to request recipes", () => {
            const expectedAction = {
                type: RecipeActionTypes.FETCH_RECIPES_REQUEST
            };
            expect(fetchRecipesAsync.request()).toEqual(expectedAction);
        });

        it("should create a successful action to receive recipes", () => {
            const recipes = {
                recipes: [{
                    id: 0,
                    name: "First",
                    description: "One"
                }, {
                    id: 1,
                    name: "Second",
                    description: "Two"
                }] as RecipeResponse[]
            };
            const expectedAction = {
                type: RecipeActionTypes.FETCH_RECIPES_SUCCESS,
                payload: recipes
            };
            expect(fetchRecipesAsync.success(recipes)).toEqual(expectedAction);
        });

        it("should create an error action to receive recipes", () => {
            const err = Error("some error");
            const expectedAction = {
                type: RecipeActionTypes.FETCH_RECIPES_FAILURE,
                payload: err
            };
            expect(fetchRecipesAsync.failure(err)).toEqual(expectedAction);
        });
    });

    describe("recipe", () => {
        it("should create an action to request a recipe", () => {
            const expectedAction = {
                type: RecipeActionTypes.FETCH_RECIPE_REQUEST,
                payload: 0
            };
            expect(fetchRecipeAsync.request(0)).toEqual(expectedAction);
        });

        it("should create a successful action to receive recipes", () => {
            let recipe = {
                id: 0,
                name: "Root Beer Float",
                description: "Delicious",
                creator: "User1",
                servings: 1,
                prep_time: "5 m",
                total_time: "5 m",
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
            const expectedAction = {
                type: RecipeActionTypes.FETCH_RECIPE_SUCCESS,
                payload: recipe
            };
            expect(fetchRecipeAsync.success(recipe)).toEqual(expectedAction);
        });

        it("should create an error action to receive recipes", () => {
            const err = Error("some error");
            const expectedAction = {
                type: RecipeActionTypes.FETCH_RECIPE_FAILURE,
                payload: err
            };
            expect(fetchRecipeAsync.failure(err)).toEqual(expectedAction);
        });
    });
});
