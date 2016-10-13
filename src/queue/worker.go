package queue

import (
	"log"
)

import (
	"core"
)

// Start: 12:41:19
// END:

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest, tsData []Testdata, output []string) (Worker, error) {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest, 1),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		TsData:      tsData,
		Output:      output}
	return worker, nil
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				var ris Response
				ris.RequestId = work.RequestId
				ris.Result, ris.Correct, ris.NSV = generateTrainset(work.TsToAnalyze, w.TsData, w.Output)
				work.Response <- ris
			case <-w.QuitChan:
				// We have been asked to stop.
				log.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// generateTrainset function use the
// inputset and try to create the best
// trainset
func generateTrainset(trainSet []core.TrainsetData, inputset []Testdata, output []string) ([]core.TrainsetData, int, int) {

	correct := 0

	param, problem, err := core.Setup(trainSet)
	if err != nil {
		log.Fatalf("Error while setup the core:", err)
	}
	model := core.Train(param, problem)
	if model.L != len(trainSet) {
		for i := 0; i < len(inputset); i++ {
			features := inputset[i].Feature
			ris := core.Predict(features, model)
			if output[int(ris)-1] == inputset[i].Name[0] {
				correct++
			}
		}
	}

	return trainSet, correct, model.L
}
