import { ActionType, getType } from "typesafe-actions";
import * as user from "./actions";
import { loginAsync, registerAsync } from "./actions";
import { UserState } from "./types";

export type UserAction = ActionType<typeof user>;

const initialState: UserState = {
    registering: false,
    loggingIn: false,
    error: ""
};

const reducer = (state = initialState, action: UserAction) => {
    switch (action.type) {
        case getType(registerAsync.request):
            return {
                registering: true
            };
        case getType(registerAsync.success):
            return {};
        case getType(registerAsync.failure):
            return {
                error: action.payload.message
            };
        case getType(loginAsync.request):
            return {
                loggingIn: true
            };
        case getType(loginAsync.success):
            return {};
        case getType(loginAsync.failure):
            return {
                error: action.payload.message
            };
        default:
            return state;
    }
};

export { reducer as userReducer };
