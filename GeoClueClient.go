package go_geoclue2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/godbus/dbus/v5"
	"time"
)

const (
	GeoClueClientInterface  = GeoClueInterface + ".Client"
	GeoClueClientObjectPath = GeoClueObjectPath + "/Client"

	/* Methods */
	GeoClueClientStart = GeoClueClientInterface + ".Start"
	GeoClueClientStop  = GeoClueClientInterface + ".Stop"

	/* Property */
	GeoClueClientPropertyLocation               = GeoClueClientInterface + ".Location"               // readable   o
	GeoClueClientPropertyDistanceThreshold      = GeoClueClientInterface + ".DistanceThreshold"      // readwrite  u
	GeoClueClientPropertyTimeThreshold          = GeoClueClientInterface + ".TimeThreshold"          //  readwrite  u
	GeoClueClientPropertyDesktopId              = GeoClueClientInterface + ".DesktopId"              // readwrite  s
	GeoClueClientPropertyRequestedAccuracyLevel = GeoClueClientInterface + ".RequestedAccuracyLevel" // readwrite  u
	GeoClueClientPropertyActive                 = GeoClueClientInterface + ".Active"                 // readable   b

)

type GeoClueClient interface {
	/* METHODS */

	// Start receiving events about the current location. Applications should hook-up to
	// LocationUpdated" signal before calling this method.
	Start() error

	// Stop receiving events about the current location.
	Stop() error

	// Current location as path to a org.freedesktop.GeoClue2.Location object. Please note that this property will
	// be set to "/" (D-Bus equivalent of null) initially, until Geoclue finds user's location.
	// You want to delay reading this property until your callback to "LocationUpdated" signal is
	// called for the first time after starting the client.
	GetLocation() (GeoClueLocation, error)

	// Contains the current distance threshold in meters. This value is used by the service when it gets new location info.
	// If the distance moved is below the threshold, it won't emit the LocationUpdated signal. The default value is 0.
	// When TimeThreshold is zero, it always emits the signal.
	GetDistanceThreshold() (uint32, error)
	SetDistanceThreshold(uint32) error

	// Contains the current time threshold in seconds. This value is used by the service when it gets new location info.
	// If the time since the last update is below the threshold, it won't emit the LocationUpdated signal.
	// The default value is 0. When TimeThreshold is zero, it always emits the signal.
	GetTimeThreshold() (uint32, error)
	SetTimeThreshold(uint32) error

	// The desktop file id (the basename of the desktop file).
	// This property must be set by applications for authorization to work.
	// e.g. firefox
	GetDesktopId() (string, error)
	SetDesktopId(string) error

	// The level of accuracy requested by client, as GClueAccuracyLevel.
	// Please keep in mind that the actual accuracy of location information is dependent on available hardware on
	// your machine, external resources and/or how much accuracy user agrees to be confortable with.
	GetRequestedAccuracyLevel() (GClueAccuracyLevel, error)
	SetRequestedAccuracyLevel(level GClueAccuracyLevel) error

	// If client is active, i-e started successfully using Start() and receiving location updates.
	// Please keep in mind that geoclue can at any time stop and start the client on user (agent) request.
	// Applications that are interested in in these changes, should watch for changes in this property.
	IsActive() (bool, error)

	// LocationUpdated signal: The signal is emitted every time the location changes. The client should set the
	// DistanceThreshold property to control how often this signal is emitted.
	// o old: old location as path to a #org.freedesktop.GeoClue2.Location object
	// o new: new location as path to a #org.freedesktop.GeoClue2.Location object
	Subscribe() <-chan *dbus.Signal
	Unsubscribe()

	MarshalJSON() ([]byte, error)
}

func NewGeoClueClient(objectPath dbus.ObjectPath) (GeoClueClient, error) {
	var gcc geoClueClient
	return &gcc, gcc.init(GeoClueInterface, objectPath)
}

type geoClueClient struct {
	dbusBase
	sigChan chan *dbus.Signal
}

func (gcc geoClueClient) Start() error {
	err := gcc.call(GeoClueClientStart)
	if err != nil {
		return err
	}
	return err
}

func (gcc geoClueClient) Stop() error {
	err := gcc.call(GeoClueClientStop)
	if err != nil {
		return err
	}
	return err
}

