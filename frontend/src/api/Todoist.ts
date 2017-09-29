import "isomorphic-fetch";
import {API_HOST} from "./index";

export function loadConfiguration() {
  return fetch(API_HOST + "/api/todoist/configuration", {
    headers: {
      "Accept": "application/json, text/plain, */*",
      "Content-Type": "application/json",
    },
  })
  .then((response) => {
    return validateResponse(response);
  })
  .then((json) => {
    validateJson(json);
    return {
      error: undefined,
      payload: json[0].api_key,
    };
  })
  .catch((err) => {
    return {
      error: err,
      payload: undefined,
    };
  });
}

export function updateConfiguration(apiKey: string) {
  return fetch(API_HOST + "/api/todoist/configuration", {
    body: JSON.stringify({api_key: apiKey}),
    headers: {
      "Accept": "application/json, text/plain, */*",
      "Content-Type": "application/json",
    },
    method: "put",
  })
  .then((response) => {
    return validateResponse(response);
  })
  .then((json) => {
    validateJson(json);
    return {
      error: undefined,
      payload: json[0].api_key,
    };
  })
  .catch((err) => {
    return {
      error: err,
      payload: undefined,
    };
  });
}

function validateResponse(response: Response): Promise<any> {
  if (! response.ok) {
    return response.text().then((text) => {
      throw new Error(text);
    });
  }
  return response.json();
}

function validateJson(json: any) {
  if (json.length < 1) {
    throw new Error("JSON is empty " + JSON.stringify(json));
  }
  if (! json[0].hasOwnProperty("api_key")) {
    throw new Error("api_key is not exists in " + JSON.stringify(json));
  }
}
