import { AxiosResponse } from "axios";
import { call, put, takeEvery } from "redux-saga/effects";
import Api from "../../../api/api";
import { fetchRecipeAsync, fetchRecipesAsync } from "./actions";
import { RecipeActionTypes, RecipeListResponse, RecipeResponse } from "./types";

export function* listRecipesSaga(): Generator {
    try {
        const response = (yield call(Api.get, "/api/v1/recipes")) as AxiosResponse;

        yield put(fetchRecipesAsync.success((response.data) as RecipeListResponse));
    } catch (err) {
        if (err.response && err.response.status === 401) {
            // log out
        }
        yield put(fetchRecipesAsync.failure(err));
    }
}

export function* getRecipeSaga(action: ReturnType<typeof fetchRecipeAsync.request>): Generator {
    try {
        const response = (yield call(Api.get, `/api/v1/recipes/${action.payload}`)) as AxiosResponse;

        yield put(fetchRecipeAsync.success((response.data) as RecipeResponse));
    } catch (err) {
        if (err.response && err.response.status === 401) {
            // log out
        }

        yield put(fetchRecipeAsync.failure(err));
    }
}

function* watchRequestRecipes() {
    yield takeEvery(RecipeActionTypes.FETCH_RECIPES_REQUEST, listRecipesSaga);
}

function* watchRequestRecipe() {
    yield takeEvery(RecipeActionTypes.FETCH_RECIPE_REQUEST, getRecipeSaga);
}

const recipeSagas = [watchRequestRecipe(), watchRequestRecipes()];
export default recipeSagas;
