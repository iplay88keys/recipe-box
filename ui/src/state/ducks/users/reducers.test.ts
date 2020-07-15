import { loginAsync, registerAsync } from "./actions";
import { userReducer } from "./reducers";
import { LoginRequest, RegisterRequest } from "./types";

describe("reducer", () => {
    describe("register", function () {
        it("should handle REGISTER_REQUEST", () => {
            const req = {
                username: "some-user",
                email: "test@example.com",
                password: "password"
            } as RegisterRequest;

            const mockSetErrors = jest.fn();

            const updatedState = userReducer(undefined, registerAsync.request(req, mockSetErrors));

            expect(updatedState).toEqual({
                "registering": true
            });
        });

        it("should handle REGISTER_SUCCESS", () => {
            const updatedState = userReducer(undefined, registerAsync.success());
            expect(updatedState).toEqual({});
        });

        it("should handle REGISTER_FAILURE", () => {
            let err = {
                message: "some error"
            } as Error;

            const updatedState = userReducer(undefined, registerAsync.failure(err));

            expect(updatedState).toEqual({
                "error": "some error"
            });
        });
    });

    describe("login", function () {
        it("should handle LOGIN_REQUEST", () => {
            const req = {
                login: "some-user",
                password: "password"
            } as LoginRequest;

            const mockSetErrors = jest.fn();

            const updatedState = userReducer(undefined, loginAsync.request(req, mockSetErrors));

            expect(updatedState).toEqual({
                "loggingIn": true
            });
        });

        it("should handle LOGIN_SUCCESS", () => {
            const updatedState = userReducer(undefined, loginAsync.success());
            expect(updatedState).toEqual({});
        });

        it("should handle LOGIN_FAILURE", () => {
            let err = {
                message: "some error"
            } as Error;

            const updatedState = userReducer(undefined, loginAsync.failure(err));
            expect(updatedState).toEqual({
                "error": "some error"
            });
        });
    });
});
