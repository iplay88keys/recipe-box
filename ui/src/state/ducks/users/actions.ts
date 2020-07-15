import { FormikErrors } from "formik";
import { createAsyncAction } from "typesafe-actions";
import { RegistrationFormValues } from "../../../views/components/Registration";
import { UserLoginFormValues } from "../../../views/components/UserLogin";
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
)<[LoginRequest, (errors: FormikErrors<UserLoginFormValues>) => void], undefined, Error>();
