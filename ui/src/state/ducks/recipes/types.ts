export interface RecipeListResponse {
    recipes: RecipeResponse[]
}

export interface RecipeResponse {
    id: number
    name: string
    description: string
    creator: string
    servings: number
    prep_time: string
    cook_time: string
    cool_time: string
    total_time: string
    source: string
    ingredients: Ingredient[]
    steps: Step[]
}

export interface Ingredient {
    ingredient: string
    ingredient_number: number
    amount: number
    measurement: string
    preparation: string
}

export interface Step {
    step_number: number
    instructions: string
}

export interface RecipeState {
    recipes: RecipeResponse[]
    recipe: RecipeResponse
    loading: boolean
    error: string
}

export enum RecipeActionTypes {
    FETCH_RECIPES_REQUEST = "@@recipes/FETCH_RECIPES_REQUEST",
    FETCH_RECIPES_SUCCESS = "@@recipes/FETCH_RECIPES_SUCCESS",
    FETCH_RECIPES_FAILURE = "@@recipes/FETCH_RECIPES_FAILURE",
    FETCH_RECIPE_REQUEST = "@@recipes/FETCH_RECIPE_REQUEST",
    FETCH_RECIPE_SUCCESS = "@@recipes/FETCH_RECIPE_SUCCESS",
    FETCH_RECIPE_FAILURE = "@@recipes/FETCH_RECIPE_FAILURE",
}
