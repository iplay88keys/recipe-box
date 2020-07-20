import axios from "axios";
import MockAdapter from "axios-mock-adapter";
import { runSaga } from "redux-saga";
import { PayloadAction } from "typesafe-actions";
import { loginAsync, registerAsync } from "./actions";
import { loginSaga, registerSaga } from "./sagas";
import { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse, UserActionTypes } from "./types";

describe.only("users", () => {
    let mock: MockAdapter;

    beforeEach(() => {
        mock = new MockAdapter(axios);
    });

    afterEach(() => {
        mock.restore();
    });

    describe.only("registerSaga", () => {
        it("dispatches the success action", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/users/register";

            const req = {
                username: "some-user",
                email: "test@example.com",
                password: "password"
            } as RegisterRequest;

            mock.onPost(url, req).reply(200);

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, registerSaga, registerAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(UserActionTypes.REGISTER_SUCCESS);

            expect(mockSetErrors).not.toHaveBeenCalled();
        });

        it("calls the provided function with the error payload if there is one", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/users/register";
            const resp = {
                errors: {
                    first: "first error",
                    second: "second error"
                }
            } as RegisterResponse;

            const req = {
                username: "some-user",
                email: "test@example.com",
                password: "password"
            } as RegisterRequest;

            mock.onPost(url, req).reply(400, resp);

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, registerSaga, registerAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(UserActionTypes.REGISTER_FAILURE);

            expect(mockSetErrors).toHaveBeenCalledTimes(1);
            expect(mockSetErrors).toHaveBeenNthCalledWith(1, resp.errors);
        });

        it("returns an error if there is no payload and there is an error", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/users/register";

            const req = {
                username: "some-user",
                email: "test@example.com",
                password: "password"
            } as RegisterRequest;

            mock.onPost(url, req).reply(500);

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, registerSaga, registerAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(UserActionTypes.REGISTER_FAILURE);

            expect(mockSetErrors).not.toHaveBeenCalled();
        });

        it("returns an error if there is a network error", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/users/register";

            const req = {
                username: "some-user",
                email: "test@example.com",
                password: "password"
            } as RegisterRequest;

            mock.onPost(url, req).networkError();

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, registerSaga, registerAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(UserActionTypes.REGISTER_FAILURE);

            expect(mockSetErrors).not.toHaveBeenCalled();
        });
    });

    describe.only("loginSaga", () => {
        it("dispatches the success action", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/users/login";
            const resp = {
                access_token: "some-token",
                refresh_token: "another-token",
                errors: {}
            } as LoginResponse;

            const req = {
                login: "some-user",
                password: "password"
            } as LoginRequest;

            mock.onPost(url, req).reply(200, resp);

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, loginSaga, loginAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(UserActionTypes.LOGIN_SUCCESS);

            expect(mockSetErrors).not.toHaveBeenCalled();
        });

        it("calls the provided function with the error payload if there is one", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/users/login";
            const resp = {
                access_token: "",
                refresh_token: "",
                errors: {
                    first: "first error",
                    second: "second error"
                }
            } as LoginResponse;

            const req = {
                login: "some-user",
                password: "password"
            } as LoginRequest;

            mock.onPost(url, req).reply(400, resp);

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, loginSaga, loginAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(UserActionTypes.LOGIN_FAILURE);

            expect(mockSetErrors).toHaveBeenCalledTimes(1);
            expect(mockSetErrors).toHaveBeenNthCalledWith(1, resp.errors);
        });

        it("returns an error if there is no payload and there is an error", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/users/login";

            const req = {
                login: "some-user",
                password: "password"
            } as LoginRequest;

            mock.onPost(url, req).reply(500);

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, loginSaga, loginAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(UserActionTypes.LOGIN_FAILURE);

            expect(mockSetErrors).not.toHaveBeenCalled();
        });

        it("returns an error if there is a network error", async () => {
            let dispatched: PayloadAction<string, string>[] = [];

            const url = "/api/v1/users/login";

            const req = {
                login: "some-user",
                password: "password"
            } as LoginRequest;

            mock.onPost(url, req).networkError();

            const mockSetErrors = jest.fn();

            await runSaga({
                dispatch: (action: PayloadAction<string, string>) => dispatched.push(action)
            }, loginSaga, loginAsync.request(req, mockSetErrors)).toPromise();

            expect(dispatched).toHaveLength(1);
            expect(dispatched[0].type).toEqual(UserActionTypes.LOGIN_FAILURE);

            expect(mockSetErrors).not.toHaveBeenCalled();
        });
    });
});
