import React from "react";
import { connect } from "react-redux";
import { ApplicationState } from "../../state/ducks";
import { fetchRecipesAsync } from "../../state/ducks/recipes/actions";
import { RecipeResponse } from "../../state/ducks/recipes/types";
import RecipeList from "../components/RecipeList";

interface PropsFromState {
    recipes: RecipeResponse[]
    loading: boolean
}

interface PropsFromDispatch {
    fetchRecipes: typeof fetchRecipesAsync.request
}

interface State {}

type AllProps = PropsFromState & PropsFromDispatch

class RecipesPage extends React.Component<AllProps, State> {
    componentDidMount() {
        const {fetchRecipes} = this.props;

        fetchRecipes();
    }

    render() {
        return <RecipeList
            recipes={this.props.recipes}
            loading={this.props.loading}
        />;
    }
}

const mapStateToProps = ({recipes}: ApplicationState) => ({
    recipes: recipes.recipes,
    loading: recipes.loading
});

const mapDispatchToProps = {
    fetchRecipes: fetchRecipesAsync.request
};

export default connect(mapStateToProps, mapDispatchToProps)(RecipesPage);
