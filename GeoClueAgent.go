package go_geoclue2
import "C"

const (
	GeoClueAgentInterface = GeoClueInterface + ".Agent"

	GeoClueAgentObjectPath = GeoClueObjectPath + "/Agent"

	/* Methods */
	GeoClueAgentAuthorizeApp = GeoClueAgentInterface + ".AuthorizeApp"

	/* Property */
	GeoClueAgentPropertyMaxAccuracyLevel = GeoClueAgentInterface + ".MaxAccuracyLevel" // readable   u
)

type GeoClueAgent interface {
	/* METHODS */

	//	This is the method that will be called by geoclue to get applications authorized to be given location information.
	//IN string desktop_id: The desktop file id (the basename of the desktop file) of the application requesting location information.
	//IN u req_accuracy_level: The level of location accuracy requested by client, as GClueAccuracyLevel.
	//OUT b authorized: Return value indicating if application should be given location information or not.
	//OUT u allowed_accuracy_level: The level of location accuracy allowed for client, as GClueAccuracyLevel.
	AuthorizeApp(string, GClueAccuracyLevel) (bool, GClueAccuracyLevel, error)

	GetMaxAccuracyLevel() (GClueAccuracyLevel, error)

	MarshalJSON() ([]byte, error)
}

func NewGeoClueAgent() (GeoClueAgent, error) {
	var gca geoClueAgent
	return &gca, gca.init(GeoClueInterface, GeoClueAgentObjectPath)
}

type geoClueAgent struct {
	dbusBase
}

func (gca geoClueAgent) AuthorizeApp(string, GClueAccuracyLevel) (bool, GClueAccuracyLevel, error) {
	panic("implement me")
}

func (gca geoClueAgent) GetMaxAccuracyLevel() (GClueAccuracyLevel, error) {
	panic("implement me")
}

func (gca geoClueAgent) MarshalJSON() ([]byte, error) {
	panic("implement me")
}
