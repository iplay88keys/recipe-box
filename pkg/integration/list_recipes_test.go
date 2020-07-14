package integration_test

import (
    "encoding/json"
    "fmt"
    "io/ioutil"

    "github.com/iplay88keys/recipe-box/pkg/api/recipes"
    . "github.com/iplay88keys/recipe-box/pkg/helpers"
    "github.com/iplay88keys/recipe-box/pkg/repositories"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("ListRecipes", func() {
    var (
        username string
        password string
        firstRecipeID int64
        secondRecipeID int64
    )

    BeforeEach(func() {
        _, err := db.Exec("DELETE FROM users WHERE id IS NOT NULL")
        Expect(err).ToNot(HaveOccurred())

        username = "list_recipes_user"
        password = "Pa3$word123"

        usersRepo := repositories.NewUsersRepository(db)
        userID, err := usersRepo.Insert(username, username+"@example.com", password)
        Expect(err).ToNot(HaveOccurred())

        res, err := db.Exec(`INSERT INTO recipes (
            creator, name, description, servings, prep_time, total_time
            ) VALUES (?, ?, ?, ?, ?, ?)`,
            userID, "Root Beer Float", "Delicious drink for a hot summer day.", 1, "5 m", "5 m")
        Expect(err).ToNot(HaveOccurred())

        firstRecipeID, err = res.LastInsertId()
        Expect(err).ToNot(HaveOccurred())

        userID, err = usersRepo.Insert("another_"+username, "another_"+username+"@example.com", password)
        Expect(err).ToNot(HaveOccurred())

        res, err = db.Exec(`INSERT INTO recipes (
            creator, name, description, servings, prep_time, cook_time, total_time
            ) VALUES (?, ?, ?, ?, ?, ?, ?)`,
            userID, "Nana's Beans", "Spruced up baked beans.", 8, "10 m", "1-2 hrs", "1-2 hrs")
        Expect(err).ToNot(HaveOccurred())

        secondRecipeID, err = res.LastInsertId()
        Expect(err).ToNot(HaveOccurred())
    })

    It("returns a list of recipes if authenticated", func() {
        resp, err := client.Get(fmt.Sprintf("http://localhost:%s/api/v1/recipes", port))
        Expect(err).ToNot(HaveOccurred())

        // Needs auth header

        defer resp.Body.Close()
        bytes, err := ioutil.ReadAll(resp.Body)
        Expect(err).ToNot(HaveOccurred())

        var recipeList recipes.RecipeListResponse
        err = json.Unmarshal(bytes, &recipeList)
        Expect(err).ToNot(HaveOccurred())

        Expect(recipeList).To(Equal(recipes.RecipeListResponse{
            Recipes: []*repositories.Recipe{{
                ID:          Int64Pointer(firstRecipeID),
                Name:        StringPointer("Root Beer Float"),
                Description: StringPointer("Delicious drink for a hot summer day."),
                Creator:     nil,
                PrepTime:    nil,
                Source:      nil,
            }, {
                ID:          Int64Pointer(secondRecipeID),
                Name:        StringPointer("Nana's Beans"),
                Description: StringPointer("Spruced up baked beans."),
                Creator:     nil,
                PrepTime:    nil,
                Source:      nil,
            }},
        }))
    })
})
