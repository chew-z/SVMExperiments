package main

import (
	"math/rand"
	"time"
)

func init() {
	firestoreClient = initFirestoreDatabase(ctx)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	defer firestoreClient.Close()

	ts := CreateQuotesData(asset)
	CreatePlanetData(ts)
	SaveTrainData(asset)
	CreateModel(asset)

	// x := make(map[int]float64)
	// // Populate x with the test vector

	// predictLabel := model.Predict(x) // Predicts a float64 label given the test vector
	// fmt.Println(predictLabel)
}
