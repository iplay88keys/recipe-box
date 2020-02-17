import React from "react";
import Enzyme, { mount } from "enzyme";
import Adapter from "enzyme-adapter-react-16";
import Nav from "react-bootstrap/Nav";
import { Navigation } from "./Navigation";

Enzyme.configure({adapter: new Adapter()});

describe("Navigation", () => {
    it("should render a list of links", () => {
        const enzymeWrapper = mount(<Navigation/>);

        expect(enzymeWrapper.find("div").find(Nav.Link)).toHaveLength(2);
    });
});
