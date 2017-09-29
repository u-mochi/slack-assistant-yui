import { all, call, fork, put, take } from "redux-saga/effects";
import * as actions from "../actions/Todoist";
import {loadConfiguration as loadConf, updateConfiguration as updateConf} from "../api/Todoist";
import todoistRoot, {loadConfiguration, updateConfiguration,
  watchLoadConfiguration, watchUpdateConfiguration} from "./Todoist";

describe("Todoist sagas", () => {
  describe("Saga: loadConfiguration", () => {
    let generator = loadConfiguration();
    const API_KEY = "Test Key";
    const ERROR = new Error("Error Message");

    beforeEach(() => {
      generator = loadConfiguration();
    });

    it("Load configuration success", () => {
      expect(generator.next().value).toEqual(call(loadConf));
      expect(generator.next({payload: API_KEY}).value).toEqual(put(actions.configuration.success(API_KEY)));
      expect(generator.next().done).toBeTruthy();
    });

    it("Load configuration fail", () => {
      expect(generator.next().value).toEqual(call(loadConf));
      expect(generator.next({error: ERROR}).value).toEqual(put(actions.configuration.failure(ERROR)));
      expect(generator.next().done).toBeTruthy();
    });
  });

  describe("Saga: updateConfiguration", () => {
    const API_KEY = "Test Key";
    const ERROR = new Error("Error Message");
    let generator = updateConfiguration(API_KEY);

    beforeEach(() => {
      generator = updateConfiguration(API_KEY);
    });

    it("Update configuration success", () => {
      expect(generator.next().value).toEqual(call(updateConf, API_KEY));
      expect(generator.next({payload: API_KEY}).value).toEqual(put(actions.configuration.success(API_KEY)));
      expect(generator.next().done).toBeTruthy();
    });

    it("Update configuration fail", () => {
      expect(generator.next().value).toEqual(call(updateConf, API_KEY));
      expect(generator.next({error: ERROR}).value).toEqual(put(actions.configuration.failure(ERROR)));
      expect(generator.next().done).toBeTruthy();
    });
  });

  describe("Saga: watchLoadConfiguration", () => {
    it("lanches loadConfiguration Task", () => {
      const generator = watchLoadConfiguration();
      expect(generator.next().value).toEqual(take(actions.LOAD_CONFIGURATION));
      expect(generator.next().value).toEqual(fork(loadConfiguration));
      expect(generator.next().done).toBeFalsy();
    });
  });

  describe("Saga: watchUpdateConfiguration", () => {
    it("lanches updateConfiguration Task", () => {
      const API_KEY = "Test Key";
      const generator = watchUpdateConfiguration();
      expect(generator.next().value).toEqual(take(actions.UPDATE_CONFIGURATION));
      expect(generator.next({apiKey: API_KEY}).value).toEqual(fork(updateConfiguration, API_KEY));
      expect(generator.next().done).toBeFalsy();
    });
  });

  describe("Saga: todoistRoot", () => {
    it("lanches child Tasks", () => {
      const generator = todoistRoot();
      expect(generator.next().value).toEqual(all([
        fork(watchLoadConfiguration),
        fork(watchUpdateConfiguration)]));
      expect(generator.next().done).toBeTruthy();
    });
  });

});
