import {CONFIGURATION, LOAD_CONFIGURATION, TodoistActions, UPDATE_CONFIGURATION} from "../actions/Todoist";

export interface ITodoistSate {
  apiKey: string;
  errorMessage: string;
  isFetching: boolean;
  message: string;
}

const initialSate: ITodoistSate = {
  apiKey: "",
  errorMessage: "",
  isFetching: false,
  message: "",
};

export default function reducer(state: ITodoistSate = initialSate, action: TodoistActions): ITodoistSate {
  switch (action.type) {
    case CONFIGURATION.REQUEST:
      return {
        ...state,
        isFetching: true,
      };
    case CONFIGURATION.SUCCESS:
      return {
        ...state,
        apiKey: action.apiKey ? action.apiKey : "",
        isFetching: false,
      };
    case CONFIGURATION.FAILURE:
      return {
        ...state,
        errorMessage: action.error ? action.error.message : "",
        isFetching: false,
      };
    case LOAD_CONFIGURATION:
    case UPDATE_CONFIGURATION:
    default:
      return state;
  }
}
