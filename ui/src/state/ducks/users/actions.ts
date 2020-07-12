import { FormikErrors } from "formik";
import { createAsyncAction } from "typesafe-actions";
import { RegistrationFormValues } from "../../../views/components/Registration";
import { RegisterRequest, UserActionTypes } from "./types";

export const registerAsync = createAsyncAction(
    UserActionTypes.REGISTER_REQUEST,
    UserActionTypes.REGISTER_SUCCESS,
    UserActionTypes.REGISTER_FAILURE
)<[RegisterRequest, (errors: FormikErrors<RegistrationFormValues>) => void], undefined, Error>();
