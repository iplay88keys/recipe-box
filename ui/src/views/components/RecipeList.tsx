import React from "react";
import Table from "react-bootstrap/Table";
import { RouteComponentProps } from "react-router";
import styled from "styled-components";
import { RecipeResponse } from "../../state/ducks/recipes/types";
import { withRouter } from "react-router-dom";

const StyledTR = styled.tr`
    cursor: pointer;
`;

interface RecipeListProps extends RouteComponentProps {
    recipes: RecipeResponse[]
    loading: boolean
}

export const RecipeList = ({recipes, loading, history}: RecipeListProps) => {
    if (loading) {
        return (
            <div>
                <p>Loading recipes</p>
            </div>
        );
    }

    return (
        <Table striped bordered hover>
            <thead>
            <tr>
                <th>Recipe</th>
                <th>Description</th>
            </tr>
            </thead>
            <tbody>
            {recipes.map((recipe: RecipeResponse) =>
                <StyledTR key={recipe.id} onClick={() => history.push(`/recipes/${recipe.id}`)}>
                    <td>{recipe.name}</td>
                    {recipe.description != null && <td>{recipe.description}</td>}
                </StyledTR>
            )}
            </tbody>
        </Table>
    );
};

export default withRouter(RecipeList);
