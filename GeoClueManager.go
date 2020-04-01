package go_geoclue2

import (
	"encoding/json"
	"github.com/godbus/dbus/v5"
	"log"
)

const (
	GeoClueInterface         = "org.freedesktop.GeoClue2"
	GeoClueManagerInterface  = GeoClueInterface + ".Manager"
	GeoClueObjectPath        = "/org/freedesktop/GeoClue2"
	GeoClueManagerObjectPath = GeoClueObjectPath + "/Manager"

	/* Methods */
	GeoClueManagerGetClient    = GeoClueManagerInterface + ".GetClient"
	GeoClueManagerCreateClient = GeoClueManagerInterface + ".CreateClient"
	GeoClueManagerDeleteClient = GeoClueManagerInterface + ".DeleteClient"
	GeoClueManagerAddAgent     = GeoClueManagerInterface + ".AddAgent"

	/* Property */
	GeoClueManagerPropertyInUse                  = GeoClueManagerInterface + ".InUse"                  // readable   b
	GeoClueManagerPropertyAvailableAccuracyLevel = GeoClueManagerInterface + ".AvailableAccuracyLevel" //readable   u              // readable   ao
)

type GeoClueManager interface {
	/* METHODS */

	// Retrieves a client object which can only be used by the calling application only.
	// On the first call from a specific D-Bus peer, this method will create the client object but subsequent
	// calls will return the path of the existing client.
	GetClient() (GeoClueClient, error)

	// Creates and retrieves a client object which can only be used by the calling application only.
	// Unlike GetClient(), this method always creates a new client.
	CreateClient() (GeoClueClient, error)

	// Use this method to explicitly destroy a client, created using GetClient() or CreateClient().
	// Long-running applications, should either use this to delete associated client(s) when not needed,
	// or disconnect from the D-Bus connection used for communicating with Geoclue (which is implicit
	// on client process termination).
	DeleteClient(GeoClueClient) error

	// An API for user authorization agents to register themselves. Each agent is responsible for the user
	// it is running as. Application developers can and should simply ignore this API.
	// IN s id: The Desktop ID (excluding .desktop) of the agent
	AddAgent(id string) error

	// Whether service is currently is use by any application.
	InUse() (bool, error)
	// The level of available accuracy, as GClueAccuracyLevel.
	GetAvailableAccuracyLevel() (GClueAccuracyLevel, error)

	MarshalJSON() ([]byte, error)
}

func NewGeoClueManager() (GeoClueManager, error) {
	var gcm geoClueManager
	return &gcm, gcm.init(GeoClueInterface, GeoClueManagerObjectPath)
}

type geoClueManager struct {
	dbusBase
}

func (gcm geoClueManager) GetAvailableAccuracyLevel() (GClueAccuracyLevel, error) {
	v, err := gcm.getUint32Property(GeoClueManagerPropertyAvailableAccuracyLevel)
	return GClueAccuracyLevel(v), err
}

func (gcm geoClueManager) InUse() (bool, error) {
	v, err := gcm.getBoolProperty(GeoClueManagerPropertyInUse)
	return v, err

}

func (gcm geoClueManager) GetClient() (GeoClueClient, error) {
	var clientPath dbus.ObjectPath
	err := gcm.callWithReturn(&clientPath, GeoClueManagerGetClient)
	if err != nil {
		log.Fatal(err.Error())
	}

	gcc, err := NewGeoClueClient(clientPath)

	return gcc, err

}

func (gcm geoClueManager) CreateClient() (GeoClueClient, error) {
	var clientPath dbus.ObjectPath
	err := gcm.callWithReturn(&clientPath, GeoClueManagerCreateClient)
	if err != nil {
		log.Fatal(err.Error())
	}
	gcc, err := NewGeoClueClient(clientPath)

	return gcc, err
}

func (gcm geoClueManager) DeleteClient(gcc GeoClueClient) error {
	err := gcm.call(GeoClueManagerDeleteClient, &gcc)
	if err != nil {
		return err
	}
	return err
}

func (gcm geoClueManager) AddAgent(id string) error {
	return gcm.call(GeoClueManagerAddAgent, id)
}

func (gcm geoClueManager) MarshalJSON() ([]byte, error) {
	inUse, err := gcm.InUse()
	if err != nil {
		return nil, err
	}
	availableAccuracyLevel, err := gcm.GetAvailableAccuracyLevel()
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"InUse":                   inUse,
		"AvailableAccuracyLevel ": availableAccuracyLevel,
	})
}
