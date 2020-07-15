import { call, put, takeEvery } from "redux-saga/effects";
import Api from "../../../api/api";
import { history } from "../../../helpers/history";
import { loginAsync, registerAsync } from "./actions";
import { LoginResponse, UserActionTypes } from "./types";

export function* registerSaga(action: ReturnType<typeof registerAsync.request>): Generator {
    try {
        yield call(Api.post, "/api/v1/users/register", JSON.stringify(action.payload), false);

        yield put(registerAsync.success());
        history.push("/login");
    } catch (err) {
        if (err.response && err.response.data && err.response.data.errors) {
            action.meta(err.response.data.errors);
        }

        yield put(registerAsync.failure(err));
    }
}

export function* loginSaga(action: ReturnType<typeof loginAsync.request>): Generator {
    try {
        const response = (yield call(Api.post, "/api/v1/users/login", JSON.stringify(action.payload), false)) as LoginResponse;

        yield put(loginAsync.success());

        localStorage.setItem("access_token", response.access_token);
        history.push("/recipes");
    } catch (err) {
        if (err.response && err.response.data && err.response.data.errors) {
            action.meta(err.response.data.errors);
        }

        yield put(registerAsync.failure(err));
    }
}

function* watchRegister() {
    yield takeEvery(UserActionTypes.REGISTER_REQUEST, registerSaga);
}

function* watchLogin() {
    yield takeEvery(UserActionTypes.LOGIN_REQUEST, loginSaga);
}

const userSagas = [watchRegister(), watchLogin()];
export default userSagas;
