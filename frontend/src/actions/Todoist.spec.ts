import {
  configuration,
  CONFIGURATION,
  ITodoistConfigurationAction,
  LOAD_CONFIGURATION,
  loadTodoistConfiguration,
  UPDATE_CONFIGURATION,
  updateTodoistConfiguration} from "./Todoist";

describe("Actions of Todoist", () => {
  let action: ITodoistConfigurationAction;
  const API_KEY = "api key";
  const ERROR = new Error("error");

  describe("Configuration Requests", () => {
    it("configuration.request", () => {
      action = configuration.request();
      expect(action.type).toEqual(CONFIGURATION.REQUEST);
    });

    it("configuration.success", () => {
      action = configuration.success(API_KEY);
      expect(action.type).toEqual(CONFIGURATION.SUCCESS);
      expect(action.apiKey).toEqual(API_KEY);
    });

    it("configuration.failure", () => {
      action = configuration.failure(ERROR);
      expect(action.type).toEqual(CONFIGURATION.FAILURE);
      expect(action.error).toEqual(ERROR);
    });
  });

  describe("loadTodoistConfiguration", () => {
    it("Load Configuration", () => {
      action = loadTodoistConfiguration();
      expect(action.type).toEqual(LOAD_CONFIGURATION);
    });
  });

  describe("updateTodoistConfiguration", () => {
    it("Update Configuration", () => {
      action = updateTodoistConfiguration(API_KEY);
      expect(action.type).toEqual(UPDATE_CONFIGURATION);
      expect(action.apiKey).toEqual(API_KEY);
    });
  });
});
