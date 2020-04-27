package geoclue2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/godbus/dbus/v5"
	"time"
)

// Paths of methods and properties
const (
	GeoclueClientInterface  = GeoclueInterface + ".Client"
	GeoclueClientObjectPath = GeoclueObjectPath + "/Client"

	/* Methods */
	GeoclueClientStart = GeoclueClientInterface + ".Start"
	GeoclueClientStop  = GeoclueClientInterface + ".Stop"

	/* Property */
	GeoclueClientPropertyLocation               = GeoclueClientInterface + ".Location"               // readable   o
	GeoclueClientPropertyDistanceThreshold      = GeoclueClientInterface + ".DistanceThreshold"      // readwrite  u
	GeoclueClientPropertyTimeThreshold          = GeoclueClientInterface + ".TimeThreshold"          //  readwrite  u
	GeoclueClientPropertyDesktopId              = GeoclueClientInterface + ".DesktopId"              // readwrite  s
	GeoclueClientPropertyRequestedAccuracyLevel = GeoclueClientInterface + ".RequestedAccuracyLevel" // readwrite  u
	GeoclueClientPropertyActive                 = GeoclueClientInterface + ".Active"                 // readable   b

	/* SIGNAL */
	GeoclueClientSignalLocationUpdated = "LocationUpdated"
)

// GeoclueClient interface you use to retrieve location information and receive location update signals from GeoClue service.
// You get the client object to use this interface on from org.freedesktop.GeoClue2.Manager.GetClient() method.
type GeoclueClient interface {

	/* METHODS */

	// Start receiving events about the current location. Applications should hook-up to
	// LocationUpdated" signal before calling this method.
	Start() error

	// Stop receiving events about the current location.
	Stop() error

	/* PROPERTIES */

	// Current location as path to a org.freedesktop.Geoclue2.Location object. Please note that this property will
	// be set to "/" (D-Bus equivalent of null) initially, until Geoclue finds user's location.
	// You want to delay reading this property until your callback to "LocationUpdated" signal is
	// called for the first time after starting the client.
	GetLocation() (GeoclueLocation, error)

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
	// your machine, external resources and/or how much accuracy user agrees to be comfortable with.
	GetRequestedAccuracyLevel() (GClueAccuracyLevel, error)
	SetRequestedAccuracyLevel(level GClueAccuracyLevel) error

	// If client is active, i-e started successfully using Start() and receiving location updates.
	// Please keep in mind that geoclue can at any time stop and start the client on user (agent) request.
	// Applications that are interested in in these changes, should watch for changes in this property.
	IsActive() (bool, error)

	MarshalJSON() ([]byte, error)

	/* SIGNALS */

	// LocationUpdated signal: The signal is emitted every time the location changes. The client should set the
	// DistanceThreshold property to control how often this signal is emitted.
	// o old: old location as path to a #org.freedesktop.Geoclue2.Location object
	// o new: new location as path to a #org.freedesktop.Geoclue2.Location object
	SubscribeLocationUpdated() <-chan *dbus.Signal
	// Parse the signal and return the old and new Location
	ParseLocationUpdated(v *dbus.Signal) (oldLocation GeoclueLocation, newLocation GeoclueLocation, err error)

	Unsubscribe()
}

// NewGeoclueClient returns new GeoclueClient Interface
func NewGeoclueClient(objectPath dbus.ObjectPath) (GeoclueClient, error) {
	var gcc geoclueClient
	return &gcc, gcc.init(GeoclueInterface, objectPath)
}

type geoclueClient struct {
	dbusBase
	sigChan chan *dbus.Signal
}

func (gcc geoclueClient) Start() error {
	err := gcc.call(GeoclueClientStart)
	if err != nil {
		return err
	}
	return err
}

func (gcc geoclueClient) Stop() error {
	err := gcc.call(GeoclueClientStop)
	if err != nil {
		return err
	}
	return err
}

