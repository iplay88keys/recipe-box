import { registerAsync } from "./actions";
import { RegisterRequest, UserActionTypes } from "./types";

describe("actions", () => {
    describe("register", () => {
        it("should create an action for registering", () => {
            const req = {
                username: "some-user",
                email: "test@example.com",
                password: "password"
            } as RegisterRequest;

            const mockSetErrors = jest.fn();

            const expectedAction = {
                type: UserActionTypes.REGISTER_REQUEST,
                payload: req,
                meta: mockSetErrors
            };

            expect(registerAsync.request(req, mockSetErrors)).toEqual(expectedAction);
        });

        it("should create a successful action for a successful registration", () => {
            const expectedAction = {
                type: UserActionTypes.REGISTER_SUCCESS
            };
            expect(registerAsync.success()).toEqual(expectedAction);
        });

        it("should create an error action for an unsuccessful signup", () => {
            const err = Error("some error");
            const expectedAction = {
                type: UserActionTypes.REGISTER_FAILURE,
                payload: err
            };
            expect(registerAsync.failure(err)).toEqual(expectedAction);
        });
    });
});
