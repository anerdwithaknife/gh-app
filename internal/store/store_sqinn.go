//go:build !cgo
// +build !cgo

package store

import (
	"github.com/cvilsmeier/sqinn-go/sqinn"
)

type sqinnStore struct {
	db *sqinn.Sqinn
}

func NewSqinnStore() StoreInterface {
	return &sqinnStore{}
}

func (s *sqinnStore) Init() error {
	var err error
	s.db, err = sqinn.Launch(sqinn.Options{})
	if err != nil {
		return err
	}

	err = s.db.Open("./store.db")
	if err != nil {
		return err
	}

	_, err = s.db.ExecOne(`
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

func (s *sqinnStore) GetAppBySlug(slug string) (*App, error) {
	rows, err := s.db.Query(
		"SELECT name, slug, app_id, client_id, private_key FROM apps WHERE slug = ?",
		[]any{slug},
		[]byte{sqinn.ValText, sqinn.ValText, sqinn.ValInt, sqinn.ValText, sqinn.ValText},
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil // not found
	}
	row := rows[0]
	app := &App{
		Name:       row.Values[0].AsString(),
		Slug:       row.Values[1].AsString(),
		AppID:      row.Values[2].AsInt(),
		ClientID:   row.Values[3].AsString(),
		PrivateKey: row.Values[4].AsString(),
	}
	return app, nil
}

func (s *sqinnStore) SaveApp(app *App) error {
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO apps (name, slug, app_id, client_id, private_key) VALUES (?, ?, ?, ?, ?)",
		1, // niterations
		5, // nparams
		[]any{app.Name, app.Slug, app.AppID, app.ClientID, app.PrivateKey},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqinnStore) GetAllApps() ([]*App, error) {
	rows, err := s.db.Query(
		"SELECT name, slug, app_id, client_id, private_key FROM apps",
		nil,
		[]byte{sqinn.ValText, sqinn.ValText, sqinn.ValInt, sqinn.ValText, sqinn.ValText},
	)
	if err != nil {
		return nil, err
	}
	apps := make([]*App, 0, len(rows))
	for _, row := range rows {
		app := &App{
			Name:       row.Values[0].AsString(),
			Slug:       row.Values[1].AsString(),
			AppID:      row.Values[2].AsInt(),
			ClientID:   row.Values[3].AsString(),
			PrivateKey: row.Values[4].AsString(),
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (s *sqinnStore) DeleteApp(slug string) error {
	_, err := s.db.Exec(
		"DELETE FROM apps WHERE slug = ?",
		1, // niterations
		1, // nparams
		[]any{slug},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqinnStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
