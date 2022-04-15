package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	libSvm "github.com/ewalker544/libsvm-go"
)

func init() {
	firestoreClient = initFirestoreDatabase(ctx)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	defer firestoreClient.Close()

	CreateData()

	f, err := os.Create("astro.train")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range planet0 {
		fmt.Fprintf(f, "%d %d:1 %d:1 %d:1 %d:1 %d:1\n", EagleOrTail(), planet0[i], planet1[i], planet2[i], indicator0[i], indicator1[i])
	}
	f.Close()

	param := libSvm.NewParameter()  // Create a parameter object with default values
	param.KernelType = libSvm.POLY  // Use the polynomial kernel
	model := libSvm.NewModel(param) // Create a model object from the parameter attributes

	// Create a problem specification from the training data and parameter attributes
	problem, _ := libSvm.NewProblem("astro.train", param)
	model.Train(problem)      // Train the model from the problem specification
	model.Dump("astro.model") // Dump the model into a user-specified file

	// x := make(map[int]float64)
	// // Populate x with the test vector

	// predictLabel := model.Predict(x) // Predicts a float64 label given the test vector
	// fmt.Println(predictLabel)
}
