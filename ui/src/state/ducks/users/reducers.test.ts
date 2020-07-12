import { registerAsync } from "./actions";
import { userReducer } from "./reducers";
import { RegisterRequest, RegisterResponse } from "./types";

describe("reducer", () => {
    it("should handle REGISTER_REQUEST", () => {
        const req = {
            username: "some-user",
            email: "test@example.com",
            password: "password"
        } as RegisterRequest;

        const mockSetErrors = jest.fn();

        const updatedState = userReducer(undefined, registerAsync.request(req, mockSetErrors));

        expect(updatedState).toEqual({
            "error": "",
            "registering": true
        });
    });

    it("should handle REGISTER_SUCCESS", () => {
        const resp = {
            errors: {
                first: "first error",
                second: "second error"
            }
        } as RegisterResponse;

        const updatedState = userReducer(undefined, registerAsync.success());
        expect(updatedState).toEqual({
            "error": "",
            "registering": false
        });
    });

    it("should handle REGISTER_FAILURE", () => {
        let err = {
            message: "some error"
        } as Error;

        const updatedState = userReducer(undefined, registerAsync.failure(err));

        expect(updatedState).toEqual({
            "error": "some error",
            "registering": false
        });
    });
});
