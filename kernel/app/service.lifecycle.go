package app

import "fmt"

func (s *Service) Init() error {
	app := s.GetKernel().GetApp()
	s.GetLogger().InfoF("[INIT] loading app %s", app.GetName())
	k := s.GetKernel()
	app.SetKernel(k)

	appLogger := s.GetLogger().CreateSubLogger(app.GetName())
	app.SetLogger(appLogger)

	components := app.GetComponents()
	for _, component := range components {

		compConf := component.GetConfig()
		if compConf != nil {
			err := s.getConfigService().LoadFromEmbedFS(fmt.Sprintf("%s/%s", app.GetName(), component.GetName()), "config", compConf)
			if err != nil {
				return fmt.Errorf("failed to load config for component %s: %w", component.GetName(), err)
			}
			s.GetLogger().TraceF("[INIT] inited config for component %s", component.GetName())
		}

		controllers := component.GetControllers()

		for _, ctrl := range controllers {
			conf := ctrl.GetConfig()
			if conf != nil {
				err := s.getConfigService().LoadFromEmbedFS(fmt.Sprintf("%s/%s/%s", app.GetName(), component.GetName(), ctrl.GetName()), "config", conf)
				if err != nil {
					return fmt.Errorf("failed to load config for controller %s in component %s: %w", ctrl.GetName(), component.GetName(), err)
				}
			}
			err := ctrl.Init()
			if err != nil {
				return fmt.Errorf("controller %s in component %s init failed: %w", ctrl.GetName(), component.GetName(), err)
			}
			s.GetLogger().TraceF("[INIT] inited controller %s", ctrl.GetName())
		}

	}
	config := app.GetConfig()

	if config != nil {

		s.GetLogger().TraceF("[INIT] init config for app %s", app.GetName())
		err := s.getConfigService().LoadFromEmbedFS("app", "config", config)

		if err != nil {
			return fmt.Errorf("failed to load config for app %s: %w", app.GetName(), err)
		}
		s.GetLogger().TraceF("[INIT] inited config for app %s", app.GetName())
	}

	s.GetLogger().TraceF("[INIT] load env vars for app %s", app.GetName())
	err := s.getConfigService().LoadFromEnvVars("host")
	if err != nil {
		return fmt.Errorf("failed to load env vars for host: %w", err)
	}

	s.GetLogger().InfoF("[INIT] app %s loaded (%d components)", app.GetName(), len(components))
	return nil
}

func (s *Service) Start() error {

	s.GetLogger().InfoF("[START] starting app controllers")
	components := s.GetKernel().GetApp().GetComponents()
	for _, component := range components {

		controllers := component.GetControllers()

		for _, ctrl := range controllers {
			err := ctrl.Start()
			if err != nil {
				return fmt.Errorf("controller %s in component %s start failed: %w", ctrl.GetName(), component.GetName(), err)
			}
			s.GetLogger().TraceF("[START] started controller %s", ctrl.GetName())
		}

	}

	return nil
}

func (s *Service) Stop() error {

	components := s.GetKernel().GetApp().GetComponents()
	for _, component := range components {

		controllers := component.GetControllers()

		for _, ctrl := range controllers {
			err := ctrl.Stop()
			if err != nil {
				s.GetLogger().ErrorF("controller %s in component %s stop failed: %s", ctrl.GetName(), component.GetName(), err.Error())
				return fmt.Errorf("controller %s in component %s stop failed: %w", ctrl.GetName(), component.GetName(), err)
			}
			s.GetLogger().TraceF("[STOP] stopped controller %s", ctrl.GetName())
		}

	}

	return nil
}
