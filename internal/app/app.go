package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"collector/internal"
	"collector/internal/controller"
	"collector/internal/service"
)

type App struct {
	port    string
	router  *mux.Router
	service service.CollectorInterface
}

func NewApp(conf internal.Config) *App {
	return &App{
		port:    conf.PortService,
		router:  mux.NewRouter(),
		service: service.New(conf.SmsPath, conf.MmsUrl, conf.ViceCallPath, conf.EmailPath, conf.BillingPath, conf.IncidentUrl, conf.SupportUrl),
	}
}

func (a *App) Run() {
	a.router.HandleFunc("/", controller.HandleConnection(a.service))
	err := http.ListenAndServe(a.port, a.router)
	if err != nil {
		fmt.Printf("can't start http service: %s\n", err.Error())
	}
}
