package command

// Standard lib
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
)

// Exsternal dependences libraries
import (
	"github.com/spf13/cobra"
)

// Internal dependences
import (
	"core"
	"queue"
)

var TeachCmd = &cobra.Command{
	Use:     "teach",
	Short:   "Teach a machine",
	Long:    `Teaching a machine in order to solve a problem`,
	Example: "mdota teach -t /path/to/testset.json -m machine_type -w worker_number",
}

var (
	MachineType  string
	TestPath     string
	WorkerNumber int
)

func init() {
	TeachCmd.PersistentFlags().StringVarP(&MachineType, "machine", "m", "svm", "the path contain user to analyze")
	TeachCmd.PersistentFlags().StringVarP(&TestPath, "tpath", "t", "", "the path contain the test set file")
	TeachCmd.PersistentFlags().IntVarP(&WorkerNumber, "worker", "w", 100, "the number of worker to setup")
	TeachCmd.RunE = teach
}

func teach(cmd *cobra.Command, args []string) error {
	wq := make(chan queue.WorkRequest, 1000)
	testdata, err := genericImportFromFile(TestPath)
	if err != nil {
		return err
	}
	var testResource queue.TestFile
	// decode json trainset
	err = json.Unmarshal(testdata, &testResource)
	if err != nil {
		return err
	}

	input, output := analyzeInputSet(testResource.Data)
	if len(input) == 0 || len(output) == 0 {
		log.Fatalf("Error: input or output is zero. Input is: %d, output is: %d", len(input), len(output))
	}

	queue.StartDispatcher(WorkerNumber, wq, testResource, output)

	nrFeature := len(testResource.Data[0].Feature)
	lenTrainSet := len(testResource.Data)
	trainSet := make([]core.TrainsetData, lenTrainSet)
	k := 0
	for i := 0; i < lenTrainSet; i++ {
		trainSet[i].Input = make([]float64, len(testResource.Data[0].Feature))
		trainSet[i].Output = make([]string, len(testResource.Data[0].Name))
		m := 0
		for j := 0; j < len(trainSet[i].Input); j++ {
			trainSet[i].Input[j] = input[m]
			m++
			if m == len(input) {
				m = 0
			}
		}
		if i%len(output) == 0 && i != 0 {
			k++
		}
		for j := 0; j < len(trainSet[i].Output); j++ {
			trainSet[i].Output[j] = output[k]
		}

	}

	lenCycle := (int(math.Pow(float64(len(input)), float64(nrFeature))) * lenTrainSet)
	maxCorrect := 0
	numberSupportVector := 0
	var bestTs []core.TrainsetData
	response := make(chan queue.Response, 1000)
	go generateWork(trainSet, input, output, lenCycle, wq, response)
	for i := 0; i < lenCycle; i++ {
		result := <-response
		if result.Correct > maxCorrect {
			maxCorrect = result.Correct
			numberSupportVector = result.NSV
			bestTs = result.Result
			log.Println("Max correct: ", maxCorrect, " nSV:", numberSupportVector)
		}
		if ((float64(i)/float64(lenCycle))*100.00)-float64(i/(lenCycle*100.00)) == 0 {
			log.Println(int((float64(i)/float64(lenCycle))*100.00), "% done")
			log.Println("Max correct: ", maxCorrect, " nSV:", numberSupportVector)
		}
	}
	log.Println(bestTs)
	return nil
}

func generateWork(trainSet []core.TrainsetData, input []float64, output []string, lenCycle int, wq chan queue.WorkRequest, response chan queue.Response) {
	id := make(chan int, 1)
	next := make(chan bool, 1)
	lenTrainSet := len(trainSet)
	computeSlice := make([]chan bool, lenTrainSet+1)
	for i := 0; i < lenTrainSet; i++ {
		computeSlice[i] = make(chan bool, 1)
	}
	ResultSlice := make(chan []float64, 1)
	for i := 0; i < len(trainSet); i++ {
		// Run the function which handle the single slice
		workingSlice := make([]float64, len(trainSet[i].Input))
		for j := 0; j < len(trainSet[i].Input); j++ {
			workingSlice[j] = trainSet[i].Input[j]
		}
		go generateSingleSlice(workingSlice, input, output, i, id, next, ResultSlice, computeSlice[i], computeSlice[i+1])
	}
	var work queue.WorkRequest
	for i := 0; i < lenCycle; i++ {
		computeSlice[0] <- true
		for {
			// Generate the new trainset
			pos := <-id
			trainSet[pos].Input = <-ResultSlice
			nextSliceIsEdit := <-next
			if !nextSliceIsEdit {
				break
			}
		}
		// Copy the trainset in order to avoid
		// information loss
		tsCpy := make([]core.TrainsetData, len(trainSet))
		for i := 0; i < len(trainSet); i++ {
			singleInput := make([]float64, len(trainSet[i].Input))
			for j := 0; j < len(trainSet[i].Input); j++ {
				singleInput[j] = trainSet[i].Input[j]
			}
			tsCpy[i].Input = singleInput
			tsCpy[i].Output = trainSet[i].Output
		}

		work.RequestId = i
		work.Response = response
		work.TsToAnalyze = tsCpy
		wq <- work
	}

}

// function that import data from a file
func genericImportFromFile(pathfile string) ([]byte, error) {
	// read the file
	data, err := ioutil.ReadFile(pathfile)
	if err != nil {
		return data, err
	}
	return data, nil
}

func analyzeInputSet(inputset []queue.Testdata) ([]float64, []string) {
	input := make([]float64, 0)
	output := make([]string, 0)
	input = append(input, inputset[0].Feature[0])
	output = append(output, inputset[0].Name[0])
	for i := 0; i < len(inputset); i++ {
		for j := 0; j < len(inputset[i].Feature); j++ {
			inserted := false
			for k := 0; k < len(input); k++ {
				if inputset[i].Feature[j] == input[k] {
					inserted = true
				}
			}
			if !inserted {
				input = append(input, inputset[i].Feature[j])
			}
			inserted = false
		}
		for j := 0; j < len(inputset[i].Name); j++ {
			inserted := false
			for k := 0; k < len(output); k++ {
				if inputset[i].Name[j] == output[k] {
					inserted = true
				}
			}
			if !inserted {
				output = append(output, inputset[i].Name[j])
			}
			inserted = false
		}
	}
	fmt.Println(input, output)
	return input, output
}

func generateSingleSlice(workingSlice []float64, input []float64, output []string, myId int, idchan chan int, nextCall chan bool, result chan []float64, generate chan bool, nextSlice chan bool) {
	m := 0
	k := 0
	for {
		// Wait turn to generate a new slice
		_ = <-generate
		workingSlice[k] = input[m]
		m++
		if m == len(input) {
			k++
			m = 0
		}
		if k == len(workingSlice) {
			k = 0
			nextSlice <- true
			nextCall <- true
		} else {
			nextCall <- false
		}
		slice := make([]float64, len(workingSlice))
		for i := 0; i < len(slice); i++ {
			slice[i] = workingSlice[i]
		}
		fmt.Println("Input number ", myId, " :", slice)
		result <- slice
		idchan <- myId
	}
}

// function that export data to a file
/*func genericExportToFile(pathfile string, data []byte) error {
	// read the file
	err := ioutil.WriteFile(pathfile, data, 0777)
	if err != nil {
		return err
	}
	return nil
}*/
