import * as fetchMock from "fetch-mock";
import {API_HOST} from "./index";
import {loadConfiguration, updateConfiguration} from "./Todoist";

describe("API Todoist", () => {

  const API_KEY = "api key";
  const ERROR_MESSAGE = "Internal Server Error";
  const RESPONSE_JSON = {
    body: [{
      UpdateDate: "2017-08-03T05:55:50.581189Z",
    }],
    status: 200,
  };

  describe("Fetch configuration", () => {

    afterEach(fetchMock.restore);

    it("Test fetch", () => {
      fetchMock.mock(API_HOST + "/api/todoist/configuration", {
        body: [{
          UpdateDate: "2017-08-03T05:55:50.581189Z",
          api_key: API_KEY,
        }],
        status: 200,
      });

      loadConfiguration()
      .then((response) => {
        expect(response.payload).toEqual(API_KEY);
      });
    });

    it("Test fetch error1 (status:500)", () => {
      fetchMock.mock(API_HOST + "/api/todoist/configuration", {
        body: ERROR_MESSAGE,
        status: 500,
      });

      loadConfiguration()
      .then((response) => {
        expect(response.error.message).toEqual(ERROR_MESSAGE);
      });
    });

    it("Test fetch error2 (invalid JSON)", () => {
      fetchMock.mock(API_HOST + "/api/todoist/configuration", RESPONSE_JSON);

      loadConfiguration()
      .then((response) => {
        expect(response.error.message).toContain(JSON.stringify(RESPONSE_JSON.body));
      });
    });
  });

  describe("Put configuration", () => {

    afterEach(fetchMock.restore);

    it("Test put", () => {
      fetchMock.put(API_HOST + "/api/todoist/configuration", {
        body: [{
          UpdateDate: "2017-08-03T05:55:50.581189Z",
          api_key: API_KEY,
        }],
        status: 200,
      });

      updateConfiguration(API_KEY)
      .then((response) => {
        expect(response.payload).toEqual(API_KEY);
      });
    });

    it("Test put error1 (status:500)", () => {
      fetchMock.put(API_HOST + "/api/todoist/configuration", {
        body: ERROR_MESSAGE,
        status: 500,
      });

      updateConfiguration(API_KEY)
      .then((response) => {
        expect(response.error.message).toEqual(ERROR_MESSAGE);
      });
    });

    it("Test put error2 (invalid JSON)", () => {
      fetchMock.put(API_HOST + "/api/todoist/configuration", RESPONSE_JSON);

      updateConfiguration(API_KEY)
      .then((response) => {
        expect(response.error.message).toContain(JSON.stringify(RESPONSE_JSON.body));
      });
    });
  });
});
