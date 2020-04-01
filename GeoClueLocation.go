package go_geoclue2

import (
	"encoding/json"
	"github.com/godbus/dbus/v5"
	"time"
)

const (
	GeoClueLocationInterface  = GeoClueInterface + ".Location"
	GeoClueLocationObjectPath = GeoClueObjectPath + "/Location"

	/* Methods */

	/* Property */
	// The latitude of the location, in degrees.
	GeoClueLocationPropertyLatitude = GeoClueLocationInterface + ".Latitude" // readable   d
	// The longitude of the location, in degrees.
	GeoClueLocationPropertyLongitude = GeoClueLocationInterface + ".Longitude" // readable   d
	// The accuracy of the location fix, in meters.
	GeoClueLocationPropertyAccuracy = GeoClueLocationInterface + ".Accuracy" // readable   d
	// The altitude of the location fix, in meters. When unknown, its set to minimum double value, -1.7976931348623157e+308.
	GeoClueLocationPropertyAltitude = GeoClueLocationInterface + ".Altitude" // readable   d
	// The speed in meters per second. When unknown, it's set to -1.0.
	GeoClueLocationPropertySpeed = GeoClueLocationInterface + ".Speed" // readable   d
	// The heading direction in degrees with respect to North direction, in clockwise order. That means North becomes 0 degree, East: 90 degrees, South: 180 degrees, West: 270 degrees and so on. When unknown, it's set to -1.0.
	GeoClueLocationPropertyHeading = GeoClueLocationInterface + ".Heading" // readable  d
	//A human-readable description of the location, if available. WARNING: Applications should not rely on this property since not all sources provide a description. If you really need a description (or more details) about current location, use a reverse-geocoding API, e.g geocode-glib.
	GeoClueLocationPropertyDescription = GeoClueLocationInterface + ".Description" // readable  s
	//The timestamp when the location was determined, in seconds and microseconds since the Epoch. This is the time of measurement if the backend provided that information, otherwise the time when GeoClue received the new location.
	// Note that GeoClue can't guarantee that the timestamp will always monotonically increase, as a backend may not respect that. Also note that a timestamp can be very old, e.g. because of a cached location.
	GeoClueLocationPropertyTimestamp = GeoClueLocationInterface + ".Timestamp" // readable   (tt)

)

type GeoClueLocation interface {
	/* METHODS */
	GetLatitude() (float64, error)
	GetLongitude() (float64, error)
	GetAccuracy() (float64, error)
	GetAltitude() (float64, error)
	GetSpeed() (float64, error)
	GetHeading() (float64, error)
	GetDescription() (string, error)
	GetTimestamp() (time.Time, error)
	MarshalJSON() ([]byte, error)
}

func NewGeoClueLocation(objectPath dbus.ObjectPath) (GeoClueLocation, error) {
	var gcl geoClueLocation
	return &gcl, gcl.init(GeoClueInterface, objectPath)
}

type geoClueLocation struct {
	dbusBase

	sigChan chan *dbus.Signal
}

func (gcl geoClueLocation) GetLatitude() (float64, error) {
	v, err := gcl.getFloat64Property(GeoClueLocationPropertyLatitude)
	return v, err
}

func (gcl geoClueLocation) GetLongitude() (float64, error) {
	v, err := gcl.getFloat64Property(GeoClueLocationPropertyLongitude)
	return v, err
}

func (gcl geoClueLocation) GetAccuracy() (float64, error) {
	v, err := gcl.getFloat64Property(GeoClueLocationPropertyAccuracy)
	return v, err
}

func (gcl geoClueLocation) GetAltitude() (float64, error) {
	v, err := gcl.getFloat64Property(GeoClueLocationPropertyAltitude)
	return v, err
}

func (gcl geoClueLocation) GetSpeed() (float64, error) {
	v, err := gcl.getFloat64Property(GeoClueLocationPropertySpeed)
	return v, err
}

func (gcl geoClueLocation) GetHeading() (float64, error) {
	v, err := gcl.getFloat64Property(GeoClueLocationPropertyHeading)
	return v, err
}

func (gcl geoClueLocation) GetDescription() (string, error) {
	v, err := gcl.getStringProperty(GeoClueLocationPropertyDescription)
	return v, err
}

func (gcl geoClueLocation) GetTimestamp() (time.Time, error) {
	v, err := gcl.getTimestampProperty(GeoClueLocationPropertyTimestamp)

	return v, err
}

func (gcl geoClueLocation) MarshalJSON() ([]byte, error) {
	latitude, err := gcl.GetLatitude()
	if err != nil {
		return nil, err
	}
	longitude, err := gcl.GetLongitude()
	if err != nil {
		return nil, err
	}
	accuracy, err := gcl.GetAccuracy()
	if err != nil {
		return nil, err
	}
	altitude, err := gcl.GetAltitude()
	if err != nil {
		return nil, err
	}

	heading, err := gcl.GetHeading()
	if err != nil {
		return nil, err
	}

	description, err := gcl.GetDescription()
	if err != nil {
		return nil, err
	}

	timestamp, err := gcl.GetTimestamp()
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"Latitude":    latitude,
		"Longitude":   longitude,
		"Accuracy":    accuracy,
		"Altitude":    altitude,
		"Heading":     heading,
		"Description": description,
		"Timestamp":   timestamp,
	})
}
