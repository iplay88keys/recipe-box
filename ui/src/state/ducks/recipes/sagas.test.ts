import axios from "axios";
import MockAdapter from "axios-mock-adapter";
import { Action } from "redux";
import { runSaga } from "redux-saga";
import { action, PayloadAction } from "typesafe-actions";
import { createRecipeAsync, fetchRecipeAsync, fetchRecipesAsync } from "./actions";
import { createRecipeSaga, getRecipeSaga, listRecipeSaga } from "./sagas";
import {
    Ingredient,
    RecipeResponse,
    RecipeActionTypes,
    RecipeCreateRequest,
    RecipeCreateResponse,
    RecipeListResponse,
    Step
} from "./types";

describe.only("recipes", () => {
    let mock: MockAdapter;

    beforeEach(() => {
        mock = new MockAdapter(axios);
    });

    afterEach(() => {
        mock.restore();
    });

    describe.only("listRecipeSaga", () => {
        it("returns the recipes and dispatches the success action", async () => {
            let dispatched: Action<string>[] = [];

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
            } as RecipeListResponse;
            const url = "/api/v1/recipes";

            mock.onGet(url).reply(200, recipes);

            await runSaga({
                dispatch: (action: Action<string>) => dispatched.push(action)
            }, listRecipeSaga).toPromise();

            expect(mock.history.get).toHaveLength(1);
            expect(mock.history.get[0].headers).toHaveProperty("Accept");
            expect(mock.history.get[0].headers["Accept"]).toEqual("application/json");

            expect(dispatched).toEqual([fetchRecipesAsync.success(recipes)] as Action<string>[]);
        });

        it("returns an error when there is a non-200 response and dispatches the error action", async () => {
            let dispatched: Action<string>[] = [];

            const url = "/api/v1/recipes";

            mock.onGet(url).reply(404);

            await runSaga({
                dispatch: (action: Action<string>) => dispatched.push(action)
            }, listRecipeSaga).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(RecipeActionTypes.FETCH_RECIPES_FAILURE);
        });
    });

    describe.only("getRecipeSaga", () => {
        it("returns a recipe and dispatches the success action", async () => {
            let dispatched: Action<string>[] = [];

            const recipe = {
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
            const url = "/api/v1/recipes/0";

            mock.onGet(url).reply(200, recipe);

            await runSaga({
                dispatch: (action: Action<string>) => dispatched.push(action)
            }, getRecipeSaga, action(RecipeActionTypes.FETCH_RECIPE_REQUEST, 0)).toPromise();

            expect(mock.history.get).toHaveLength(1);
            expect(mock.history.get[0].headers).toHaveProperty("Accept");
            expect(mock.history.get[0].headers["Accept"]).toEqual("application/json");

            expect(dispatched).toEqual([fetchRecipeAsync.success(recipe)] as Action<string>[]);
        });

        it("returns an error when there is a non-200 response and dispatches the error action", async () => {
            let dispatched: Action<string>[] = [];

            const url = "/api/v1/recipes/0";

            mock.onGet(url).reply(404);

            await runSaga({
                dispatch: (action: Action<string>) => dispatched.push(action)
            }, getRecipeSaga, action(RecipeActionTypes.FETCH_RECIPE_REQUEST, 0)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(RecipeActionTypes.FETCH_RECIPE_FAILURE);
        });
    });

    describe.only("createRecipeSaga", () => {
        it("dispatches the success action", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/recipes";
            const resp = {
                recipe_id: 10,
                errors: {}
            } as RecipeCreateResponse;

            const req = {
                name: "Root Beer",
                description: "Delicious",
                servings: 4
            } as RecipeCreateRequest;

            mock.onPost(url, req).reply(200, resp);

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, createRecipeSaga, createRecipeAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(RecipeActionTypes.CREATE_RECIPE_SUCCESS);

            expect(mockSetErrors).not.toHaveBeenCalled();
        });

        it("calls the provided function with the error payload if there is one", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/recipes";
            const resp = {
                recipe_id: 0,
                errors: {
                    first: "first error",
                    second: "second error"
                }
            } as RecipeCreateResponse;

            const req = {
                name: "Root Beer",
                description: "Delicious",
                servings: 4
            } as RecipeCreateRequest;

            mock.onPost(url, req).reply(400, resp);

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, createRecipeSaga, createRecipeAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(RecipeActionTypes.CREATE_RECIPE_FAILURE);

            expect(mockSetErrors).toHaveBeenCalledTimes(1);
            expect(mockSetErrors).toHaveBeenNthCalledWith(1, resp.errors);
        });

        it("returns an error if there is no payload and there is an error", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/recipes";

            const req = {
                name: "Root Beer",
                description: "Delicious",
                servings: 4
            } as RecipeCreateRequest;

            mock.onPost(url, req).reply(500);

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, createRecipeSaga, createRecipeAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(RecipeActionTypes.CREATE_RECIPE_FAILURE);

            expect(mockSetErrors).not.toHaveBeenCalled();
        });

        it("returns an error if there is a network error", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/recipes";

            const req = {
                name: "Root Beer",
                description: "Delicious",
                servings: 4
            } as RecipeCreateRequest;

            mock.onPost(url, req).networkError();

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, createRecipeSaga, createRecipeAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(RecipeActionTypes.CREATE_RECIPE_FAILURE);

            expect(mockSetErrors).not.toHaveBeenCalled();
        });
    });
});
