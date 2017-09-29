import { all, fork } from "redux-saga/effects";
import rootSaga from "./index";
import todoistRoot from "./Todoist";

describe("Root sagas", () => {
  describe("Saga: Root", () => {
    it("lanches child Tasks", () => {
      const generator = rootSaga();
      expect(generator.next().value).toEqual(all([
        fork(todoistRoot)]));
      expect(generator.next().done).toBeTruthy();
    });
  });
});
