export interface RegisterRequest {
    username: string
    email: string
    password: string
}

export interface RegisterResponse {
    errors: APIError
}

export interface APIError {
    [key: string]: string
}

export interface UserState {
    registering: boolean
    error: string
}

export enum UserActionTypes {
    REGISTER_REQUEST = "@@user/REGISTER_REQUEST",
    REGISTER_SUCCESS = "@@user/REGISTER_SUCCESS",
    REGISTER_FAILURE = "@@user/REGISTER_FAILURE",
}
