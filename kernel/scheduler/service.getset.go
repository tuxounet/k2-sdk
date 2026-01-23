package scheduler

import "github.com/robfig/cron/v3"

func (s *Service) getCronRunner() *cron.Cron {
	return s.GetData("cron").(*cron.Cron)
}

func (s *Service) setCronRunner(c *cron.Cron) {
	s.SetData("cron", c)
}
