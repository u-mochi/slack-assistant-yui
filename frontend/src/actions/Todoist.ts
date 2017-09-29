import {Action} from "redux";
import * as actions from "./index";

// Action Types
export const CONFIGURATION = actions.createRequestTypes("TODOIST_CONFIGURATION");
export const LOAD_CONFIGURATION = "TODOIST_CONFIGURATION_LOAD";
export const UPDATE_CONFIGURATION = "TODOIST_CONFIGURATION_UPDATE";

// Action
export interface ITodoistConfigurationAction extends Action {
  apiKey?: string;
  error?: Error;
}

export type TodoistActions = ITodoistConfigurationAction;

// Action Creators
export const configuration = {
  failure: (error: Error): ITodoistConfigurationAction => actions.action(CONFIGURATION.FAILURE, {error}),
  request: (): ITodoistConfigurationAction => actions.action(CONFIGURATION.REQUEST, {}),
  success: (apiKey: string): ITodoistConfigurationAction => actions.action(CONFIGURATION.SUCCESS, {apiKey}),
};
export const loadTodoistConfiguration = (): ITodoistConfigurationAction => actions.action(LOAD_CONFIGURATION, {});
export const updateTodoistConfiguration =
  (apiKey: string): ITodoistConfigurationAction => actions.action(UPDATE_CONFIGURATION, {apiKey});
