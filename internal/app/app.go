package app

import (
	"collector/pkg"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"collector/internal/controller"
	"collector/internal/service"
)

type App struct {
	Url     string
	Router  *mux.Router
	Service service.CollectorInterface
}

func NewApp(conf pkg.Config) *App {
	return &App{
		Url:     conf.UrlService,
		Router:  mux.NewRouter(),
		Service: service.New(conf),
	}
}

func (a *App) Run() {
	a.Router.HandleFunc("/", controller.HandleConnection(a.Service))
	err := http.ListenAndServe(a.Url, a.Router)
	if err != nil {
		fmt.Printf("can't start http Service: %s\n", err.Error())
	}
}
