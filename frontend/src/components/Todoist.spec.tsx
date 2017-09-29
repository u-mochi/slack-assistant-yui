import {mount} from "enzyme";
import * as React from "react";
import {TodoistApiKeyComponent} from "./todoist";

describe("Todoist components", () => {
  describe("Todoist API key component", () => {
    const API_KEY1 = "Test Key1";
    const API_KEY2 = "Test Key2";
    const spy1 = jasmine.createSpy("onClickSpy");
    const spy2 = jasmine.createSpy("onLoadSpy");
    const apiKeyComponent = mount(<TodoistApiKeyComponent
      onClick={(apiKey: string) => { spy1(); }}
      onLoad={() => { spy2(); }}
      apiKey="" />);

    it("Test set API key", () => {
      const apiKeyInput = apiKeyComponent.find("#todoist-apikey-input").getDOMNode() as HTMLInputElement;
      expect(spy2).toHaveBeenCalled();
      apiKeyComponent.setProps({apiKey: API_KEY1});
      expect(apiKeyInput.value).toEqual(API_KEY1);
      apiKeyComponent.setProps({apiKey: API_KEY2});
      expect(apiKeyInput.value).toEqual(API_KEY2);
    });
    it("Test click", () => {
      apiKeyComponent.find("#todoist-apikey-post").simulate("click");
      expect(spy1).toHaveBeenCalled();
    });
  });
});
