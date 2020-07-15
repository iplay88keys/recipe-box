import React from "react";
import { connect } from "react-redux";
import { ApplicationState } from "../../state/ducks";
import { loginAsync } from "../../state/ducks/users/actions";
import UserLogin from "../components/UserLogin";

interface PropsFromDispatch {
    login: typeof loginAsync.request
}

interface State {}

type AllProps = PropsFromDispatch & State

class Login extends React.Component<AllProps, State> {
    render() {
        return (
            <UserLogin
                login={this.props.login}
            />
        );
    }
}

const mapStateToProps = ({}: ApplicationState) => ({});

const mapDispatchToProps = {
    login: loginAsync.request
};

export default connect(mapStateToProps, mapDispatchToProps)(Login);
