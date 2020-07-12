import { Reducer } from "redux";
import { ActionType, getType } from "typesafe-actions";
import * as user from "./actions";
import { registerAsync } from "./actions";
import { UserState } from "./types";

export type UserAction = ActionType<typeof user>;

const initialState: UserState = {
    registering: false,
    error: ""
};

const reducer: Reducer<UserState, UserAction> = (state = initialState, action: UserAction) => {
    switch (action.type) {
        case getType(registerAsync.request):
            return {
                ...state,
                registering: true
            };
        case getType(registerAsync.success):
            return {
                ...state,
                registering: false
            };
        case getType(registerAsync.failure):
            return {
                ...state,
                registering: false,
                error: action.payload.message
            };
        default:
            return state;
    }
};

export { reducer as userReducer };
