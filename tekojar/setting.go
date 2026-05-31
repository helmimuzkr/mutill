package tekojar

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Setting struct {
	path           string
	currentSetting *TekojarSetting
}

type TekojarSetting struct {
	Command         string           `json:"command"`
	AutoShutdown    bool             `json:"auto_shutdown"`
	ServiceSettings []ServiceSetting `json:"service_settings"`
}

type ServiceSetting struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	SkipFlag bool   `json:"skip_flag"`
	Delay    int    `json:"delay"`
}

func LoadSetting(path string) (*Setting, error) {
	path, _ = defaultSettingPath()
	s := &Setting{path: path}
	if _, err := s.Load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Setting) Load() (*TekojarSetting, error) {
	f, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read setting: %w", err)
	}

	var tekojarSetting TekojarSetting
	if err := json.Unmarshal(f, &tekojarSetting); err != nil {
		return nil, fmt.Errorf("failed to parse setting: %w", err)
	}

	s.currentSetting = &tekojarSetting
	return s.currentSetting, nil
}

func (s *Setting) Save(ts TekojarSetting) error {
	data, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal setting: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write setting: %w", err)
	}

	s.currentSetting = &ts
	return nil
}

func (s *Setting) Current() *TekojarSetting {
	if s.currentSetting == nil {
		ts, _ := s.Load()
		s.currentSetting = ts
	}
	return s.currentSetting
}

func defaultSettingPath() (string, error) {
	path := ""

	if devPath := os.Getenv("DEV_SETTING_PATH"); devPath != "" {
		path = devPath
		PrintLog(SYSTEM, 0, "Dev Path", path)
	}

	if path == "" {
		exe, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("failed to get executable path: %w", err)
		}
		path = filepath.Join(filepath.Dir(exe), "settings.json")
		PrintLog(SYSTEM, 0, "Prod Path", path)
	}

	return path, nil
}
