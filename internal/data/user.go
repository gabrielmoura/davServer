package data

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"log"
	"os"
)

// User represents a user with a username and password.
type User struct {
	Username string
	Password string
}

// ResponseMap is a generic map for JSON responses.
type ResponseMap map[string]interface{}

// GetValidUsers returns the list of valid users.
func GetValidUsers() []User {
	users, _ := readUsers(dB)
	return users
}

// CreateUser adds a new user to the list.
func CreateUser(user User) []User {
	users, _ := readUsers(dB)
	users = append(users, user)
	if err := writeUsers(dB, users); err != nil {
		fmt.Println(err)
	}
	return users
}

// DeleteUser removes a user from the list.
func DeleteUser(username string) []User {
	users, _ := readUsers(dB)
	for i, u := range users {
		if u.Username == username {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
	if err := writeUsers(dB, users); err != nil {
		fmt.Println(err)
	}
	return users
}

// GetUser retrieves a user by username.
func GetUser(username string) (User, error) {
	users, _ := readUsers(dB)
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
		log.Printf("Erro ao criar pasta do usuário: %v\n", err)
		log.Printf("Pasta: %s\n", path)
		return "", err
	}
	return path, nil
}

func writeUsers(db *badger.DB, user []User) error {
	return db.Update(func(txn *badger.Txn) error {
		data, err := serializeUser(user)
		if err != nil {
			return err
		}
		return txn.Set([]byte("users"), data)
	})
}
func readUsers(db *badger.DB) ([]User, error) {
	var user []User
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("users"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			user, err = deserializeUser(val)
			return err
		})
	})
	return user, err
}
func serializeUser(user []User) ([]byte, error) {
	return json.Marshal(user)
}
func deserializeUser(data []byte) ([]User, error) {
	var user []User
	err := json.Unmarshal(data, &user)
	return user, err
}
