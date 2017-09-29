import { all, call, fork, put, take } from "redux-saga/effects";
import * as actions from "../actions/Todoist";
import {loadConfiguration as loadConf, updateConfiguration as updateConf} from "../api/Todoist";

export function* loadConfiguration() {
  const { payload, error } = yield call(loadConf);
  if (payload && !error) {
    yield put(actions.configuration.success(payload));
  } else {
    yield put(actions.configuration.failure(error));
  }
}

export function* updateConfiguration(apiKey: string) {
  const { payload, error } = yield call(updateConf, apiKey);
  if (payload && !error) {
    yield put(actions.configuration.success(payload));
  } else {
    yield put(actions.configuration.failure(error));
  }
}

export function* watchLoadConfiguration() {
  while (true) {
    yield take(actions.LOAD_CONFIGURATION);
    yield fork(loadConfiguration);
  }
}

export function* watchUpdateConfiguration() {
  while (true) {
    const {apiKey} = yield take(actions.UPDATE_CONFIGURATION);
    yield fork(updateConfiguration, apiKey);
  }
}

export default function* todoistRoot() {
  yield all([
    fork(watchLoadConfiguration),
    fork(watchUpdateConfiguration),
  ]);
}
