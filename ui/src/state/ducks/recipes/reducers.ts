import { Reducer } from "redux";
import { RecipeActionTypes, RecipeResponse, RecipesState } from "./types";

const initialState: RecipesState = {
    recipes: [],
    recipe: {} as RecipeResponse,
    loading: true,
    error: ""
};

const reducer: Reducer<RecipesState> = (state = initialState, action) => {
    switch (action.type) {
        case RecipeActionTypes.FETCH_RECIPES_REQUEST: {
            return {...state, loading: true};
        }
        case RecipeActionTypes.FETCH_RECIPES_SUCCESS: {
            return {...state, loading: false, recipes: action.payload.recipes};
        }
        case RecipeActionTypes.FETCH_RECIPES_FAILURE: {
            return {...state, loading: false, error: action.payload.error};
        }
        case RecipeActionTypes.FETCH_RECIPE_REQUEST: {
            return {...state, loading: true};
        }
        case RecipeActionTypes.FETCH_RECIPE_SUCCESS: {
            return {...state, loading: false, recipe: action.payload};
        }
        case RecipeActionTypes.FETCH_RECIPE_FAILURE: {
            return {...state, loading: false, error: action.payload.error};
        }
        default: {
            return state;
        }
    }
};

export { reducer as recipesReducer };
