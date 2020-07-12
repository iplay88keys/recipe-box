import { call, put, takeEvery } from "redux-saga/effects";
import Api from "../../../api/api";
import { history } from "../../../helpers/history";
import { registerAsync } from "./actions";
import { UserActionTypes } from "./types";

export function* registerSaga(action: ReturnType<typeof registerAsync.request>): Generator {
    try {
        yield call(Api.post, "/api/v1/users/register", JSON.stringify(action.payload));

        yield put(registerAsync.success());
        history.push("/login");
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

const userSagas = [watchRegister()];
export default userSagas;
