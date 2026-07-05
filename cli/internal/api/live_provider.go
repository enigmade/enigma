package api

import (
	"enigma/cli/internal/ports"
	"enigma/hwdetect/pkg/hardware"
)

// LiveProvider is the production StateProvider. It reads real system state:
// project ports from ports.db and GPU info from hardware.toml. Services and
// runtimes are enumerated by the caller-supplied hooks so this stays
// testable and free of hard systemctl/mise dependencies.
type LiveProvider struct {
	HardwarePath string // /etc/enigma/hardware.toml

	// ServiceLister returns currently managed services (systemctl --user).
	// Nil is treated as "no services".
	ServiceLister func() ([]Service, error)

	// RuntimeLister returns installed runtimes (mise ls). Nil => none.
	RuntimeLister func() ([]Runtime, error)

	// ModelLister returns installed AI models (ollama list). Nil => none.
	ModelLister func() ([]Model, error)
}

// Snapshot assembles the full State from all live sources.
func (p *LiveProvider) Snapshot() (*State, error) {
	state := &State{
		Runtimes: []Runtime{},
		Services: []Service{},
		Projects: []Project{},
		Models:   []Model{},
	}

	projectPorts, err := ports.LoadDB()
	if err != nil {
		return nil, err
	}
	for _, pp := range projectPorts {
		state.Projects = append(state.Projects, Project{
			Path: pp.ProjectPath,
			Port: pp.Port,
			URL:  pp.URL,
		})
	}

	if p.HardwarePath != "" {
		if cfg, err := hardware.ReadFromFile(p.HardwarePath); err == nil && len(cfg.GPU) > 0 {
			gpu := cfg.GPU[0]
			state.GPU = &GPU{
				Vendor:  gpu.Vendor,
				Model:   gpu.Name,
				VRAMGiB: gpu.VRAMMb / 1024,
			}
		}
	}

	if p.ServiceLister != nil {
		services, err := p.ServiceLister()
		if err != nil {
			return nil, err
		}
		state.Services = services
	}

	if p.RuntimeLister != nil {
		runtimes, err := p.RuntimeLister()
		if err != nil {
			return nil, err
		}
		state.Runtimes = runtimes
	}

	if p.ModelLister != nil {
		models, err := p.ModelLister()
		if err != nil {
			return nil, err
		}
		state.Models = models
	}

	return state, nil
}
