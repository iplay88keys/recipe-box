import { combineReducers } from "redux";
import { all } from "redux-saga/effects";

import { recipesReducer } from "./recipes/reducers";
import recipesSagas from "./recipes/sagas";
import { RecipesState } from "./recipes/types";

export interface ApplicationState {
    recipes: RecipesState
}

export const createRootReducer = () => combineReducers({
    recipes: recipesReducer
});

export function* rootSaga() {
    yield all([...recipesSagas]);
}

