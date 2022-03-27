package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"strconv"

)

type BusStop struct {
	External_id  string `json:"external_id"`
	Forecast     []forecast `json:"forecast"` //contains predicted waiting time for buses 
	Geometry     []geometry `json:"geometry"`
	Id           int `json:"id"`
	Name         string `json:"name"`
	Name_en      string `json:"Name_en"`
	Name_ru      string `json:"Name_ru"`
	Nameslug     string `json:"Nameslug"`
	Reousrce_uri string `json:"Resource_uri"`
}

type forecast struct {
	Forecast_seconds float64 `json:"forecast_seconds"`
	Route            route `json:"route"`
	Rv_id            int `json:"rv_id"` //busline id 
	Total_pass       float64 `json:"total_pass"`
	Vehicle          string `json:"vehicle"`
	Vehicle_id       int `json:"vehicle_id"`
}

type route struct {
	Id         int `json:"id"`
	Name       string `json:"name"`
	Short_name string `json:"short_name"`
}

type geometry struct {
	External_id string `json:"external_id"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Seq         int `json:"seq"`
}

type BusLine struct {
	External_id  string `json:"external_id"`
	Id           int `json:"id"`
	Name         string `json:"name"`
	Name_en      string `json:"name_en"`
	Name_ru      string `json:"name_ru"`
	Nameslug     string `json:"nameslug"`
	Resource_uri string `json:"Resource_uri"`
	Routename    string `json:"routename"`
	Vehicles     []vehicle `json:"vehicles"`
	Via          string `json:"via"`
}

type vehicle struct {
	Bearing           int `json:"bearing"`
	Device_ts         string `json:"device_ts"`
	Enterprise        enterprise `json:"enterprise"`
	Lat               string `json:"lat"`
	Lon               string `json:"lon"`
	Park              park `json:"park"`
	Position          position `json:"position"`
	Projection        projection `json:"projection"`
	Registration_code string `json:"registration_code"`
	Routevariant_id   int `json:"routevariant_id"`
	Speed             string `json:"speed"`
	Stats             stats `json:"stats"`
	Ts                string `json:"ts"`
	Vehicle_id        int `json:"vehicle_id"`
}

type enterprise struct {
	Enterprise_id   int `json:"enterprise_id"`
	Enterprise_name string `json:"enterprise_name"`
}

type park struct {
	Park_id   int `json:"park_id"`
	Park_name string `json:"park_name"`
}

type position struct {
	Bearing   int `json:"bearing"`
	Device_ts int `json:"device_ts"`
	Lat       string `json:"lat"`
	Lon       string `json:"lon"`
	Speed     int `json:"speed"`
	Ts        int `json:"ts"`
}

type projection struct {
	Edge_distance      string `json:"edge_distance"`
	Edge_id            int `json:"edge_id"`
	Edge_projection    string `json:"edge_protection"`
	Edge_start_node_id int `json:"edge_start_node_id"`
	Edge_stop_node_id  int `json:"edge_stop_node_id"`
	Lat                string `json:"lat"`
	Lon                string `json:"lon"`
	Orig_lat           string `json:"orig_lat"`
	Orig_lon           string `json:"orig_lon"`
	Routevariant_id    int `json:"routevariant_id"`
	Ts                 int `json:"ts"`
}

type stats struct {
	Avg_speed     string `json:"avg_speed"`
	Bearing       int `json:"bearing"`
	Cumm_speed_10 string `json:"cumm_speed_10"`
	Cumm_speed_2  string `json:"cumm_speed_2"`
	Device_ts     int `json:"device_ts"`
	Lat           string `json:"lat"`
	Lon           string `json:"lon"`
	Speed         int `json:"speed"`
	Ts            int `json:"ts"`
}


var busStops = []string {"378204", "383050", "378202", "383049", "382998", "378237", "378233", "378230", "378229", "378228", "378227", "382995", "378224", "378226", "383010", "383009",
	"383006", "383004", "378234", "383003", "378222", "383048", "378203", "382999",
	"378225", "383014", "383013", "383011", "377906", "383018", "383015", "378207"}

var busLines = []string{"44478", "44479", "44480", "44481"}

func returnMap (busStops [] string) map[string]string { //returns a map with key : bus location name and value : bus location id 

	var busStopNames = make(map[string]string)

	for i := 0; i < len(busStops); i++ {

		response, err := http.Get("https://dummy.uwave.sg/busstop/" + busStops[i])

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var BusStop BusStop
		json.Unmarshal(responseData, &BusStop)
		name := BusStop.Name 

		busStopNames[name] = busStops[i]

    }
	
	return busStopNames 

} 


func routeBusStops (busStops [] string, busLineId int) []BusStop {//returns the bus stops where buses of a particular busline stop at 

	routeStops := make([]BusStop, 10)

	for i := 0; i < len(busStops); i++ {

		response, err := http.Get("https://dummy.uwave.sg/busstop/" + busStops[i])

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var BusStop BusStop
		json.Unmarshal(responseData, &BusStop)
		forecasts := BusStop.Forecast

		for j := 0; j < len(forecasts); j++ {

			if forecasts[j].Rv_id == busLineId{
				routeStops = append(routeStops, BusStop)
		}
	}
    }
	return routeStops
}

func returnRoute (routeID string) string { //returns the name of the busline given the id 

		var route string

		switch routeID {

			case "44478":
				route = "Campus Loop Red"
			case "44479":
				route = "Campus Loop Blue"
			case "44480":
				route = "Campus Rider Green"
			case "44481":
				route = "Campus Weekend Rider Brown"
		}

		return route
	}

func returnBusListByLocation (w  http.ResponseWriter, r *http.Request){ //returns the buses and their estimated waiting time at a particular bus stop
    
	busStopNames := returnMap (busStops)

	vars := mux.Vars(r)
    key := vars["name"]

	id, exist :=  busStopNames[key]
    if exist {
    
		endpoint := "https://dummy.uwave.sg/busstop/" + id

		response, err := http.Get(endpoint)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Fprintf(w,string(responseData))

		var BusStop BusStop
		json.Unmarshal(responseData, &BusStop)


		for i := 0; i < len(BusStop.Forecast); i++ {

			bus := BusStop.Forecast[i]
			fmt.Fprintf(w, "Bus ID: " + strconv.Itoa(bus.Vehicle_id))
			fmt.Fprintf(w, " running in " + bus.Route.Short_name)
			fmt.Fprintf(w," is arriving in ")
			waitingTime := int(bus.Forecast_seconds)
			waitingTime /= 60
			fmt.Fprintf(w,strconv.Itoa(waitingTime) + " minutes")
			fmt.Fprintln(w, "")

		}
	} else {
		fmt.Fprintf(w, key + " is not valid. Please type in one of the following :\n")
		fmt.Fprintln(w, "")
			for busStopName, _ := range busStopNames {
				fmt.Fprintln(w,"Bus Stop Name: " + busStopName)
			}
		
	}
}


func returnBusStopByName(w  http.ResponseWriter, r *http.Request){ //returns the bus stop based on the name 
    
	busStopNames := returnMap (busStops)

	vars := mux.Vars(r)
    key := vars["name"]

	id, exist :=  busStopNames[key]
    if exist {
    

		endpoint := "https://dummy.uwave.sg/busstop/" + id

		response, err := http.Get(endpoint)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w,string(responseData))

		var BusStop BusStop
		json.Unmarshal(responseData, &BusStop)
		} else {
			fmt.Fprintf(w, key + " is not valid. Please type in one of the following :\n")
			fmt.Fprintln(w, "")
			for busStopName, _ := range busStopNames {
				fmt.Fprintln(w,"Bus Stop Name: " + busStopName)
			}
			
		}
	}



func returnBusStopById(w http.ResponseWriter, r *http.Request){ //returns the bus stop based on the id 
    vars := mux.Vars(r)
    key := vars["id"]

	endpoint := "https://dummy.uwave.sg/busstop/" + key

	response, err := http.Get(endpoint)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w,string(responseData))

	var BusStop BusStop
	json.Unmarshal(responseData, &BusStop)

	}

func returnBusListByLocationAndBusLineId(w http.ResponseWriter, r *http.Request){ //returns a specific bus and its estimated waiting time 
    
	busStopNames := returnMap (busStops)

	vars := mux.Vars(r)
    key := vars["name"]
	busLineId := vars["id"]

	id, exist :=  busStopNames[key]

    if exist {

		endpoint := "https://dummy.uwave.sg/busstop/" + id

		response, err := http.Get(endpoint)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Fprintf(w,string(responseData))

		var BusStop BusStop
		json.Unmarshal(responseData, &BusStop)

		found := false

		for i := 0; i < len(BusStop.Forecast); i++ {

			bus := BusStop.Forecast[i]
			if (strconv.Itoa(bus.Rv_id) == busLineId) {
				fmt.Fprintf(w, "Bus ID: " + strconv.Itoa(bus.Vehicle_id))
				fmt.Fprintf(w, " running in " + bus.Route.Short_name)
				fmt.Fprintf(w," is arriving in ")
				waitingTime := int(bus.Forecast_seconds)
				waitingTime /= 60
				fmt.Fprintf(w,strconv.Itoa(waitingTime) + " minutes")
				found = true
				break
			}
		}

		if found == false {
			fmt.Fprintf(w, "There are no buses running in " + returnRoute(busLineId) + " which stop at this bus stop.")
		}

	
	} else {
		fmt.Fprintf(w, key + " is not valid. Please type in one of the following :\n")
		fmt.Fprintln(w, "")
		for busStopName, _ := range busStopNames {
			fmt.Fprintln(w,"Bus Stop Name: " + busStopName)
		}
		
	}
}

func returnBusLineById(w http.ResponseWriter, r *http.Request) { //returns all the buses that are running in the bus line and their estimated waiting time for each bus stop

	vars := mux.Vars(r)
    key := vars["id"]

	endpoint := "https://dummy.uwave.sg/busline/" + key

	response, err := http.Get(endpoint)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Fprintf(w,string(responseData))

	var BusLine BusLine
	json.Unmarshal(responseData, &BusLine)

	routeStops := routeBusStops(busStops, BusLine.Id)

	for i := 0; i < len(BusLine.Vehicles); i++ {

		bus := BusLine.Vehicles[i]
		fmt.Fprintf(w,"Bus Id: " + strconv.Itoa(bus.Vehicle_id))
		fmt.Fprintf(w," Latitude: " + bus.Position.Lat)
		fmt.Fprintf(w, " Longitude: " + bus.Position.Lon)
		fmt.Fprintln(w, " ")

		busStopTimings := make(map[string]int)

		for j:= 0; j < len(routeStops); j++ {
			routeBusStop := routeStops[j]
			busStopName := routeBusStop.Name
			if _, exist := busStopTimings[busStopName]; exist {
				continue
			} else {
				for k:= 0; k < len(routeBusStop.Forecast); k++ {
					if routeBusStop.Forecast[k].Vehicle_id == bus.Vehicle_id {
						time := int(routeBusStop.Forecast[k].Forecast_seconds) /  60
						busStopTimings[busStopName] = time
						fmt.Fprintf(w, "Arriving at " + busStopName + " in " + strconv.Itoa(time) + " minutes")
						fmt.Fprintln(w," ")
						break
					}
				}
			}
			
		}
		fmt.Fprintln(w," ")

	}

	if len(BusLine.Vehicles) == 0 {
		fmt.Fprintf(w, "There are currently no buses running in this bus line.")
	}
}

func homePage(w http.ResponseWriter, r *http.Request){ //list instructions and details of endpoints 
	http.ServeFile(w, r, r.URL.Path[1:])
}


func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
   
	myRouter.HandleFunc("/BusStop/{id}", returnBusStopById)

	myRouter.HandleFunc("/BusStopName/{name}", returnBusStopByName)

	myRouter.HandleFunc("/BusStopLists/{name}", returnBusListByLocation)

	myRouter.HandleFunc("/BusStopList/{name}/{id}", returnBusListByLocationAndBusLineId)

	myRouter.HandleFunc("/BusLine/{id}", returnBusLineById)

    // finally, instead of passing in nil, we want
    // to pass in our newly created router as the second
    // argument
    log.Fatal(http.ListenAndServe(":8080", myRouter))

}

func main() { 

	handleRequests()

}
