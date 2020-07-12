import React from "react";
import { connect } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router";
import { ApplicationState } from "../../state/ducks";
import { registerAsync } from "../../state/ducks/users/actions";
import Registration from "../components/Registration";

interface PropsFromDispatch {
    register: typeof registerAsync.request
}

interface State {}

type AllProps = PropsFromDispatch & State & RouteComponentProps

class Register extends React.Component<AllProps, State> {
    constructor(props: AllProps) {
        super(props);
    }

    render() {
        return (
            <Registration
                register={this.props.register}
            />
        );
    }
}

const mapStateToProps = ({}: ApplicationState) => ({});

const mapDispatchToProps = {
    register: registerAsync.request
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Register));
