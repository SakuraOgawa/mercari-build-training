package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	// STEP 5-1: uncomment this line
	// _ "github.com/mattn/go-sqlite3"
)

var errImageNotFound = errors.New("image not found")
var errItemNotFound = errors.New("item not found")

type Item struct {
	ID       int    `db:"id" json:"-"`
	Name     string `db:"name" json:"name"`
	Category string `db:"category" json:"category"`
	Image    string `db:"image" json:"image"`
}

// Please run `go generate ./...` to generate the mock implementation
// ItemRepository is an interface to manage items.
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=${GOPACKAGE} -destination=./mock_$GOFILE
type ItemRepository interface {
	Insert(ctx context.Context, item *Item) error
	List(ctx context.Context) ([]*Item, error)
	Select(ctx context.Context, id int) ([]*Item, error)
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	// fileName is the path to the JSON file storing items.
	fileName string
}

// NewItemRepository creates a new itemRepository.
func NewItemRepository() ItemRepository {
	return &itemRepository{fileName: "items.json"}
}

// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	// STEP 4-1: add an implementation to store an item

	var data struct {
		Items []*Item `json:"items"`
	}

	oldData, err := os.ReadFile(i.fileName)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(oldData) > 0 {
		if err := json.Unmarshal(oldData, &data); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	data.Items = append(data.Items, item)
	fmt.Println("Updated data:", data)

	newData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(i.fileName, newData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println("Item successfully added")
	return nil
}

// StoreImage stores an image and returns an error if any.
// This package doesn't have a related interface for simplicity.
func StoreImage(fileName string, image []byte) error {
	// STEP 4-4: add an implementation to store an image
	//画像を保存する
	if err := os.WriteFile(fileName, image, 0644); err != nil {
		return fmt.Errorf("failed to image file: %w", err)
	}

	return nil
}

func (i *itemRepository) List(ctx context.Context) ([]*Item, error) {
	fileData, err := os.ReadFile(i.fileName)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var jsonData struct {
		Items []*Item `json:"items"`
	}

	if err := json.Unmarshal(fileData, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmashal JSON: %w", err)
	}
	return jsonData.Items, nil
}

func (i *itemRepository) Select(ctx context.Context, id int) ([]*Item, error) {
	if id <= 0 {
		return nil, errItemNotFound
	}

	items, err := i.List(ctx)
	if err != nil {
		return nil, err
	}

	if len(items) < id {
		return nil, errItemNotFound
	}

	return []*Item{items[id-1]}, nil

}
