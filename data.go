package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	libSvm "github.com/ewalker544/libsvm-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markcheno/go-talib"
)

/*Quotes ...
 */
type Quotes struct {
	DefaultChartInterval string          `json:"_default_chart_interval,omitempty"`
	RefPrice             float64         `json:"_ref_price,omitempty"`
	Bars                 [][]interface{} `json:"_d"`
}

/*Candle ...
 */
type Candle struct {
	Time string
	// BOSSA API is OHLC go-echarts OCLH
	OHLC [4]float64
}

var (
	apiURL = os.Getenv("API_URL")
	asset  = os.Getenv("ASSET")
	city   = os.Getenv("CITY")
	ctx    = context.Background()
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
	firestoreClient *firestore.Client
	kd              [100]Candle
	location, _     = time.LoadLocation(city)
	price           [100]int
	indicator0      [100]int
	indicator1      [100]int
	indicator2      [100]int
	indicator3      [100]int
	planet0         [100]int
	planet1         [100]int
	planet2         [100]int
	position0       [100]int
	position1       [100]int
	position2       [100]int
	typical         [100]float64
	timeframe       = os.Getenv("TIMEFRAME")
	tolerance       = int64(10800000) // 3 hours in miliseconds
	userAgent       = randUserAgent()
)

/*
 */
func CreateModel(asset string) {
	dataPath := fmt.Sprintf("%s.train", asset)
	modelPath := fmt.Sprintf("%s.model", asset)
	param := libSvm.NewParameter()  // Create a parameter object with default values
	param.KernelType = libSvm.POLY  // Use the polynomial kernel
	model := libSvm.NewModel(param) // Create a model object from the parameter attributes

	// Create a problem specification from the training data and parameter attributes
	problem, _ := libSvm.NewProblem(dataPath, param)
	model.Train(problem)  // Train the model from the problem specification
	model.Dump(modelPath) // Dump the model into a user-specified file

}

/*
 */
func SaveTrainData(asset string) {
	dataPath := fmt.Sprintf("%s.train", asset)
	f, err := os.Create(dataPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range planet0 {
		sig := Signum(indicator3[i] - indicator2[i])
		fmt.Fprintf(f, "%d %d:1 %d:1 %d:1 %d:1 %d:1 %d:1 %d:1 %d:1 %d:1 %d:1 %d:1\n", sig, planet0[i], planet1[i], planet2[i], position0[i], position1[i], position2[i], price[i], indicator0[i], indicator1[i], indicator2[i], indicator3[i])
	}
	f.Close()

}

/*CreateData - create sample data.
 */
func CreatePlanetData(timeseries []float64) {
	moon := IlluminationTimeseries("Moon", timeseries)
	mercury := IlluminationTimeseries("Mercury", timeseries)
	venus := IlluminationTimeseries("Venus", timeseries)
	moonPos := PositionTimeseries("Moon", timeseries)
	mercuryPos := PositionTimeseries("Mercury", timeseries)
	venusPos := PositionTimeseries("Venus", timeseries)
	for i, _ := range timeseries {
		planet0[i] = Scale(moon[i])
		planet1[i] = Scale(mercury[i])
		planet2[i] = Scale(venus[i])
		position0[i] = Scale(Normalize(moonPos[i], 0.0, 360.0))
		position1[i] = Scale(Normalize(mercuryPos[i], 0.0, 360.0))
		position2[i] = Scale(Normalize(venusPos[i], 0.0, 360.00))
	}
}

func CreateQuotesData(asset string) []float64 {
	var high, low [100]float64
	var open, clos [100]float64
	var timeseries []float64
	if quotes := getQuotes(asset, timeframe); quotes != nil {
		for i, bar := range quotes.Bars {
			var tmp Candle
			tm := int64(bar[0].(float64))
			time := time.Unix(0, tm*int64(time.Millisecond))
			if timeframe == "D1" {
				tmp.Time = time.In(location).Format("Jan _2")
			} else {
				tmp.Time = time.In(location).Format("Jan _2 15:04")
			}
			o, _ := bar[1].(float64)
			h, _ := bar[2].(float64)
			l, _ := bar[3].(float64)
			c, _ := bar[4].(float64)
			open[i] = o
			high[i] = h
			low[i] = l
			clos[i] = c
			typical[i] = (h + l + c) / 3.0    // typical price - HLC/3
			tmp.OHLC = [4]float64{o, c, l, h} // OHLC to OCLH
			kd[i] = tmp
			timeseries = append(timeseries, bar[0].(float64))
		}
		ma0 := talib.Ma(typical[:], 10, talib.SMA)
		dx := talib.Dx(high[:], low[:], clos[:], 10)
		aroonDown, aroonUp := talib.Aroon(high[:], low[:], 10)
		min, max := MinMax(ma0[10:])
		for i, _ := range high {
			price[i] = Scale(Normalize(typical[i], min, max))
			indicator0[i] = Scale(Normalize(typical[i], min, max))
			indicator1[i] = Scale((Normalize(dx[i], 0.0, 100.0)))
			indicator2[i] = Scale(Normalize(aroonDown[i], 0.0, 100.0))
			indicator3[i] = Scale(Normalize(aroonUp[i], 0.0, 100.0))
		}
		return timeseries
	} else {
		log.Println("Something went wrong")
		return nil
	}
}

func getQuotes(asset string, timeframe string) *Quotes {
	var quotes Quotes
	apiURL := fmt.Sprintf("%s%s.", apiURL, asset)
	if timeframe != "" {
		apiURL += "/" + timeframe
	}
	request, _ := http.NewRequest("GET", apiURL, nil)
	request.Header.Set("User-Agent", userAgent)
	if response, err := client.Do(request); err == nil {
		if err := json.NewDecoder(response.Body).Decode(&quotes); err != nil {
			log.Fatal(err)
		}
	}
	return &quotes
}
