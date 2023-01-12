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
	port    string
	router  *mux.Router
	service service.CollectorInterface
}

func NewApp(conf pkg.Config) *App {
	return &App{
		port:    conf.UrlService,
		router:  mux.NewRouter(),
		service: service.New(conf.SmsPath, conf.MmsUrl, conf.VoiceCallPath, conf.EmailPath, conf.BillingPath, conf.IncidentUrl, conf.SupportUrl),
	}
}

func (a *App) Run() {
	a.router.HandleFunc("/", controller.HandleConnection(a.service))
	err := http.ListenAndServe(a.port, a.router)
	if err != nil {
		fmt.Printf("can't start http service: %s\n", err.Error())
	}
}
