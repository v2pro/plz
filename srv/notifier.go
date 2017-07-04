package srv

type serverNotifier struct {
	ch chan bool
}

func (signal *serverNotifier) Stop() {
	defer func() {
		recover()
	}()
	if signal.IsStopped() {
		return
	}
	close(signal.ch) // might panic in race condition
}

func (signal *serverNotifier) IsStopped() bool {
	select {
	case <-signal.ch:
		return true
	default:
		return false
	}
}

func (signal *serverNotifier) Wait() {
	<-signal.ch
}

type multipleServerNotifiers struct {
	notifiers []Notifier
}

func (mss *multipleServerNotifiers) Stop() {
	for _, notifier := range mss.notifiers {
		notifier.Stop()
	}
}

func (mss *multipleServerNotifiers) IsStopped() bool {
	for _, notifier := range mss.notifiers {
		if !notifier.IsStopped() {
			return false
		}
	}
	return true
}

func (mss *multipleServerNotifiers) Wait() {
	for _, notifier := range mss.notifiers {
		notifier.Wait()
	}
}