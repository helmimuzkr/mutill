package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Mutill struct {
	config   *MutillConfig
	services map[string]*Service

	mu sync.RWMutex

	logStarted bool
	logCh      chan string
	ShutdownCh chan os.Signal
}

func (m *Mutill) RegisterService(config *MutillConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if config == nil {
		return
	}

	m.config = config

	if m.config.Services == nil {
		panic("no service inputed")
	}

	PrintLog(SYSTEM, 999, fmt.Sprintf("config -> %v", config))

	if m.services == nil {
		m.services = make(map[string]*Service)
	}

	for _, s := range m.config.Services {
		if s.Skip {
			continue
		}
		m.services[s.Name] = InitService(s.Name, s.Path, s.Args)
	}
}

func (m *Mutill) GetService(name string) *Service {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.services[name]
}

func (m *Mutill) StartAll() {
	PrintLog(SYSTEM, 999, "starting application...")

	m.mu.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, v := range m.services {
		if v.Status == INACTIVE {
			services = append(services, v)
		}
	}
	m.mu.RUnlock()
	m.listenLog()

	for _, s := range services {
		cmd := "cmd"
		if m.config.Command != "" {
			cmd = m.config.Command
		}
		args := []string{"/c", s.Path}
		if s.Args != nil {
			args = s.Args
		}
		go s.StartProcess(cmd, args...)
	}

	go func() {
		for log := range m.ReadLog() {
			fmt.Println(log)
		}
	}()

	m.ListenShutdown()
}

func (m *Mutill) StopAll() {
	m.mu.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, v := range m.services {
		services = append(services, v)
	}
	m.mu.RUnlock()

	for _, v := range services {
		v.StopProcess()
	}

	actives, inactives := m.GetTotalStatusServices()
	PrintLog(SYSTEM, 999, "application stopped", fmt.Sprintf("actives: %d", actives), fmt.Sprintf("inactives: %d", inactives))
}

func (m *Mutill) listenLog() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.logStarted {
		panic("ServiceChannelLog() already called")
	}

	m.logStarted = true

	m.logCh = make(chan string, 100*len(m.services))

	var wg sync.WaitGroup

	for _, s := range m.services {
		wg.Add(1)
		go func(svc *Service) {
			defer wg.Done()
			for log := range svc.LogCh {
				m.logCh <- log
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(m.logCh)
	}()
}

func (m *Mutill) ReadLog() chan string {
	if m.logCh == nil {
		m.logCh = make(chan string, 100*len(m.services))
	}
	return m.logCh
}

func (m *Mutill) GetTotalStatusServices() (int, int) {
	active := 0
	inactive := 0

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, s := range m.services {
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
func (m *Mutill) ListenShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	if m.config.AutoShutdown {
		m.AutomaticShutDownTicker(5*time.Second, sigChan)
	}

	<-sigChan
	m.StopAll()
}

func (m *Mutill) AutomaticShutDownTicker(interval time.Duration, sigChan chan os.Signal) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			_, totalInactive := m.GetTotalStatusServices()
			if len(m.services) == totalInactive {
				sigChan <- syscall.SIGQUIT
			}
		}
	}()
}