func (gcc geoclueClient) GetLocation() (GeoclueLocation, error) {
	var objPath dbus.ObjectPath
	cActive, err := gcc.IsActive()
	if err != nil {
		return nil, err
	}
	if !cActive {
		return nil, errors.New("client must be started before gathering the location")
	}

	timeout := time.After(5 * time.Second)
	// Please note that this property will be set to "/" (D-Bus equivalent of null) initially,
	// until Geoclue finds user's location. You want to delay reading this property until
	// your callback to "LocationUpdated" signal is called for the first time after starting the client.
	c := gcc.SubscribeLocationUpdated()
	select {
	case <-timeout:
		return nil, errors.New("timed out gathering location object")
	case <-c:
		objPath, err = gcc.getObjectProperty(GeoclueClientPropertyLocation)
		if err != nil {
			return nil, err
		}
		if len(fmt.Sprint(objPath)) > 1 {
			break
		}
	}
	gcc.Unsubscribe()

	return NewGeoclueLocation(objPath)
}

func (gcc geoclueClient) GetDistanceThreshold() (uint32, error) {
	v, err := gcc.getUint32Property(GeoclueClientPropertyDistanceThreshold)
	return v, err
}

func (gcc geoclueClient) SetDistanceThreshold(value uint32) error {
	return gcc.setProperty(GeoclueClientPropertyDistanceThreshold, value)

}

func (gcc geoclueClient) GetTimeThreshold() (uint32, error) {
	v, err := gcc.getUint32Property(GeoclueClientPropertyTimeThreshold)
	return v, err
}

func (gcc geoclueClient) SetTimeThreshold(value uint32) error {
	return gcc.setProperty(GeoclueClientPropertyTimeThreshold, value)
}

func (gcc geoclueClient) GetDesktopId() (string, error) {
	v, err := gcc.getStringProperty(GeoclueClientPropertyDesktopId)
	return v, err
}

func (gcc geoclueClient) SetDesktopId(value string) error {
	return gcc.setProperty(GeoclueClientPropertyDesktopId, value)
}

func (gcc geoclueClient) GetRequestedAccuracyLevel() (GClueAccuracyLevel, error) {
	v, err := gcc.getUint32Property(GeoclueClientPropertyRequestedAccuracyLevel)
	return GClueAccuracyLevel(v), err
}

func (gcc geoclueClient) SetRequestedAccuracyLevel(level GClueAccuracyLevel) error {
	return gcc.setProperty(GeoclueClientPropertyRequestedAccuracyLevel, level)
}

func (gcc geoclueClient) IsActive() (bool, error) {
	v, err := gcc.getBoolProperty(GeoclueClientPropertyActive)
	return v, err
}

func (gcc geoclueClient) SubscribeLocationUpdated() <-chan *dbus.Signal {
	if gcc.sigChan != nil {
		return gcc.sigChan
	}
	rule := fmt.Sprintf("type='signal', member='%s',path_namespace='%s'", GeoclueClientSignalLocationUpdated, GeoclueClientObjectPath)
	gcc.conn.BusObject().Call(dbusMethodAddMatch, 0, rule)
	gcc.sigChan = make(chan *dbus.Signal, 10)
	gcc.conn.Signal(gcc.sigChan)
	return gcc.sigChan
}
func (gcc geoclueClient) ParseLocationUpdated(v *dbus.Signal) (oldLocation GeoclueLocation, newLocation GeoclueLocation, err error) {
	if len(v.Body) != 2 {
		err = errors.New("error by parsing activation changed signal")
		return
	}
	oPath, ok := v.Body[0].(dbus.ObjectPath)
	if !ok {
		err = errors.New("error by parsing old object path")
		return
	}
	oldLocation, err = NewGeoclueLocation(oPath)
	if err != nil {
		return
	}
	nPath, ok := v.Body[1].(dbus.ObjectPath)
	if !ok {
		err = errors.New("error by parsing new object path")
		return
	}
	newLocation, err = NewGeoclueLocation(nPath)
	if err != nil {
		return
	}
	return
}

func (gcc geoclueClient) Unsubscribe() {
	gcc.conn.RemoveSignal(gcc.sigChan)
	gcc.sigChan = nil
}

func (gcc geoclueClient) MarshalJSON() ([]byte, error) {
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
