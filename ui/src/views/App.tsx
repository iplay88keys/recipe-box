import React from "react";
import "bootstrap/dist/css/bootstrap.min.css";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import { Navigation } from "./components/Navigation";
import Recipes from "./pages/recipes";
import Recipe from "./pages/recipe";
import styled from "styled-components";

const StyledApp = styled.div`
  height: 100%;
  width: 75%;
  margin: auto;
  padding-top: 30px;
`;

const App: React.FC = () => {
    return (
        <BrowserRouter>
            <div>
                <Navigation/>
                <StyledApp>
                    <Switch>
                        <Route exact path="/" component={Recipes}/>
                        <Route exact path="/recipes" component={Recipes}/>
                        <Route exact path="/recipes/:recipeID" component={Recipe}/>
                    </Switch>
                </StyledApp>
            </div>
        </BrowserRouter>
    );
};

export default App;
