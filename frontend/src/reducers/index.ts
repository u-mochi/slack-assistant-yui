import {Action, combineReducers} from "redux";
import {TodoistActions} from "../actions/Todoist";
import todoist, {ITodoistSate} from "./Todoist";

export default combineReducers({
  todoist,
});

export interface IReduxState {
  todoist: ITodoistSate;
}

export type ReduxAction = TodoistActions | Action;
