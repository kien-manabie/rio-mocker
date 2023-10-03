package setup

import (
	"context"

	"github.com/kien-manabie/rio-mocker"
	"github.com/kien-manabie/rio-mocker/internal/cache"
	"github.com/kien-manabie/rio-mocker/internal/config"
	"github.com/kien-manabie/rio-mocker/internal/database"
	fs "github.com/kien-manabie/rio-mocker/internal/storage"
)

func ProvideStubStore(ctx context.Context, cfg *config.Config) (rio.StubStore, error) {
	db, err := database.NewStubDBStore(ctx, cfg.DB)
	if err != nil {
		return nil, err
	}

	return cache.NewStubCache(db, db, cfg), nil
}

func ProvideFileStorage(ctx context.Context, cfg *config.Config) (fs.FileStorage, error) {
	return fs.NewFileStorage(cfg.FileStorage), nil
}
