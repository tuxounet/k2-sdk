package scheduler

import "time"

func (s *Service) Init() error {
	return nil
}

func (s *Service) Start() error {

	cronInstance := s.getCronRunner()
	kernel := s.GetKernel()
	app := kernel.GetApp()

	components := app.GetComponents()
	registrations := 0

	for _, component := range components {
		schedules := component.GetControllers()
		for _, controller := range schedules {
			for _, schedule := range controller.GetSchedules() {
				s.GetLogger().InfoF("Registering schedule '%s' with cron '%s'", schedule.GetName(), schedule.GetCron())
				cronInstance.AddFunc(schedule.GetCron(), func() {

					fireTime := time.Now()
					s.GetLogger().DebugF("Schedule '%s' fired at %s", schedule.GetName(), fireTime.Format(time.RFC3339))
					result := schedule.GetTaskHandler()(fireTime.Format(time.RFC3339))
					endTime := time.Now()
					duration := endTime.Sub(fireTime)

					s.GetLogger().DebugF("Schedule '%s' completed at %s (duration: %s) with result: %s", schedule.GetName(), endTime.Format(time.RFC3339), duration.String(), result)
				})
				registrations++
			}
		}
	}
	if registrations == 0 {
		s.GetLogger().WarnF("No schedules registered, cron runner will not start")
		return nil
	}
	s.GetLogger().InfoF("Starting cron runner with %d registered schedules", registrations)
	cronInstance.Start()

	return nil
}

func (s *Service) Stop() error {
	cronInstance := s.getCronRunner()
	cronInstance.Stop()
	return nil
}
