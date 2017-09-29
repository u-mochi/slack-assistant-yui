import {action, createRequestTypes} from "./index";

describe("Actions", () => {
  it("Create Request Type", () => {
    const BASE = "base";
    const requestType = createRequestTypes(BASE);
    expect(requestType.FAILURE).toEqual(`${BASE}_FAILURE`);
    expect(requestType.REQUEST).toEqual(`${BASE}_REQUEST`);
    expect(requestType.SUCCESS).toEqual(`${BASE}_SUCCESS`);
  });

  it("Action", () => {
    const TYPE = "type";
    const PAYLOAD: { [key: string]: any } = {
      a: 1,
      b: 2,
      c: 3,
    };
    const result = action(TYPE, PAYLOAD);
    expect(result.type).toEqual(TYPE);
    for (const prop in PAYLOAD) {
      if (result.hasOwnProperty(prop)) {
        expect(result[prop]).toEqual(PAYLOAD[prop]);
      } else {
        fail("Property " + prop + "not exists.");
      }
    }
  });
});
