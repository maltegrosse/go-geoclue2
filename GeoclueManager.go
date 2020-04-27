package geoclue2

import (
	"encoding/json"
	"github.com/godbus/dbus/v5"
	"log"
)

// Paths of methods and properties
const (
	GeoclueInterface         = "org.freedesktop.GeoClue2"
	GeoclueManagerInterface  = GeoclueInterface + ".Manager"
	GeoclueObjectPath        = "/org/freedesktop/GeoClue2"
	GeoclueManagerObjectPath = GeoclueObjectPath + "/Manager"

	/* Methods */
	GeoclueManagerGetClient    = GeoclueManagerInterface + ".GetClient"
	GeoclueManagerCreateClient = GeoclueManagerInterface + ".CreateClient"
	GeoclueManagerDeleteClient = GeoclueManagerInterface + ".DeleteClient"
	GeoclueManagerAddAgent     = GeoclueManagerInterface + ".AddAgent"

	/* Property */
	GeoclueManagerPropertyInUse                  = GeoclueManagerInterface + ".InUse"                  // readable   b
	GeoclueManagerPropertyAvailableAccuracyLevel = GeoclueManagerInterface + ".AvailableAccuracyLevel" //readable   u              // readable   ao
)

// GeoclueManager interface you use to talk to main GeoClue2 manager object at path "/org/freedesktop/GeoClue2/Manager".
// The only thing you do with this interface is to call GetClient() or CreateClient() on it to get your application specific client object(s).
type GeoclueManager interface {
	/* METHODS */

	// Retrieves a client object which can only be used by the calling application only.
	// On the first call from a specific D-Bus peer, this method will create the client object but subsequent
	// calls will return the path of the existing client.
	GetClient() (GeoclueClient, error)

	// Creates and retrieves a client object which can only be used by the calling application only.
	// Unlike GetClient(), this method always creates a new client.
	CreateClient() (GeoclueClient, error)

	// Use this method to explicitly destroy a client, created using GetClient() or CreateClient().
	// Long-running applications, should either use this to delete associated client(s) when not needed,
	// or disconnect from the D-Bus connection used for communicating with Geoclue (which is implicit
	// on client process termination).
	DeleteClient(GeoclueClient) error

	// An API for user authorization agents to register themselves. Each agent is responsible for the user
	// it is running as. Application developers can and should simply ignore this API.
	// IN s id: The Desktop ID (excluding .desktop) of the agent
	AddAgent(id string) error

	/* PROPERTIES */

	// Whether service is currently is use by any application.
	InUse() (bool, error)
	// The level of available accuracy, as GClueAccuracyLevel.
	GetAvailableAccuracyLevel() (GClueAccuracyLevel, error)

	MarshalJSON() ([]byte, error)
}

// NewGeoclueManager returns new GeoclueManager Interface
func NewGeoclueManager() (GeoclueManager, error) {
	var gcm geoclueManager

	return &gcm, gcm.init(GeoclueInterface, GeoclueManagerObjectPath)
}

type geoclueManager struct {
	dbusBase
}

func (gcm geoclueManager) GetAvailableAccuracyLevel() (GClueAccuracyLevel, error) {
	v, err := gcm.getUint32Property(GeoclueManagerPropertyAvailableAccuracyLevel)
	return GClueAccuracyLevel(v), err
}

func (gcm geoclueManager) InUse() (bool, error) {
	v, err := gcm.getBoolProperty(GeoclueManagerPropertyInUse)
	return v, err

}

func (gcm geoclueManager) GetClient() (GeoclueClient, error) {
	var clientPath dbus.ObjectPath
	err := gcm.callWithReturn(&clientPath, GeoclueManagerGetClient)
	if err != nil {
		log.Fatal(err.Error())
	}
	gcc, err := NewGeoclueClient(clientPath)

	return gcc, err
}

func (gcm geoclueManager) CreateClient() (GeoclueClient, error) {
	var clientPath dbus.ObjectPath
	err := gcm.callWithReturn(&clientPath, GeoclueManagerCreateClient)
	if err != nil {
		log.Fatal(err.Error())
	}
	gcc, err := NewGeoclueClient(clientPath)

	return gcc, err
}

func (gcm geoclueManager) DeleteClient(gcc GeoclueClient) error {
	err := gcm.call(GeoclueManagerDeleteClient, &gcc)
	if err != nil {
		return err
	}
	return err
}

func (gcm geoclueManager) AddAgent(id string) error {
	return gcm.call(GeoclueManagerAddAgent, id)
}

func (gcm geoclueManager) MarshalJSON() ([]byte, error) {
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
