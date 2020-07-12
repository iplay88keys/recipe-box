import { Reducer } from "redux";
import { ActionType, getType } from "typesafe-actions";
import * as recipe from "./actions";
import { fetchRecipeAsync, fetchRecipesAsync } from "./actions";
import { RecipeResponse, RecipeState } from "./types";

export type RecipeAction = ActionType<typeof recipe>;

const initialState: RecipeState = {
    recipes: [],
    recipe: {} as RecipeResponse,
    loading: true,
    error: ""
};

const reducer: Reducer<RecipeState, RecipeAction> = (state = initialState, action: RecipeAction) => {
    switch (action.type) {
        case getType(fetchRecipesAsync.request):
            return {
                ...state,
                loading: true
            };
        case getType(fetchRecipesAsync.success):
            return {
                ...state,
                loading: false,
                recipes: action.payload.recipes
            };
        case getType(fetchRecipesAsync.failure):
            return {
                ...state,
                loading: false,
                error: action.payload.message
            };
        case getType(fetchRecipeAsync.request):
            return {
                ...state,
                loading: true
            };
        case getType(fetchRecipeAsync.success):
            return {
                ...state,
                loading: false,
                recipe: action.payload
            };
        case getType(fetchRecipeAsync.failure):
            return {
                ...state,
                loading: false,
                error: action.payload.message
            };
        default:
            return state
    }
};

export { reducer as recipeReducer };
