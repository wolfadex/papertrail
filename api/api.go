package api

import (
	"path"

	"github.com/efritz/nacelle"

	"github.com/aphistic/papertrail/storage"
	"github.com/aphistic/papertrail/storage/sqlite"
)

type Initializer struct {
	Container *nacelle.ServiceContainer `service:"container"`
	Logger    nacelle.Logger            `service:"logger"`
}

func NewInitializer() *Initializer {
	return &Initializer{}
}

func (i *Initializer) Init(config nacelle.Config) error {
	rawConfig, err := config.Get(ConfigToken)
	if err != nil {
		return err
	}
	cfg := rawConfig.(*Config)

	fs, err := storage.NewFileLocal(cfg.StorageRoot)
	if err != nil {
		return err
	}

	ds, err := sqlite.NewClient(path.Join(cfg.StorageRoot, "papertrail.db"))
	if err != nil {
		return err
	}
	err = ds.Migrate()
	if err != nil {
		return err
	}

	err = i.Container.Set("api", &Client{
		cfg:    cfg,
		logger: i.Logger,

		fileStorage: fs,
		dataStorage: ds,
	})
	if err != nil {
		return err
	}

	return nil
}

type Client struct {
	cfg    *Config
	logger nacelle.Logger

	fileStorage storage.File
	dataStorage storage.Data
}
