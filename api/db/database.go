package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
    "git.alnav.dev/alnav3/splitweb/api/encryption"
)

func Connect() (*sql.DB, error) {
    // connect to the database
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    token := os.Getenv("TURSO_AUTH_TOKEN")
    dbUrl := os.Getenv("TURSO_DATABASE_URL")
    url := dbUrl + "?authToken=" + token
    db, err := sql.Open("libsql", url)
    if err != nil {
        log.Panic(fmt.Sprintf("Error connecting to database %s: %s", dbUrl, err))
    }
    fmt.Println(db)
    return db, nil
}

func AuthUser(username, password string, db *sql.DB) (string, error) {
	var encryptedPass string
	err := db.QueryRow("SELECT encryptedPassword FROM user WHERE user = ?", username).Scan(&encryptedPass)

    if err != nil {
        return "", handleSqlErr(err)
    }

    err = encryption.Compare(password, encryptedPass)
    if err != nil {
        return "", errors.New("password is incorrect")
    }

    secret := os.Getenv("TURSO_JWT_SECRET")
    jwt := encryption.GenerateJWT(username, secret)
    // authenticate the user
    return jwt , nil
}

func SignupUser(username string, hashedPassword string, db *sql.DB) error {
    result, err := db.Exec("INSERT INTO user (user, encryptedPassword) VALUES (?, ?)", username, hashedPassword)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected != 1 {
        return fmt.Errorf("expected to affect 1 row, affected %d", rowsAffected)
    }

    return nil
}

func Validate(tokenString string, database *sql.DB) bool {
    secret := os.Getenv("TURSO_JWT_SECRET")
    username, err := encryption.ValidateToken(tokenString, secret)
    if err != nil {
        return false
    }
    return usernameExists(username, database)
}

func usernameExists(username string, database *sql.DB) bool {
    var user string
    err := database.QueryRow("SELECT user FROM user WHERE user = ?", username).Scan(&user)
    if err != nil {
        return false
    }
    return true
}

func handleSqlErr(err error) error {
    if err == sql.ErrNoRows {
        // No user with the given username was found
        return errors.New("user doesn't exist")
    }
    // A database error occurred
    return err
}
