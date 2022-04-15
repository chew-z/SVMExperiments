package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Illumination struct {
	Date         time.Time `firestore:"date,omitempty"`
	DateUnix     int64     `firestore:"date_unix,omitempty"`
	Illumination float64   `firestore:"illumination,omitempty"`
}

/* IlluminationTimeseries - return Illumination for planet
on given days (from quotes timeseries)
*/
func IlluminationTimeseries(planet string, tq *[]float64) *[]float64 {
	t1 := int64((*tq)[0]) - tolerance
	t2 := int64((*tq)[len(*tq)-1]) + tolerance
	illArr, _ := retrieveIlluminationRangeByUnix(planet, t1, t2)
	iL2 := SearchAndMatch(tq, illArr)
	// TODO - add error handling here when don't match
	// log.Printf("Records for %s # - quotes: %d\tsearchAndMatch(): %d", planet, len(*tq), len(*iL2))
	return iL2
}

/* SearchAndMatch - search two slices and match aproximated dates
one from quotes and one from astro tables
*/
func SearchAndMatch(tq *[]float64, illArr *[]Illumination) *[]float64 {
	var result []float64
	for i, t := range *tq {
		for _, d := range (*illArr)[i:] { // This is crude algo but we are talking 100-150 items
			a := int64(t)
			b := d.DateUnix
			if isEqualInt64(a, b, tolerance) {
				result = append(result, d.Illumination)
				break
			}
		}
	}
	return &result
}

func retrieveIlluminationRangeByUnix(planet string, start int64, end int64) (*[]Illumination, error) {
	path := fmt.Sprintf("planets/%s/illumination", planet)
	var illArr []Illumination
	it := firestoreClient.Collection(path).Where("date_unix", ">=", start).Where("date_unix", "<=", end).OrderBy("date_unix", firestore.Asc).Documents(ctx)
	i := 0
	for {
		doc, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("documents iterator: %v", err)
		}
		i++
		var ill Illumination
		doc.DataTo(&ill)
		illArr = append(illArr, ill)
	}
	// log.Printf("Retrieved %d records for %s", i, planet)
	return &illArr, nil
}

/*initFirestoreDatabase - as the name says creates Firestore client
in Google Cloud it is using project ID, on localhost credentials file
It works for AppEngine, CloudRun/Docker and local testing
*/
func initFirestoreDatabase(ctx context.Context) *firestore.Client {
	// Default - local testing
	sa := option.WithCredentialsFile(".firebase-credentials.json")
	// Cloud credentials and roles
	firestoreClient, err := firestore.NewClient(ctx, firestore.DetectProjectID, sa)
	if err != nil {
		log.Panic(err)
	}
	return firestoreClient
}
