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

func NewTekojarApp() (*TekojarApp, error) {
	s, err := tekojar.LoadSetting("./settings.json")
	if err != nil {
		return nil, err
	}

	tj, err := tekojar.NewWithSetting(s)
	if err != nil {
		return nil, err
	}
	return &TekojarApp{tekojar: tj}, nil
}

func (ta *TekojarApp) Startup(ctx context.Context) {
	ta.ctx = ctx
}

func (ta *TekojarApp) Shutdown(ctx context.Context) {
	ta.tekojar.StopAll()
}

func (ta *TekojarApp) GetSetting() *tekojar.TekojarSetting {
	return ta.tekojar.Setting.Current()
}

func (ta *TekojarApp) SaveSetting(ts tekojar.TekojarSetting) error {
	return ta.tekojar.Setting.Save(ts)
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

func (ta *TekojarApp) Get(name string) (ServiceView, error) {
	s, err := ta.tekojar.GetService(name)
	if err != nil {
		return ServiceView{}, err
	}
	return ServiceView{
		Name:   s.Name,
		Status: string(s.Status),
		Logs:   []LogView{},
	}, nil
}

func (ta *TekojarApp) Start(name string) error {
	ta.tekojar.Start(name)

	ch, err := ta.tekojar.WatchService(name)
	if err != nil {
		return err
	}

	go func() {
		for log := range ch {
			runtime.EventsEmit(ta.ctx, "service:log", map[string]interface{}{
				"name": name,
				"logView": LogView{
					IsError: ta.containsIgnoreCase(log, "error"),
					Log:     log,
				},
			})
		}
	}()

	return nil
}

func (ta *TekojarApp) containsIgnoreCase(str string, char string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(char))
}

func (ta *TekojarApp) Stop(name string) {
	ta.tekojar.Stop(name)
}
