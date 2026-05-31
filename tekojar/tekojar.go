package tekojar

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

type Tekojar struct {
	Setting *Setting

	services map[string]*Service

	mu sync.RWMutex
}

func New() (*Tekojar, error) {
	// go run . --settings /path/to/setting.json
	settingPath := flag.String("settings", "./settings.json", "path to settings file")
	flag.Parse()

	s, err := LoadSetting(*settingPath)
	if err != nil {
		PrintErr(SYSTEM, 0, err.Error())
		return nil, err
	}

	return NewWithSetting(s)
}

func NewWithSetting(Setting *Setting) (*Tekojar, error) {
	t := &Tekojar{
		services: make(map[string]*Service),
		Setting:  Setting,
	}

	t.registerServices()

	return t, nil
}

func (t *Tekojar) registerServices() {
	t.mu.Lock()
	defer t.mu.Unlock()

	ts := t.Setting.Current()
	PrintLog(SYSTEM, 0, fmt.Sprintf("Tekojar Setting: %#v", ts))

	if ts.ServiceSettings == nil {
		PrintErr(SYSTEM, 0, "no service inputed")
	}

	if t.services == nil {
		t.services = make(map[string]*Service)
	}

	for _, s := range ts.ServiceSettings {
		if t.services[s.Name] != nil {
			continue
		}
		t.services[s.Name] = InitService(s.Name, s.Path, s.SkipFlag)
	}
}

func (t *Tekojar) GetService(name string) (*Service, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	s := t.services[name]
	if s == nil {
		PrintErr(SYSTEM, 0, "Service not found")
		return nil, errors.New("Service not found")
	}
	return t.services[name], nil
}

func (t *Tekojar) WatchService(name string) (chan string, error) {
	s, err := t.GetService(name)
	if err != nil {
		return nil, err
	}
	return s.Subscribe(), nil
}

func (t *Tekojar) UnwatchService(name string) error {
	s, err := t.GetService(name)
	if err != nil {
		return err
	}
	s.Unsubscribe()
	return nil
}

func (t *Tekojar) Start(name string) error {
	s, err := t.GetService(name)
	if err != nil {
		return err
	}
	ts := t.Setting.Current()
	go s.StartProcess(ts.Command)
	go t.ListenShutdown(name)
	return nil
}

func (t *Tekojar) Stop(name string) error {
	s, err := t.GetService(name)
	if err != nil {
		return err
	}
	s.StopProcess()
	return nil
}

func (t *Tekojar) GetAll() []*Service {
	t.registerServices()

	t.mu.Lock()
	services := make([]*Service, 0, len(t.services))
	for _, v := range t.services {
		services = append(services, v)
	}
	t.mu.Unlock()

	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})

	return services
}

func (t *Tekojar) StartAll() {
	PrintLog(SYSTEM, 0, "starting application...")

	t.mu.RLock()
	services := make([]*Service, 0, len(t.services))
	for _, v := range t.services {
		if v.GetStatus() == INACTIVE {
			services = append(services, v)
		}
	}
	t.mu.RUnlock()

	ts := t.Setting.Current()

	for _, s := range services {
		if s.isSkip {
			continue
		}
		go s.StartProcess(ts.Command)
	}

	t.ListenAndShutdownAll()
}

func (t *Tekojar) StopAll() {
	t.mu.RLock()
	services := make([]*Service, 0, len(t.services))
	for _, v := range t.services {
		services = append(services, v)
	}
	t.mu.RUnlock()

	for _, v := range services {
		v.StopProcess()
	}

	actives, inactives := t.GetTotalStatusServices()
	PrintLog(SYSTEM, 0, "application stopped", fmt.Sprintf("actives: %d", actives), fmt.Sprintf("inactives: %d", inactives))
}

func (t *Tekojar) GetTotalStatusServices() (int, int) {
	active := 0
	inactive := 0

	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, s := range t.services {
		switch s.GetStatus() {
		case ACTIVE:
			active++
		case INACTIVE:
			inactive++
		}
	}

	return active, inactive
}

// ListenShutdown will block goroutines
func (t *Tekojar) ListenAndShutdownAll() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	ts := t.Setting.Current()

	if ts.AutoShutdown {
		t.AutomaticShutDownTicker("", 5*time.Second, sigChan)
	}

	<-sigChan

	t.StopAll()
}

// ListenShutdown will block goroutines
func (t *Tekojar) ListenShutdown(name string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	ts := t.Setting.Current()

	if ts.AutoShutdown {
		t.AutomaticShutDownTicker(name, 5*time.Second, sigChan)
	}

	<-sigChan

	t.Stop(name)
}

func (t *Tekojar) AutomaticShutDownTicker(name string, interval time.Duration, sigChan chan os.Signal) {
	ticker := time.NewTicker(interval)

	go func() {
		defer func() {
			ticker.Stop()
			PrintLog(SYSTEM, 0, "Shutdown automaticlly")
		}()

		for range ticker.C {
			s, _ := t.GetService(name)
			_, totalInactive := t.GetTotalStatusServices()
			if name != "" && s.Status == INACTIVE {
				sigChan <- syscall.SIGTERM
				return // exit goroutine after triggering
			} else if len(t.services) == totalInactive {
				sigChan <- syscall.SIGTERM
				return // exit goroutine after triggering
			}
		}
	}()
}
