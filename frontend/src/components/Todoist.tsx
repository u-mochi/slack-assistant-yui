import * as React from "react";

export interface ITodoistApiKeyComponentProps {
    apiKey: string;
    onClick: (apiKey: string) => void;
    onLoad: () => void;
}

export class TodoistApiKeyComponent extends React.Component<ITodoistApiKeyComponentProps, any> {
  public refs: {
    apikey: HTMLInputElement;
  };
  public componentDidMount() {
    this.props.onLoad();
  }
  public componentWillReceiveProps(nextProps: ITodoistApiKeyComponentProps, nextContext: any) {
    this.refs.apikey.value = nextProps.apiKey;
  }
  public render() {
    return (
      <div id="todoist-apikey">
        <div className="row">
          <div className="columns small-4">
            <label htmlFor="todoist-apikey-input" className="text-right middle">Todoist API Key:</label>
          </div>
          <div className="columns small-8">
            <input
              id="todoist-apikey-input"
              type="text"
              placeholder="Your Todoist API key here."
              ref="apikey"
              className="eight"/>
          </div>
        </div>
        <div className="row">
          <div className="columns small-2 small-offset-10">
            <button
              id="todoist-apikey-post"
              type="button"
              className="success button float-right"
              onClick={() => {
                this.props.onClick(this.refs.apikey.value.trim());
              }}>
              Save
            </button>
          </div>
        </div>
      </div>
    );
  }
}
