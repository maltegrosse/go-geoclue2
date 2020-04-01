package main

import (
	"fmt"
	"github.com/maltegrosse/go-geoclue2"
	"log"
	"os"
)

func main() {
	// create new Instance of GeoClue Manager
	gcm, err := go_geoclue2.NewGeoClueManager()
	if err != nil {
		log.Fatal(err.Error())

	}
	// create new Instance of GeoClue Client
	client, err := gcm.GetClient()
	if err != nil {
		log.Fatal(err.Error())
	}
	// desktop id is required to start the client
	// (double check your geoclue.conf file)
	err = client.SetDesktopId("firefox")
	if err != nil {
		log.Fatal(err.Error())
	}
	// client must be started before requesting the location
	err = client.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	// create new Instance of GeoClue Location
	location, err := client.GetLocation()
	if err != nil {
		log.Fatal(err.Error())
	}
	// get latitude
	latitude, err := location.GetLatitude()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(latitude)

	// get longitude
	longitude, err := location.GetLongitude()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(longitude)

	// get accuracy
	accuracy, err := location.GetAccuracy()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(accuracy)


/*
	// listen location updates
	c := client.Subscribe()
	for v := range c {
		fmt.Println(v)
	}
*/
	os.Exit(0)
}
