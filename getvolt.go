package main

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"math"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"time"
)

// Define a struct to match the structure of the JSON response
type sensor_vals struct {
	StatusSNS struct {
		Time   string `json:"Time"`
		Ina219 struct {
			ID      int     `json:"Id"`
			Voltage float64 `json:"Voltage"`
			Current float64 `json:"Current"`
			Power   int     `json:"Power"`
		} `json:"INA219"`
	} `json:"StatusSNS"`
}

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type voltCollector struct {
	voltMetric *prometheus.Desc
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newVoltCollector() *voltCollector {
	return &voltCollector{
		voltMetric: prometheus.NewDesc("volt_metric",
			"Voltage tasmota",
			nil, nil,
		),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *voltCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.voltMetric
}

//Collect implements required collect function for all promehteus collectors
func (collector *voltCollector) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue float64
	metricValue += getVolt()

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	m1 := prometheus.MustNewConstMetric(collector.voltMetric, prometheus.GaugeValue, metricValue)
//	m2 := prometheus.MustNewConstMetric(collector.barMetric, prometheus.GaugeValue, metricValue)
	m1 = prometheus.NewMetricWithTimestamp(time.Now(), m1)
//	m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
	ch <- m1
//	ch <- m2
}


func getVolt() float64 {
	// Define the URL of the API endpoint you want to request
	url := "http://172.16.0.101/cm?user=admin&password=sti789654!&cmnd=Status%2010"

	// Send an HTTP GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		//fmt.Println("Failed to make the HTTP request:", err)
		log.Fatalf("Failed to make the HTTP request:", err)
		return 0
	}
	defer response.Body.Close()

	// Check if the request was successful (status code 200)
	if response.StatusCode != http.StatusOK {
		//fmt.Printf("Failed to retrieve data. Status code: %d\n", response.StatusCode)
		log.Fatalf("Failed to retrieve data. Status code: %d\n", response.StatusCode)
		return 0
	}

	// Decode the JSON response into the Post struct
	var vals sensor_vals
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&vals); err != nil {
		//fmt.Println("Failed to decode JSON:", err)
		log.Fatalf("Failed to decode JSON:", err)
		return 0
	}

	// Extract and print the specific value from the JSON response
	volt := math.Round(vals.StatusSNS.Ina219.Voltage*10)/10

	//fmt.Printf("Voltage: %f\n",vals.StatusSNS.Ina219.Voltage*)
	//fmt.Printf("Voltage: %.1f\n", volt)A
	return volt
}

func main() {
	volt := newVoltCollector()
	prometheus.MustRegister(volt)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
