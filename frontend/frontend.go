package frontend

import (
	"fmt"
	"net/http"

	"github.com/gernest/utron"
	"github.com/svarogg/dedagger/backend"
	"github.com/svarogg/dedagger/frontend/controllers"
)

type Frontend struct {
	backend *backend.Backend
}

func NewFrontend(be *backend.Backend) *Frontend {
	return &Frontend{backend: be}
}

func (fe *Frontend) Start() error {
	app, err := utron.NewMVC("frontend/config")
	if err != nil {
		return err
	}

	app.AddController(controllers.NewStores(fe.backend))

	fmt.Printf("Listening on %s\n", app.Config.BaseURL)
	return http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), app)
}
