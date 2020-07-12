import { Button, Link, Typography } from "@material-ui/core";
import Enzyme, { mount } from "enzyme";
import Adapter from "enzyme-adapter-react-16";
import React from "react";
import { Navigation } from "./Navigation";

Enzyme.configure({adapter: new Adapter()});

describe("Navigation", () => {
    it("should render a list of links", () => {
        const enzymeWrapper = mount(<Navigation/>);

        expect(enzymeWrapper.find(Typography).at(0).text()).toEqual("Recipe Box");
        expect(enzymeWrapper.find(Link)).toHaveLength(2);
        expect(enzymeWrapper.find(Link).at(0).text()).toEqual("Home");
        expect(enzymeWrapper.find(Link).at(1).text()).toEqual("Recipes");
        expect(enzymeWrapper.find(Button)).toHaveLength(1);
        expect(enzymeWrapper.find(Button).at(0).text()).toEqual("Register");
    });
});
