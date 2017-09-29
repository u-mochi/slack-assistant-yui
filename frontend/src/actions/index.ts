const REQUEST = "REQUEST";
const SUCCESS = "SUCCESS";
const FAILURE = "FAILURE";

interface IRequestActions {
  [key: string]: string;
  FAILURE: string;
  REQUEST: string;
  SUCCESS: string;
}

export function createRequestTypes(base: string) {
  return [REQUEST, SUCCESS, FAILURE].reduce((acc, type) => {
    acc[type] = `${base}_${type}`;
    return acc;
  }, {} as IRequestActions);
}

export function action(type: string, payload = {}): {[key: string]: any, type: string } {
  return {type, ...payload};
}
