package store

import (
	"fmt"
	"os"
)

type App struct {
	Name       string
	Slug       string
	AppID      int
	ClientID   string
	PrivateKey string
}

type StoreInterface interface {
	Init() error
	GetAppBySlug(slug string) (*App, error)
	SaveApp(app *App) error
	GetAllApps() ([]*App, error)
	DeleteApp(slug string) error
	Close() error
}

func NewDefaultStore() (StoreInterface, error) {
	storePath := os.Getenv("APP_STORE_PATH")

	if storePath == "" {
		storePath = os.Getenv("HOME") + "/.gh-app.yaml"
	}

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("store path error: %w", err)
	}

	db := NewYAMLStore(storePath)
	if err := db.Init(); err != nil {
		return nil, fmt.Errorf("error initializing store: %w", err)
	}

	return db, nil
}
