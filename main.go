package main

import (
	"github.com/SpeedVan/go-common/app/web"
	"github.com/SpeedVan/go-common/config/env"
	"github.com/SpeedVan/go-common/log"
	"github.com/SpeedVan/http-proxy/controller"
)

func main() {

	if cfg, err := env.LoadAllWithoutPrefix("CALLBACK_"); err == nil {
		logger := log.NewCommon(log.Debug)

		app := web.New(cfg, logger)

		// pc := controller.NewPC(httpClient, cfg.Get("SVC"))

		// esc := controller.NewESC(cfg, table)

		app.HandleController(controller.New(cfg))
		// app.Router.Use(mux.CORSMethodMiddleware(app.Router))
		app.Run(log.Debug)
	}
}
