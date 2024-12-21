package build

import (
	"github.com/eryk-vieira/mango/internal/types"
)

type build struct {
	Settings *types.Settings
}

func NewBuilder(settings *types.Settings) *build {
	return &build{
		Settings: settings,
	}
}

func (b *build) Build() ([]Route, []Errors) {
	builder := routerBuilder{
		Settings: b.Settings,
	}

	routes, errorList := builder.Build()

	serverBuilder := serverBuilder{
		Settings:  b.Settings,
		debugMode: true,
	}

	serverBuilder.Build(routes, &errorList)

	if len(errorList) > 0 {
		return []Route{}, errorList
	}

	return routes, errorList
}
