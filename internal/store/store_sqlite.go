//go:build cgo
// +build cgo

package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore() *SqliteStore {
	return &SqliteStore{}
}

func (s *SqliteStore) Init() error {
	var err error
	s.db, err = sql.Open("sqlite3", "./store.db")
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS apps (
		name TEXT,
		slug TEXT NOT NULL PRIMARY KEY,
		app_id INTEGER NOT NULL,
		client_id TEXT,
		private_key TEXT
	);
	`)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStore) GetAppBySlug(slug string) (*App, error) {
	app := &App{}
	err := s.db.QueryRow("SELECT name, slug, app_id, client_id, private_key FROM apps WHERE slug = ?", slug).Scan(&app.Name, &app.Slug, &app.AppID, &app.ClientID, &app.PrivateKey)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (s *SqliteStore) SaveApp(app *App) error {
	_, err :=
		s.db.Exec("INSERT OR REPLACE INTO apps (name, slug, app_id, client_id, private_key) VALUES (?, ?, ?, ?, ?)",
			app.Name, app.Slug, app.AppID, app.ClientID, app.PrivateKey)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteStore) GetAllApps() ([]*App, error) {
	rows, err := s.db.Query("SELECT name, slug, app_id, client_id, private_key FROM apps")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*App
	for rows.Next() {
		app := &App{}
		err := rows.Scan(&app.Name, &app.Slug, &app.AppID, &app.ClientID, &app.PrivateKey)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (s *SqliteStore) DeleteApp(slug string) error {
	_, err := s.db.Exec("DELETE FROM apps WHERE slug = ?", slug)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
