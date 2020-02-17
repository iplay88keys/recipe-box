import { createLocation, createMemoryHistory, Location, MemoryHistory } from "history";
import React from "react";
import Enzyme, { shallow } from "enzyme";
import Adapter from "enzyme-adapter-react-16";
import Table from "react-bootstrap/Table";
import { match } from "react-router";
import { RecipeResponse } from "../../state/ducks/recipes/types";
import { RecipeList } from "./RecipeList";

Enzyme.configure({adapter: new Adapter()});

describe("RecipeList", () => {
    let history: MemoryHistory;
    let matchParam: match<{ id: string }>;
    let location: Location;

    beforeEach(() => {
        history = createMemoryHistory();
        const path = `/route/:id`;

        matchParam = {
            isExact: false,
            path,
            url: path.replace(":id", "1"),
            params: {id: "1"}
        } as match<{ id: string }>;

        location = createLocation(matchParam.url);
    });

    it("should render a list of recipes", () => {
        const recipes = [{
            id: 0,
            name: "First",
            description: "One"
        }, {
            id: 1,
            name: "Second",
            description: "Two"
        }] as RecipeResponse[];

        const enzymeWrapper = shallow(
            <RecipeList
                recipes={recipes}
                loading={false}
                history={history}
                match={matchParam}
                location={location}
            />
        );

        expect(enzymeWrapper.find(Table)).toHaveLength(1);
        expect(enzymeWrapper.find("td")).toHaveLength(4);
        expect(enzymeWrapper.find("tbody").children()).toHaveLength(2);
        expect(enzymeWrapper.find("tbody").childAt(0).childAt(0).text()).toEqual("First");
        expect(enzymeWrapper.find("tbody").childAt(0).childAt(1).text()).toEqual("One");
        expect(enzymeWrapper.find("tbody").childAt(1).childAt(0).text()).toEqual("Second");
        expect(enzymeWrapper.find("tbody").childAt(1).childAt(1).text()).toEqual("Two");
    });

    it("does not render missing data", () => {
        const recipes = [{
            id: 0,
            name: "First"
        }] as RecipeResponse[];

        const enzymeWrapper = shallow(
            <RecipeList
                recipes={recipes}
                loading={false}
                history={history}
                match={matchParam}
                location={location}
            />
        );

        expect(enzymeWrapper.find(Table)).toHaveLength(1);
        expect(enzymeWrapper.find("td")).toHaveLength(1);
        expect(enzymeWrapper.find("tbody").children()).toHaveLength(1);
        expect(enzymeWrapper.find("tbody").childAt(0).children()).toHaveLength(1);
        expect(enzymeWrapper.find("tbody").childAt(0).childAt(0).text()).toEqual("First");
    });

    it("should load the single recipe page when the row is clicked", () => {
        const recipes = [{
            id: 0,
            name: "First",
            description: "One"
        }, {
            id: 1,
            name: "Second",
            description: "Two"
        }] as RecipeResponse[];

        const historyMock = {
            length: {} as any,
            action: {} as any,
            location: {} as any,
            push: jest.fn(),
            replace: jest.fn(),
            go: jest.fn(),
            goBack: jest.fn(),
            goForward: jest.fn(),
            block: jest.fn(),
            listen: jest.fn(),
            createHref: jest.fn()
        };

        const enzymeWrapper = shallow(
            <RecipeList
                recipes={recipes}
                loading={false}
                history={historyMock}
                match={matchParam}
                location={location}
            />
        );

        enzymeWrapper.find("tbody").childAt(0).simulate("click");
        expect(historyMock.push.mock.calls[0]).toEqual(["/recipes/0"]);
    });

    it("should render loading info when loading", () => {
        const props = {
            recipes: [] as RecipeResponse[]
        };

        const enzymeWrapper = shallow(
            <RecipeList
                recipes={props.recipes}
                loading={true}
                history={history}
                match={matchParam}
                location={location}
            />
        );

        expect(enzymeWrapper.find("div")).toHaveLength(1);
        expect(enzymeWrapper.find("div p").text()).toBe("Loading recipes");
    });
});
