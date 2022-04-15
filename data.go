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
	indicator0      [100]int
	indicator1      [100]int
	planet0         [100]int
	planet1         [100]int
	planet2         [100]int
	typical         [100]float64
	timeframe       = os.Getenv("TIMEFRAME")
	tolerance       = int64(10800000) // 3 hours in miliseconds
	userAgent       = randUserAgent()
)

/*CreateData - create sample data.
 */
func CreateData() {
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
		// ma1 := talib.Ma(typical[:], 20, talib.SMA)
		dx := talib.Dx(high[:], low[:], clos[:], 10)
		// copy(indicator0[:], ma0)
		// copy(indicator1[:], dx)

		moon := IlluminationTimeseries("Moon", &timeseries)
		mercury := IlluminationTimeseries("Mercury", &timeseries)
		venus := IlluminationTimeseries("Venus", &timeseries)
		min, max := MinMax(ma0[10:])
		for i, _ := range high {
			indicator0[i] = Scale(Normalize(typical[i], min, max))
			indicator1[i] = Scale((Normalize(dx[i], 0.0, 100.0)))
			planet0[i] = Scale((*moon)[i])
			planet1[i] = Scale((*mercury)[i])
			planet2[i] = Scale((*venus)[i])
		}

		// copy(planet0[:], (*moon))
		// copy(planet1[:], (*mercury))
		// copy(planet2[:], (*venus))
	} else {
		log.Println("Something went wrong")
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
