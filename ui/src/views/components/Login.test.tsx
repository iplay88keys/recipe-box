import { TextField } from "@material-ui/core";
import { act } from "@testing-library/react";
import Enzyme, { mount, ReactWrapper } from "enzyme";
import Adapter from "enzyme-adapter-react-16";
import React, { ChangeEvent, FocusEvent } from "react";
import { loginAsync } from "../../state/ducks/users/actions";
import Login, { LoginFormInner } from "./Login";

Enzyme.configure({adapter: new Adapter()});

describe("Login", () => {
    it("should render a form for logging in", () => {
        let login = jest.fn(loginAsync.request);

        const enzymeWrapper = mount(
            <Login
                login={login}
            />
        );

        expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(0).text()).toEqual("Username/Email Address *Username/Email Address *");
        expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(1).text()).toEqual("Password *Password *");
    });

    describe("form validation errors", () => {
        let enzymeWrapper: ReactWrapper;

        beforeEach(() => {
            let login = jest.fn(loginAsync.request);

            enzymeWrapper = mount(
                <Login
                    login={login}
                />
            );

        });

        describe("login", () => {
            beforeEach(async () => {
                await act(async () => {
                    enzymeWrapper.find(LoginFormInner).find(TextField).at(0).props().onBlur!({
                        preventDefault() {},
                        target: {
                            name: "login"
                        }
                    } as FocusEvent<HTMLInputElement>);

                    enzymeWrapper.update();
                });
            });

            it("is required", async () => {
                await act(async () => {
                    enzymeWrapper.find(LoginFormInner).find(TextField).at(0).props().onChange!({
                        preventDefault() {},
                        target: {
                            name: "login",
                            value: ""
                        }
                    } as ChangeEvent<HTMLInputElement>);
                });

                enzymeWrapper.update();
                expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(0).props().value).toEqual("");
                expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(0).props().error).toEqual(true);
                expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(0).props().helperText)
                    .toEqual("Required");
            });
        });

        describe("password", () => {
            beforeEach(async () => {
                await act(async () => {
                    enzymeWrapper.find(LoginFormInner).find(TextField).at(1).props().onBlur!({
                        preventDefault() {},
                        target: {
                            name: "password"
                        }
                    } as FocusEvent<HTMLInputElement>);

                    enzymeWrapper.update();
                });
            });

            it("is required", async () => {
                await act(async () => {
                    enzymeWrapper.find(LoginFormInner).find(TextField).at(1).props().onChange!({
                        preventDefault() {},
                        target: {
                            name: "password",
                            value: ""
                        }
                    } as ChangeEvent<HTMLInputElement>);
                });

                enzymeWrapper.update();
                expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(1).props().value).toEqual("");
                expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(1).props().error).toEqual(true);
                expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(1).props().helperText)
                    .toEqual("Required");
            });
        });
    });

    describe("api validation errors", () => {
        it("displays username errors", async () => {
            let login = jest.fn(loginAsync.request);

            let enzymeWrapper = mount(
                <Login
                    login={login}
                />
            );

            await act(async () => {
                enzymeWrapper.find(LoginFormInner).props().setStatus({"login": "Api Error"});
            });

            enzymeWrapper.update();
            expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(0).props().value)
                .toEqual("");
            expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(0).props().error).toEqual(true);
            expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(0).props().helperText)
                .toEqual("Api Error");
        });

        it("displays password errors", async () => {
            let login = jest.fn(loginAsync.request);

            let enzymeWrapper = mount(
                <Login
                    login={login}
                />
            );

            await act(async () => {
                enzymeWrapper.find(LoginFormInner).props().setStatus({"password": "Api Error"});
            });

            enzymeWrapper.update();
            expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(1).props().value)
                .toEqual("");
            expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(1).props().error).toEqual(true);
            expect(enzymeWrapper.find(LoginFormInner).find(TextField).at(1).props().helperText)
                .toEqual("Api Error");
        });
    });
});
