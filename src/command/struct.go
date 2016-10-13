package command

type Testdata struct {
	Feature []float64 `json:"i,omitempty"`
	Name    []string  `json:"o,omitempty"`
}

type TestFile struct {
	Users []Testdata `json:"training_data,omitempty"`
}
