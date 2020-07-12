package repositories

import (
    "database/sql"
    "errors"
    "fmt"

    "golang.org/x/crypto/bcrypt"
)

type User struct {
    Username string
    Email    string
}

type UsersRepository struct {
    db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
    return &UsersRepository{db: db}
}

func (u *UsersRepository) Authenticate(username, email, password string) (bool, error) {

    return false, nil
}

func (u *UsersRepository) ExistsByUsername(username string) (bool, error) {
    if username == "" {
        fmt.Println("Could not check for user: username required")
        return false, errors.New("could not check for user: username required")
    }

    userRow := u.db.QueryRow(existsByUsernameQuery, username)

    var user User
    err := userRow.Scan(&user.Username)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }

        fmt.Printf("Failed to scan row for user by username '%s': %s\n", username, err.Error())
        return false, errors.New("failed to query for user by username")
    }

    return true, nil
}

func (u *UsersRepository) ExistsByEmail(email string) (bool, error) {
    if email == "" {
        fmt.Println("Could not check for user: username required")
        return false, errors.New("could not check for user: email required")
    }

    userRow := u.db.QueryRow(existsByEmailQuery, email)

    var user User
    err := userRow.Scan(&user.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }

        fmt.Printf("Failed to scan row for user by email '%s': %s\n", email, err.Error())
        return false, errors.New("failed to query for user by email")
    }

    return true, nil
}

func (u *UsersRepository) Insert(username, email, password string) (int64, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
        fmt.Printf("Could not hash the user password")
        return -1, errors.New("could not has password")
    }

    res, err := u.db.Exec(insertUserQuery,
        username,
        email,
        string(hashedPassword),
    )

    if err != nil {
        fmt.Printf("User could not be saved: %s\n", err.Error())
        return -1, errors.New("user could not be added")
    }

    id, err := res.LastInsertId()
    if err != nil {
        fmt.Printf("Recipe was not saved correctly: %s\n", err.Error())
        return 0, errors.New("user was not saved correctly")
    }

    return id, nil
}

const existsByUsernameQuery = "SELECT username FROM users WHERE username=?"
const existsByEmailQuery = "SELECT email FROM users WHERE email=?"
const insertUserQuery = `INSERT INTO users
(username,
email,
password_hash)
VALUES (?, ?, ?)
`
