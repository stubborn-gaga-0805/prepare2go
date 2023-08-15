package runtime

var (
	stopListenerChan = make(chan bool, 0)

	dataWatcherChan = make(chan string, 0)
	mqWatcherChan   = make(chan string, 0)
)
