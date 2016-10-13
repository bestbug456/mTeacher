package queue

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworkers int, WorkQueue chan WorkRequest, testResource TestFile, output []string) error {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		worker, err := NewWorker(i+1, WorkerQueue, testResource.Data, output)
		if err != nil {
			return err
		}
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				worker := <-WorkerQueue
				worker <- work
			}
		}
	}()

	return nil
}
