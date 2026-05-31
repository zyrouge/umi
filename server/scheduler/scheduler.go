package scheduler

func StartSchedulers() error {
	if err := StartEventDeletionScheduler(); err != nil {
		return err
	}
	return nil
}
