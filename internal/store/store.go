package store

import (
	"fmt"
	"os"
)

type App struct {
	Name         string
	Slug         string
	AppID        int
	ClientID     string
	ClientSecret string
	PrivateKey   string
}

type StoreInterface interface {
	Init() error
	GetAppBySlug(slug string) (*App, error)
	SaveApp(app *App) error
	GetAllApps() ([]*App, error)
	DeleteApp(slug string) error
	Close() error
}

func getStorePath() string {
	storePath := os.Getenv("GH_APP_STORE_PATH")
	if storePath == "" {
		storePath = os.Getenv("HOME") + "/.gh-app.yaml"
	}
	return storePath
}

func NewDefaultStore(createIfMissing bool) (StoreInterface, error) {
	storePath := getStorePath()

	if _, err := os.Stat(storePath); os.IsNotExist(err) && !createIfMissing {
		return nil, fmt.Errorf("store file does not exist: %s", storePath)
	}

	db := NewYAMLStore(storePath)
	if err := db.Init(); err != nil {
		return nil, fmt.Errorf("error initializing store: %w", err)
	}

	return db, nil
}
