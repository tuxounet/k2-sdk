package app

import "fmt"

func (s *Service) Init() error {
	s.GetLogger().TraceF("[INIT]")
	k := s.GetKernel()

	app := s.GetKernel().GetApp()
	app.SetKernel(k)

	appLogger := s.GetLogger().CreateSubLogger(app.GetName())
	app.SetLogger(appLogger)

	config := app.GetConfig()

	if config != nil {

		s.GetLogger().TraceF("[INIT] init config for app %s", app.GetName())
		err := s.getConfigService().LoadFromEmbedFS("app", "config", config)

		if err != nil {
			return fmt.Errorf("failed to get config for app %s: %s", app.GetName(), err.Error())
		}
		s.GetLogger().TraceF("[INIT] inited config for app %s", app.GetName())
	}
	components := app.GetComponents()
	for _, component := range components {

		compConf := component.GetConfig()
		if compConf != nil {
			s.GetLogger().TraceF("[INIT] init config for component %s", component.GetName())
			err := s.getConfigService().LoadFromEmbedFS(fmt.Sprintf("%s/%s", app.GetName(), component.GetName()), "config", compConf)
			if err != nil {
				return fmt.Errorf("failed to get config for component %s: %s", component.GetName(), err.Error())
			}
			s.GetLogger().TraceF("[INIT] inited config for component %s", component.GetName())
		}

		controllers := component.GetControllers()

		for _, ctrl := range controllers {
			s.GetLogger().TraceF("[INIT] init controller %s", ctrl.GetName())
			conf := ctrl.GetConfig()
			if conf != nil {
				err := s.getConfigService().LoadFromEmbedFS(fmt.Sprintf("%s/%s/%s", app.GetName(), component.GetName(), ctrl.GetName()), "config", conf)
				if err != nil {
					return fmt.Errorf("failed to get config for controller %s: %s", ctrl.GetName(), err.Error())
				}
			}
			err := ctrl.Init()
			if err != nil {
				return fmt.Errorf("controller %s init failed: %w", ctrl.GetName(), err)
			}
			s.GetLogger().TraceF("[INIT] inited controller %s", ctrl.GetName())
		}

	}
	s.GetLogger().TraceF("[INIT] complete")
	return nil
}

func (s *Service) Start() error {

	s.GetLogger().TraceF("[START]")
	components := s.GetKernel().GetApp().GetComponents()
	for _, component := range components {

		controllers := component.GetControllers()

		for _, ctrl := range controllers {
			s.GetLogger().TraceF("[START] start controller %s", ctrl.GetName())
			err := ctrl.Start()
			if err != nil {
				return fmt.Errorf("controller %s start failed: %w", ctrl.GetName(), err)
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
			s.GetLogger().TraceF("stop controller %s", ctrl.GetName())
			err := ctrl.Stop()
			if err != nil {
				s.GetLogger().ErrorF("controller %s Stop failed: %s", ctrl.GetName(), err.Error())
				return err
			}
			s.GetLogger().TraceF("stopped controller %s", ctrl.GetName())
		}

	}

	return nil
}
