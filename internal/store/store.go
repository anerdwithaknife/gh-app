package store

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

func NewDefaultStore() StoreInterface {
	return NewYAMLStore("./store.yaml")
}
