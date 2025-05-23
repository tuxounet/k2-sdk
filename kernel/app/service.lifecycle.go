package app

func (s *Service) Init() error {

	k := s.GetKernel()

	app := s.GetKernel().GetApp()
	app.SetKernel(k)

	appLogger := s.GetLogger().CreateSubLogger(app.GetName())
	app.SetLogger(appLogger)

	config := app.GetConfig()

	if config != nil {

		err := s.getConfigService().LoadFromEmbedFS("config", config)

		if err != nil {
			s.GetLogger().ErrorF("Error loading config for app %s: %s", app.GetName(), err.Error())
			return err
		}
	}
	components := app.GetComponents()
	for _, component := range components {

		compConf := component.GetConfig()
		if compConf != nil {
			err := s.getConfigService().LoadFromEmbedFS("config", compConf)
			if err != nil {
				s.GetLogger().ErrorF("failed to get config for component %s: %s", component.GetName(), err.Error())
				return err
			}
		}

		controllers := component.GetControllers()

		for _, ctrl := range controllers {
			conf := ctrl.GetConfig()
			if conf != nil {
				err := s.getConfigService().LoadFromEmbedFS("config", conf)
				if err != nil {
					s.GetLogger().ErrorF("failed to get config for controller %s: %s", ctrl.GetName(), err.Error())
					return err
				}
			}
			err := ctrl.Init()
			if err != nil {
				s.GetLogger().ErrorF("controller %s init failed: %s", ctrl.GetName(), err.Error())
				return err
			}
		}

	}
	return nil
}

func (s *Service) Start() error {

	components := s.GetKernel().GetApp().GetComponents()
	for _, component := range components {

		controllers := component.GetControllers()

		for _, ctrl := range controllers {
			s.GetLogger().TraceF("start controller %s", ctrl.GetName())
			err := ctrl.Start()
			if err != nil {
				s.GetLogger().ErrorF("controller %s start failed: %s", ctrl.GetName(), err.Error())
				return err
			}
			s.GetLogger().TraceF("started controller %s", ctrl.GetName())
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
