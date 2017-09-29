import * as React from "react";
import * as ReactDOM from "react-dom";
import {Provider} from "react-redux";
import Todoist from "./containers/Todoist";

import configureStore from "./store";

ReactDOM.render(
  <Provider store={configureStore()}>
    <Todoist />
  </Provider>
  , document.getElementById("todoist-configuration"),
);
