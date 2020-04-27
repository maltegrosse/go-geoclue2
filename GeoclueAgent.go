package geoclue2

import "encoding/json"

// todo: not implemented

// Paths of methods and properties
const (
	GeoclueAgentInterface = GeoclueInterface + ".Agent"

	GeoclueAgentObjectPath = GeoclueObjectPath + "/Agent"

	/* Methods */
	GeoclueAgentAuthorizeApp = GeoclueAgentInterface + ".AuthorizeApp"

	/* Property */
	GeoclueAgentPropertyMaxAccuracyLevel = GeoclueAgentInterface + ".MaxAccuracyLevel" // readable   u
)

// GeoclueAgent interface all application-authorizing agents must implement. There must be a separate agent object for every logged-in user on path "/org/freedesktop/GeoClue2/Agent".
type GeoclueAgent interface {
	/* METHODS */

	//	This is the method that will be called by geoclue to get applications authorized to be given location information.
	//IN string desktop_id: The desktop file id (the basename of the desktop file) of the application requesting location information.
	//IN u req_accuracy_level: The level of location accuracy requested by client, as GClueAccuracyLevel.
	//OUT b authorized: Return value indicating if application should be given location information or not.
	//OUT u allowed_accuracy_level: The level of location accuracy allowed for client, as GClueAccuracyLevel.
	AuthorizeApp(string, GClueAccuracyLevel) (bool, GClueAccuracyLevel, error)

	/* PROPERTIES */
	// The global maximum level of accuracy allowed for all clients. Since agents are per-user, this can be different for each user. See GClueAccuracyLevel for possible values.
	GetMaxAccuracyLevel() (GClueAccuracyLevel, error)

	MarshalJSON() ([]byte, error)
}

// NewGeoclueAgent returns new NewGeoclueAgent Interface
func NewGeoclueAgent() (GeoclueAgent, error) {
	var gca geoclueAgent
	return &gca, gca.init(GeoclueInterface, GeoclueAgentObjectPath)
}

type geoclueAgent struct {
	dbusBase
}

func (gca geoclueAgent) AuthorizeApp(desktopId string, reqLevel GClueAccuracyLevel) (authorized bool, allowedLevel GClueAccuracyLevel, err error) {
	var tmpUint uint32
	err = gca.callWithReturn2(&authorized, &tmpUint, GeoclueAgentAuthorizeApp, &desktopId, &reqLevel)
	if err != nil {
		return false, GClueAccuracyLevelNone, err
	}
	allowedLevel = GClueAccuracyLevel(tmpUint)
	return
}

func (gca geoclueAgent) GetMaxAccuracyLevel() (GClueAccuracyLevel, error) {
	res, err := gca.getUint32Property(GeoclueAgentPropertyMaxAccuracyLevel)
	if err != nil {
		return GClueAccuracyLevelNone, err
	}
	return GClueAccuracyLevel(res), nil
}

func (gca geoclueAgent) MarshalJSON() ([]byte, error) {
	maxAccuracyLevel, err := gca.GetMaxAccuracyLevel()
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]interface{}{
		"MaxAccuracyLevel": maxAccuracyLevel,
	})
}
