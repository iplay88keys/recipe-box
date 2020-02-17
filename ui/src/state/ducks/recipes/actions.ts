import { createAsyncAction } from "typesafe-actions";
import { RecipeResponse, RecipeListResponse, RecipeActionTypes } from "./types";

export const fetchRecipesAsync = createAsyncAction(
    RecipeActionTypes.FETCH_RECIPES_REQUEST,
    RecipeActionTypes.FETCH_RECIPES_SUCCESS,
    RecipeActionTypes.FETCH_RECIPES_FAILURE
)<undefined, RecipeListResponse, Error>();

export const fetchRecipeAsync = createAsyncAction(
    RecipeActionTypes.FETCH_RECIPE_REQUEST,
    RecipeActionTypes.FETCH_RECIPE_SUCCESS,
    RecipeActionTypes.FETCH_RECIPE_FAILURE
)<number, RecipeResponse, Error>();
