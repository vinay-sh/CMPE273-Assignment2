package main

import (
"fmt"
"log"
"io/ioutil"
"net/http"
"encoding/json"
"github.com/jasonwinn/geocoder"
"math/rand"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
"github.com/julienschmidt/httprouter"
"strconv"
)


type Coordinates struct{
	Lat float64 `json: "Lat"`
	Lng float64 `json: "Lng"`
}

type Response struct{
	Id int `json: "Id"`
	Name string `json: "Name"`
	Address string `json: "Address"`
	City string `json: "City"`
	State string `json: "State"`
	Zip string `json: "Zip"`
	Coordinate Coordinates `json: "Coordinate"`
}


func CreateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var idata Response

	session :=connectToDb()

    c:= session.DB("vinaysh").C("admin")


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

	randnum := rand.Int();
	fmt.Println("Random number generated is",randnum)

	query:= idata.Address + " " + idata.City + " " + idata.State + " " + idata.Zip
	fmt.Println("Query is", query)
	lat1, lng1, err := geocoder.Geocode(query)
  	if err != nil {
    	panic("THERE WAS SOME ERROR!!")
  }

	fmt.Println("Lat Long Coordinates")

	idata.Coordinate.Lat = lat1
	idata.Coordinate.Lng = lng1

	fmt.Println("Latitude is",idata.Coordinate.Lat)
	fmt.Println("Longitude is",idata.Coordinate.Lng)


	lol:= Coordinates{idata.Coordinate.Lat, idata.Coordinate.Lng}
	resp1, err := json.Marshal(lol)
    if err != nil {
    panic(err)
    }
    fmt.Println("Marshaled lol is", string(resp1))


	if err := c.Insert(&Response{randnum,idata.Name,idata.Address,idata.City,idata.State,idata.Zip,idata.Coordinate});err != nil {
    	rw.WriteHeader(404)
    	return 
    }


	rData:= Response{randnum, idata.Name, idata.Address, idata.City, idata.State, idata.Zip, idata.Coordinate}
	fmt.Println("Response is",rData)

	resp, err := json.Marshal(rData)
    if err != nil {
    panic(err)
    }

    fmt.Println("Marshaled Json Response is", string(resp))

    rw.Header().Set("Content-Type", "application/json")
    rw.Write(resp)

}


func ReadLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

	session :=connectToDb()

	c:= session.DB("vinaysh").C("admin")
	var result Response
    id,_ := strconv.Atoi(p.ByName("id"))

    err := c.Find(bson.M{"id":id }).One(&result)

    if err != nil {
        rw.WriteHeader(404)
        return 
            }
        js, err := json.Marshal(result)
        if err != nil {
        panic(err)
        }
        rw.Write(js)

}


func UpdateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    var idata Response

    session :=connectToDb()

    c:= session.DB("vinaysh").C("admin")


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

	fmt.Println("Unmarshaled string is", idata.Address, idata.City, idata.State, idata.Zip)

	query:= idata.Address + " " + idata.City + " " + idata.State + " " + idata.Zip
	fmt.Println("Query is", query)
	lat1, lng1, err := geocoder.Geocode(query)
  	if err != nil {
    	panic("THERE WAS SOME ERROR!!")
  }

	fmt.Println("Lat Long Coordinates")

	idata.Coordinate.Lat = lat1
	idata.Coordinate.Lng = lng1

	id,_ := strconv.Atoi(p.ByName("id"))

    if err := c.Update(bson.M{"id": id}, bson.M{"id": id,"name":idata.Name,"address":idata.Address,"city":idata.City,"state":idata.State,"zip":idata.Zip,"coordinate":idata.Coordinate});err != nil {
     	rw.WriteHeader(404)
        return 
            }


	rData:= Response{id, idata.Name, idata.Address, idata.City, idata.State, idata.Zip, idata.Coordinate}
	fmt.Println("Response is",rData)

	resp, err := json.Marshal(rData)
    if err != nil {
    panic(err)
    }

    fmt.Println("Marshaled Json Response is", string(resp))

    rw.Header().Set("Content-Type", "application/json")
    rw.Write(resp)

}


func DeleteLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    session :=connectToDb()

    c:= session.DB("vinaysh").C("admin")

    fmt.Println("Name",p.ByName("id"))
    id,_ := strconv.Atoi(p.ByName("id"))
    err := c.Remove(bson.M{"id":id})
    if err != nil {
        rw.WriteHeader(404)
        return 
    }
}

func connectToDb() *mgo.Session{
    session, err := mgo.Dial("mongodb://admin:admin@ds045054.mongolab.com:45054/vinaysh")
        if err != nil {
                panic("Couldn't connect to the database")
        }
           session.SetMode(mgo.Monotonic, true)
           fmt.Println("Session is ",session)
    return session
            
}

func main(){
     mux := httprouter.New()
     mux.POST("/location",CreateLocation) 
     mux.GET("/location/:id",ReadLocation)
     mux.PUT("/location/:id",UpdateLocation)
     mux.DELETE("/location/:id",DeleteLocation)
     server := http.Server{
             Addr:        "127.0.0.1:8080",
             Handler: mux,
     }
     server.ListenAndServe()
}