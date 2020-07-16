import { createMuiTheme, CssBaseline } from "@material-ui/core";
import { ThemeProvider } from "@material-ui/styles";
import React from "react";
import { connect } from "react-redux";
import { Redirect, Route, Router, Switch } from "react-router-dom";
import styled from "styled-components";
import { history } from "../helpers/history";
import { ApplicationState } from "../state/ducks";
import { LoggedInRedirect } from "./components/LoggedInRedirect";
import { Navigation } from "./components/Navigation";
import { PrivateRoute } from "./components/PrivateRoute";
import Login from "./pages/LoginPage";
import Recipe from "./pages/RecipePage";
import Recipes from "./pages/RecipesPage";
import Register from "./pages/RegisterPage";

const StyledApp = styled.div`
  height: 100%;
  width: 75%;
  margin: auto;
  padding-top: 30px;
`;

interface PropsFromState {
    loggedIn: boolean
}

interface PropsFromDispatch {}

interface State {}

type AllProps = PropsFromState & PropsFromDispatch & State

class App extends React.Component<AllProps, State> {
    // constructor(props: AllProps) {
    //     super(props);
    //
    //     history.listen(() => {
    //     });
    // }

    render() {
        const theme = createMuiTheme({
            palette: {
                type: "light"
            }
        });

        return (
            <Router history={history}>
                <ThemeProvider theme={theme}>
                    <CssBaseline/>
                    <div>
                        <Navigation loggedIn={this.props.loggedIn}/>
                        <StyledApp>
                            <Switch>
                                <Route exact path="/" component={Recipes}/>
                                <LoggedInRedirect exact path="/register" component={Register}/>
                                <Route exact path="/login" component={Login}/>
                                <PrivateRoute exact path="/recipes" component={Recipes}/>
                                <PrivateRoute exact path="/recipes/:recipeID" component={Recipe}/>
                                <Redirect from="*" to="/"/>
                            </Switch>
                        </StyledApp>
                    </div>
                </ThemeProvider>
            </Router>
        );
    }
}

const mapStateToProps = ({users}: ApplicationState) => ({
    loggedIn: users.loggedIn
});

const mapDispatchToProps = {};

export default connect(mapStateToProps, mapDispatchToProps)(App);
