import { FormikErrors } from "formik";
import { createAction, createAsyncAction } from "typesafe-actions";
import { LoginFormValues } from "../../../views/components/Login";
import { RegistrationFormValues } from "../../../views/components/Registration";
import { LoginRequest, RegisterRequest, UserActionTypes } from "./types";

export const registerAsync = createAsyncAction(
    UserActionTypes.REGISTER_REQUEST,
    UserActionTypes.REGISTER_SUCCESS,
    UserActionTypes.REGISTER_FAILURE
)<[RegisterRequest, (errors: FormikErrors<RegistrationFormValues>) => void], undefined, Error>();

export const loginAsync = createAsyncAction(
    UserActionTypes.LOGIN_REQUEST,
    UserActionTypes.LOGIN_SUCCESS,
    UserActionTypes.LOGIN_FAILURE
)<[LoginRequest, (errors: FormikErrors<LoginFormValues>) => void], undefined, Error>();

export const logout = createAction(UserActionTypes.LOGOUT)<void>();