func (gcc geoClueClient) GetLocation() (GeoClueLocation, error) {
	var objPath dbus.ObjectPath
	var err error

	cActive,err :=gcc.IsActive()
	if !cActive{
		return nil, errors.New("client must be started before gathering the location")
	}

	timeout := time.After(5 * time.Second)
	// Please note that this property will be set to "/" (D-Bus equivalent of null) initially,
	// until Geoclue finds user's location. You want to delay reading this property until
	// your callback to "LocationUpdated" signal is called for the first time after starting the client.
	c := gcc.Subscribe()
	select {
	case <-timeout:
		return nil, errors.New("timed out gathering location object")
	case <-c:
		objPath, err = gcc.getObjectProperty(GeoClueClientPropertyLocation)
		if err != nil {
			return nil, err
		}
		if len(fmt.Sprint(objPath)) > 1 {
			break
		}
	}
	gcc.Unsubscribe()

	gcl, err := NewGeoClueLocation(objPath)

	return gcl, err
}

func (gcc geoClueClient) GetDistanceThreshold() (uint32, error) {
	v, err := gcc.getUint32Property(GeoClueClientPropertyDistanceThreshold)
	return v, err
}

func (gcc geoClueClient) SetDistanceThreshold(value uint32) error {
	return gcc.setProperty(GeoClueClientPropertyDistanceThreshold, value)

}

func (gcc geoClueClient) GetTimeThreshold() (uint32, error) {
	v, err := gcc.getUint32Property(GeoClueClientPropertyTimeThreshold)
	return v, err
}

func (gcc geoClueClient) SetTimeThreshold(value uint32) error {
	return gcc.setProperty(GeoClueClientPropertyTimeThreshold, value)
}

func (gcc geoClueClient) GetDesktopId() (string, error) {
	v, err := gcc.getStringProperty(GeoClueClientPropertyDesktopId)
	return v, err
}

func (gcc geoClueClient) SetDesktopId(value string) error {
	return gcc.setProperty(GeoClueClientPropertyDesktopId, value)
}

func (gcc geoClueClient) GetRequestedAccuracyLevel() (GClueAccuracyLevel, error) {
	v, err := gcc.getUint32Property(GeoClueClientPropertyRequestedAccuracyLevel)
	return GClueAccuracyLevel(v), err
}

func (gcc geoClueClient) SetRequestedAccuracyLevel(level GClueAccuracyLevel) error {
	return gcc.setProperty(GeoClueClientPropertyRequestedAccuracyLevel, level)
}

func (gcc geoClueClient) IsActive() (bool, error) {
	v, err := gcc.getBoolProperty(GeoClueClientPropertyActive)
	return v, err
}

func (gcc geoClueClient) Subscribe() <-chan *dbus.Signal {
	if gcc.sigChan != nil {
		return gcc.sigChan
	}

	gcc.subscribeNamespace(GeoClueClientObjectPath)
	gcc.sigChan = make(chan *dbus.Signal, 10)
	gcc.conn.Signal(gcc.sigChan)

	return gcc.sigChan
}

func (gcc geoClueClient) Unsubscribe() {
	gcc.conn.RemoveSignal(gcc.sigChan)
	gcc.sigChan = nil
}

func (gcc geoClueClient) MarshalJSON() ([]byte, error) {
	location, err := gcc.GetLocation()
	if err != nil {
		return nil, err
	}
	mLocation, err := location.MarshalJSON()
	if err != nil {
		return nil, err
	}
	distanceThreshold, err := gcc.GetDistanceThreshold()
	if err != nil {
		return nil, err
	}

	timeThreshold, err := gcc.GetTimeThreshold()
	if err != nil {
		return nil, err
	}
	desktopId, err := gcc.GetDesktopId()
	if err != nil {
		return nil, err
	}
	requestedAccuracyLevel, err := gcc.GetRequestedAccuracyLevel()
	if err != nil {
		return nil, err
	}
	isActive, err := gcc.IsActive()
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"Location":               mLocation,
		"DistanceThreshold  ":    distanceThreshold,
		"TimeThreshold":          timeThreshold,
		"DesktopId":              desktopId,
		"RequestedAccuracyLevel": requestedAccuracyLevel,
		"Active":                 isActive,
	})
}
