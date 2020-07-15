package recipes_test

import (
    "database/sql"
    "errors"
    "net/http"
    "net/http/httptest"

    "golang.org/x/net/context"

    "github.com/iplay88keys/my-recipe-library/pkg/api/auth"

    "github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
    . "github.com/iplay88keys/my-recipe-library/pkg/helpers"
    "github.com/iplay88keys/my-recipe-library/pkg/repositories"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("ListRecipes", func() {
    It("returns the list of recipes", func() {
        listRecipes := func(userID int64) ([]*repositories.Recipe, error) {
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
                Creator:     StringPointer("Some Creator"),
                Source:      nil,
            }}, nil
        }

        req := httptest.NewRequest("GET", "/recipes", nil)
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

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
                "creator": "Some Creator"
            }]
        }`))
    })

    It("returns no content if there are no recipes", func() {
        listRecipes := func(userID int64) ([]*repositories.Recipe, error) {
            return nil, sql.ErrNoRows
        }

        req := httptest.NewRequest("GET", "/recipes", nil)
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.ListRecipes(listRecipes).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusNoContent))
    })

    It("returns an error if the repository call fails", func() {
        listRecipes := func(userID int64) ([]*repositories.Recipe, error) {
            return nil, errors.New("some error")
        }

        req := httptest.NewRequest("GET", "/recipes", nil)
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.ListRecipes(listRecipes).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
    })

    It("returns an error if the userID does not exist on the context", func() {
        listRecipes := func(userID int64) ([]*repositories.Recipe, error) {
            return nil, nil
        }

        req := httptest.NewRequest("GET", "/recipes", nil)

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.ListRecipes(listRecipes).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
    })
})
