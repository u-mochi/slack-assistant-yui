import { all, fork } from "redux-saga/effects";
import todoist from "./Todoist";

export default function* root() {
  yield all([
    fork(todoist),
  ]);
}
