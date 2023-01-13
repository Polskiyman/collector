package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"collector/internal/controller"
	"collector/internal/service"
	"collector/pkg/config"
)

type App struct {
	Url     string
	Router  *mux.Router
	Service service.CollectorInterface
}

func NewApp(conf config.Config) *App {
	return &App{
		Url:     conf.AppUrl,
		Router:  mux.NewRouter(),
		Service: service.New(conf.Adapters),
	}
}

func (a *App) Run() {
	a.Router.HandleFunc("/", controller.HandleConnection(a.Service))
	err := http.ListenAndServe(a.Url, a.Router)
	if err != nil {
		fmt.Printf("can't start http Service: %s\n", err.Error())
	}
}
