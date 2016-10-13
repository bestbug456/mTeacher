package queue

import (
	"core"
)

type WorkRequest struct {
	TsToAnalyze []core.TrainsetData
	Response    chan Response
	RequestId   int
}

type Response struct {
	Result    []core.TrainsetData
	Correct   int
	NSV       int
	RequestId int
}

type TrainResult struct {
	Feature []float64 `json:"i,omitempty"`
	Name    []string  `json:"o,omitempty"`
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
	TsData      []Testdata
	Output      []string
}

type Testdata struct {
	Feature []float64 `json:"i,omitempty"`
	Name    []string  `json:"o,omitempty"`
}

type TestFile struct {
	Data []Testdata `json:"training_data,omitempty"`
}
