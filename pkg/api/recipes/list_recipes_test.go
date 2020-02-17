package recipes_test

import (
	"errors"
	"github.com/iplay88keys/recipe-box/pkg/api/recipes"
	. "github.com/iplay88keys/recipe-box/pkg/helpers"
	"github.com/iplay88keys/recipe-box/pkg/repositories"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListRecipes", func() {
	It("returns the list of recipes", func() {
		listRecipes := func() ([]*repositories.Recipe, error) {
			return []*repositories.Recipe{{
				ID:          Int64Pointer(1),
				Name:        StringPointer("First"),
				Description: StringPointer("One"),
				Creator:     StringPointer("Some Creator"),
				Source:      StringPointer("Some Website"),
			}, {
				ID:          Int64Pointer(2),
				Name:        StringPointer("Second"),
				Description: StringPointer("Two"),
				Creator:     StringPointer("Another Creator"),
				Source:      nil,
			}}, nil
		}

		req := httptest.NewRequest("GET", "/recipes", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(recipes.ListRecipes(listRecipes).Handler)

		handler.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(rr.Body.String()).To(MatchJSON(`{
			"recipes": [{
				"id": 1,
				"name": "First",
				"description": "One",
				"creator": "Some Creator",
                "source": "Some Website"
			}, {
				"id": 2,
				"name": "Second",
				"description": "Two",
				"creator": "Another Creator"
			}]
		}`))
	})

	It("returns an error if the repository call fails", func() {
		listRecipes := func() ([]*repositories.Recipe, error) {
			return nil, errors.New("some error")
		}

		req := httptest.NewRequest("GET", "/recipes", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(recipes.ListRecipes(listRecipes).Handler)

		handler.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusInternalServerError))
	})
})
