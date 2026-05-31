package scheduler

import (
	"time"

	"zyrouge.me/umi/application"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/utils"
)

func StartEventDeletionScheduler() error {
	config, err := application.GetConfig()
	if err != nil {
		return err
	}
	if config.Database.RetentionDays > 0 {
		go func() {
			ticker := time.NewTicker(1 * time.Hour)
			defer ticker.Stop()
			for range ticker.C {
				if err := repository.DeleteExpiredEventsByRetentionDays(config.Database.RetentionDays); err != nil {
					utils.Logger.Error().Err(err).Msg("failed to delete expired events")
				}
			}
		}()
	}
	return nil
}
