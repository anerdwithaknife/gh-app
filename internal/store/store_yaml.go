package store

import (
	"errors"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type yamlStore struct {
	filePath string
	apps     map[string]*App
	mu       sync.Mutex
}

func NewYAMLStore(filePath string) *yamlStore {
	return &yamlStore{
		filePath: filePath,
		apps:     make(map[string]*App),
	}
}

func (s *yamlStore) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.apps = make(map[string]*App)
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var appList []*App
	if err := decoder.Decode(&appList); err != nil && err.Error() != "EOF" {
		return err
	}

	s.apps = make(map[string]*App)
	for _, app := range appList {
		s.apps[app.Slug] = app
	}
	return nil
}

func (s *yamlStore) persist() error {
	appList := make([]*App, 0, len(s.apps))
	for _, app := range s.apps {
		appList = append(appList, app)
	}
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	return encoder.Encode(appList)
}

func (s *yamlStore) GetAppBySlug(slug string) (*App, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	app, ok := s.apps[slug]
	if !ok {
		return nil, nil
	}
	return app, nil
}

func (s *yamlStore) SaveApp(app *App) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.apps[app.Slug] = app
	return s.persist()
}

func (s *yamlStore) GetAllApps() ([]*App, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	appList := make([]*App, 0, len(s.apps))
	for _, app := range s.apps {
		appList = append(appList, app)
	}
	return appList, nil
}

func (s *yamlStore) DeleteApp(slug string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.apps[slug]; !ok {
		return errors.New("app not found")
	}
	delete(s.apps, slug)
	return s.persist()
}

func (s *yamlStore) Close() error {
	// No resources to close for YAML store
	return nil
}
