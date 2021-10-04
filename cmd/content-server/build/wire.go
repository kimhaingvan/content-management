// +build wireinject

package build

import (
	"cms-project/cmd/content-server/config"
	"cms-project/pkg/application"
	"cms-project/pkg/database"

	"github.com/google/wire"
)

func InitApp(cfg config.Config) (*App, error) {
	wire.Build(
		database.WireSet,
		//storage.WireSet,
		//application.WireSet,
		buildApplication,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}

type App struct {
	Db  *database.Database
	App *application.Application
}
