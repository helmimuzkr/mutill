package backend

import (
	"context"
	"strings"

	"tekotools/tekojar"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TekojarApp struct {
	ctx     context.Context
	tekojar *tekojar.Tekojar
}

type ServiceView struct {
	Name   string    `json:"name"`
	Status string    `json:"status"`
	Logs   []LogView `json:"logs"`
}

type LogView struct {
	IsError bool   `json:"is_error"`
	Log     string `json:"log"`
}

func NewTekojarApp() *TekojarApp {
	cfg, err := tekojar.LoadSetting("./settings.json")
	if err != nil {
		tekojar.PrintErr(tekojar.SYSTEM, 0, err.Error())
		panic(err)
	}

	tj := tekojar.New(cfg)

	return &TekojarApp{
		tekojar: tj,
	}
}

func (ta *TekojarApp) Startup(ctx context.Context) {
	ta.ctx = ctx
}

func (ta *TekojarApp) Shutdown(ctx context.Context) {
	ta.tekojar.StopAll()
}

func (ta *TekojarApp) GetAll() []ServiceView {
	services := ta.tekojar.GetAll()

	servicesView := make([]ServiceView, 0, len(services))
	for _, s := range services {
		servicesView = append(servicesView, ServiceView{
			Name:   s.Name,
			Status: string(s.Status),
			Logs:   []LogView{},
		})
	}

	return servicesView
}

func (ta *TekojarApp) Get(name string) ServiceView {
	s := ta.tekojar.GetService(name)
	return ServiceView{
		Name:   s.Name,
		Status: string(s.Status),
		Logs:   []LogView{},
	}
}

func (ta *TekojarApp) Start(name string) {
	ta.tekojar.Start(name)

	go func() {
		for log := range ta.tekojar.WatchService(name) {
			runtime.EventsEmit(ta.ctx, "service:log", map[string]interface{}{
				"name": name,
				"logView": LogView{
					IsError: ta.containsIgnoreCase(log, "error"),
					Log:     log,
				},
			})
		}
	}()
}

func (ta *TekojarApp) containsIgnoreCase(str string, char string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(char))
}

func (ta *TekojarApp) Stop(name string) {
	ta.tekojar.Stop(name)
}
