package carddav

import (
	"context"
	"fmt"
	"github.com/emersion/go-vcard"
	"github.com/emersion/go-webdav/carddav"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"log"
	"path/filepath"
)

type CardBackend struct {
	addressBooks []carddav.AddressBook
	objectMap    map[string][]carddav.AddressObject
}

func (c *CardBackend) AddressbookHomeSetPath(ctx context.Context) (string, error) {
	user, ok := ctx.Value("user").(data.User)
	if !ok {
		return "", fmt.Errorf("Usuário não encontrado")
	}
	fullPath := filepath.Join(config.Conf.ShareRootDir, user.Username, "address")
	return fullPath, nil
}
func (c *CardBackend) AddressBook(ctx context.Context) (*carddav.AddressBook, error) {
	return nil, nil
}
func (c *CardBackend) GetAddressObject(ctx context.Context, path string, req *carddav.AddressDataRequest) (*carddav.AddressObject, error) {
	for _, objs := range c.objectMap {
		for _, obj := range objs {
			if obj.Path == path {
				return &obj, nil
			}
		}
	}
	log.Printf("Couldn't find address object at: %s", path)
	return nil, nil
}
func (c *CardBackend) ListAddressObjects(ctx context.Context, req *carddav.AddressDataRequest) ([]carddav.AddressObject, error) {
	user, ok := ctx.Value("user").(data.User)
	if !ok {
		return nil, fmt.Errorf("Usuário não encontrado")
	}
	fullPath := filepath.Join(config.Conf.ShareRootDir, user.Username, "address")
	addressObject, err := c.GetAddressObject(ctx, fullPath, req)
	if err != nil {
		return nil, err
	}

	return []carddav.AddressObject{*addressObject}, nil
}
func (c *CardBackend) QueryAddressObjects(ctx context.Context, query *carddav.AddressBookQuery) ([]carddav.AddressObject, error) {
	return nil, nil
}
func (c *CardBackend) PutAddressObject(ctx context.Context, path string, card vcard.Card, opts *carddav.PutAddressObjectOptions) (loc string, err error) {
	return "", nil
}
func (c *CardBackend) DeleteAddressObject(ctx context.Context, path string) error {
	return nil
}
func (c *CardBackend) CurrentUserPrincipal(ctx context.Context) (string, error) {
	user, ok := ctx.Value("user").(data.User)
	if !ok {
		return "", fmt.Errorf("Usuário não encontrado")
	}
	fullPath := filepath.Join(config.Conf.ShareRootDir, user.Username, "/")
	return fullPath, nil
}
