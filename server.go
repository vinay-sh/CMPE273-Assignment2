package main

import (
"fmt"
//"net"
//"net/rpc/jsonrpc"
"log"
"io/ioutil"
"net/http"
"encoding/json"
"github.com/jasonwinn/geocoder"
"math/rand"
)

type Request struct{
	Name string
	Address string
	City string
	State string
	Zip string
}

type Coordinates struct{
	lat float64
	lng float64
}

type Response struct{
	Id int
	Name string
	Address string
	City string
	State string
	Zip string
	Coordinate Coordinates
}

func TripPlanner(rw http.ResponseWriter, req *http.Request){

	var idata Request

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal("Read Error"+err.Error())
	}

	if err:= req.Body.Close();
	err != nil{
		log.Fatal("Body Close Error"+err.Error())
	}

		fmt.Println("Input JSON data", string(body))

	if err:= json.Unmarshal(body, &idata);
	err != nil {
		panic(err)
		fmt.Println("Unmarshal failed")
	}

	fmt.Println("Unmarshaled string is", idata.Name, idata.Address, idata.City, idata.State, idata.Zip)

	query:= idata.Address + " " + idata.City + " " + idata.State + " " + idata.Zip
	fmt.Println("Query is", query)
	lat1, lng1, err := geocoder.Geocode(query)
  	if err != nil {
    	panic("THERE WAS SOME ERROR!!!!!")
  }

	fmt.Println("Lat Long Coordinates")
	latlng := Coordinates{lat1,lng1}
	fmt.Println("Latlng",latlng)

	randnum := rand.Int();
	fmt.Println("Random number generated is",randnum)

	rData:= Response{randnum, idata.Name, idata.Address, idata.City, idata.State, idata.Zip, latlng}
	fmt.Println("Response is",rData)

	resp, err := json.Marshal(rData)
    if err != nil {
    panic(err)
    }

    fmt.Println("Marshaled Json Response is", string(resp))

    rw.Header().Set("Content-Type", "application/json")
    rw.Write(resp)
}


func main(){
	http.HandleFunc("/locations", TripPlanner)
	err:=http.ListenAndServe(":8080",nil)
	if(err!=nil){
		log.Fatal("Connection Error"+err.Error())
	}
}