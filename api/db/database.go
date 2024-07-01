package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"git.alnav.dev/alnav3/splitweb/api/encryption"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
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
    return db, nil
}

func GetSessionToken() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    return os.Getenv("TURSO_JWT_SECRET")
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

    fmt.Println("password encriptada: ", encryptedPass)
    secret := os.Getenv("TURSO_JWT_SECRET")
    jwt := encryption.GenerateJWT(username, secret)
    // authenticate the user
    return jwt , nil
}

func SignupUser(username string, hashedPassword string, db *sql.DB) error {
    result, err := db.Exec("INSERT INTO user (user, encryptedPassword) VALUES (?, ?)", username, hashedPassword)
    if err != nil {
        if strings.Contains(err.Error(), "UNIQUE constraint failed: user.user") {
            fmt.Println("Error: User already exists")
            return errors.New("user already exists")
        }
        fmt.Println("Error inserting user: ", err)
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        fmt.Println("Error getting rows affected: ", err)
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
