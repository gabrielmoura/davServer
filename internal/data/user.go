package data

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
)

var users []User

// User represents a user with a username and password.
type User struct {
	Username string
	Password string
}

// ResponseMap is a generic map for JSON responses.
type ResponseMap map[string]interface{}

// GetValidUsers returns the list of valid users.
func GetValidUsers() []User {
	return users
}

// CreateUser adds a new user to the list.
func CreateUser(user User) []User {
	users = append(users, user)
	return users
}

// DeleteUser removes a user from the list.
func DeleteUser(username string) []User {
	for i, u := range users {
		if u.Username == username {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
	return users
}

// GetUser retrieves a user by username.
func GetUser(username string) (User, error) {
	for _, u := range users {
		if u.Username == username {
			return u, nil
		}
	}
	return User{}, fmt.Errorf("usuário não encontrado")
}

// GenerateMD5Hash generates an MD5 hash of a given password.
func GenerateMD5Hash(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// CreateUserDirectory creates a directory for a user.
func CreateUserDirectory(path string) (string, error) {
	if err := os.Mkdir(path, 0755); err != nil {
		return "", err
	}
	return path, nil
}
