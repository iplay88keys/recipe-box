export interface RegisterRequest {
    username: string
    email: string
    password: string
}

export interface RegisterResponse {
    errors: APIError
}

export interface LoginRequest {
    login: string
    password: string
}

export interface LoginResponse {
    access_token: string
    refresh_token: string
    errors: APIError
}

export interface APIError {
    [key: string]: string
}

export interface UserState {
    registering: boolean
    loggingIn: boolean
    error: string
}

export enum UserActionTypes {
    REGISTER_REQUEST = "@@user/REGISTER_REQUEST",
    REGISTER_SUCCESS = "@@user/REGISTER_SUCCESS",
    REGISTER_FAILURE = "@@user/REGISTER_FAILURE",
    LOGIN_REQUEST = "@@user/LOGIN_REQUEST",
    LOGIN_SUCCESS = "@@user/LOGIN_SUCCESS",
    LOGIN_FAILURE = "@@user/LOGIN_FAILURE",
}
