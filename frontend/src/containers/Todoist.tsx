import * as React from "react";
import {connect} from "react-redux";
import {Dispatch} from "redux";
import {loadTodoistConfiguration, updateTodoistConfiguration} from "../actions/Todoist";
import {TodoistApiKeyComponent} from "../components/Todoist";
import {IReduxState, ReduxAction} from "../reducers";

const TodoistApiKey = connect(
  (state: IReduxState) => {
    return {apiKey: state.todoist.apiKey};
  },
  (dispatch: Dispatch<ReduxAction>) => {
    return {
      onClick: (apiKey: string) => {
          dispatch(updateTodoistConfiguration(apiKey));
      },
      onLoad: () => {
        dispatch(loadTodoistConfiguration());
      },
    };
  },
)(TodoistApiKeyComponent);

export default class TodoistContainer extends React.Component {
  public render() {
    return(
      <div className="columns">
        <div className="row">
          <h2>Todoist</h2>
        </div>
        <form>
          <TodoistApiKey />
        </form>
      </div>
    );
  }
}
