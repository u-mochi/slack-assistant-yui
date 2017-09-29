import {configuration} from "../actions/Todoist";
import reducer from "./Todoist";

describe("Reducers Todoist", () => {
  const API_KEY = "api key";
  const ERROR_MESSAGE = "error occurred";

  it("configuration.request", () => {
    expect(reducer(undefined, configuration.request()).isFetching).toBeTruthy();
  });

  it("configuration.success", () => {
    const state = reducer(undefined, configuration.success(API_KEY));
    expect(state.apiKey).toEqual(API_KEY);
    expect(state.isFetching).toBeFalsy();
  });

  it("configuration.failure", () => {
    const state = reducer(undefined, configuration.failure(new Error(ERROR_MESSAGE)));
    expect(state.errorMessage).toEqual(ERROR_MESSAGE);
    expect(state.isFetching).toBeFalsy();
  });
});
