package backup

import (
	"encoding/json"
	"fmt"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/go/pkg/ternary"
	"os"
	"path/filepath"
)

func HandleBackup() error {
	if *config.ImportUsers || *config.ExportUsers {
		args := os.Args[2:]
		if *config.ImportUsers {
			return inputUser(args)
		}
		if *config.ExportUsers {
			return exportUser(args)
		}
		panic("Erro ao realizar backup")
	}
	return nil
}

func inputUser(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("nenhum arquivo de entrada fornecido")
	}
	path, err := filepath.Abs(ternary.OrString(args[0], "users.json"))
	if err != nil {
		return err
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var userList []data.User
	if err := json.Unmarshal(file, &userList); err != nil {
		return err
	}
	for _, user := range userList {
		data.CreateUser(user)
	}
	fmt.Println("Usuários importados com sucesso")
	os.Exit(0)
	return nil
}

func exportUser(args []string) error {
	path := filepath.Join(".", ternary.OrString(args[0], "users.json"))
	users := data.GetValidUsers()
	file, err := json.Marshal(users)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, file, 0644); err != nil {
		return err
	}
	fmt.Println("Usuários exportados com sucesso")
	os.Exit(0)
	return nil
}
