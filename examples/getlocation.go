package main

import (
	"fmt"
	"github.com/maltegrosse/go-geoclue2"

	"log"
	"os"
)

func main() {
	// create new Instance of Geoclue Manager
	gcm, err := geoclue2.NewGeoclueManager()
	if err != nil {
		log.Fatal(err.Error())

	}
	// req availableAccuracyLevel
	aAccuracyLevel, err := gcm.GetAvailableAccuracyLevel()
	if err != nil {
		log.Fatal(err.Error())

	}
	fmt.Println("Available Accuracy Level: ", aAccuracyLevel)
	// create new Instance of Geoclue Client
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

	// Set RequestedAccuracyLevel
	err = client.SetRequestedAccuracyLevel(geoclue2.GClueAccuracyLevelExact)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get RequestedAccuracyLevel
	level, err := client.GetRequestedAccuracyLevel()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Requested Accuracy Level: ", level)

	// client must be started before requesting the location
	err = client.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	// create new Instance of Geoclue Location
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

/*	// get accuracy
	accuracy, err := location.GetAccuracy()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(accuracy)*/

	// listen location updates
	c := client.SubscribeLocationUpdated()
	for v := range c {
		fmt.Println(v)
		oldl, newl, err := client.ParseLocationUpdated(v)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("old")
		oldjson, err := oldl.MarshalJSON()
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(string(oldjson))
		fmt.Println("---")
		fmt.Println("new")
		newjson, err := newl.MarshalJSON()
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(string(newjson))

		fmt.Println("---")

	}

	os.Exit(0)
}
