package geoclue2

import (
	"encoding/json"
	"github.com/godbus/dbus/v5"
	"time"
)

// Paths of methods and properties
const (
	GeoclueLocationInterface = GeoclueInterface + ".Location"

	/* Methods */

	/* Property */

	GeoclueLocationPropertyLatitude = GeoclueLocationInterface + ".Latitude" // readable   d

	GeoclueLocationPropertyLongitude = GeoclueLocationInterface + ".Longitude" // readable   d

	GeoclueLocationPropertyAccuracy = GeoclueLocationInterface + ".Accuracy" // readable   d

	GeoclueLocationPropertyAltitude = GeoclueLocationInterface + ".Altitude" // readable   d

	GeoclueLocationPropertySpeed = GeoclueLocationInterface + ".Speed" // readable   d

	GeoclueLocationPropertyHeading = GeoclueLocationInterface + ".Heading" // readable  d

	GeoclueLocationPropertyDescription = GeoclueLocationInterface + ".Description" // readable  s

	GeoclueLocationPropertyTimestamp = GeoclueLocationInterface + ".Timestamp" // readable   (tt)

)

// GeoclueLocation interface you use on location objects.
type GeoclueLocation interface {
	/* METHODS */
	// The latitude of the location, in degrees.
	GetLatitude() (float64, error)
	// The longitude of the location, in degrees.
	GetLongitude() (float64, error)
	// The accuracy of the location fix, in meters.
	GetAccuracy() (float64, error)
	// The altitude of the location fix, in meters. When unknown, its set to minimum double value, -1.7976931348623157e+308.
	GetAltitude() (float64, error)
	// The speed in meters per second. When unknown, it's set to -1.0.
	GetSpeed() (float64, error)
	// The heading direction in degrees with respect to North direction, in clockwise order. That means North becomes 0 degree, East: 90 degrees, South: 180 degrees, West: 270 degrees and so on. When unknown, it's set to -1.0.
	GetHeading() (float64, error)
	//A human-readable description of the location, if available. WARNING: Applications should not rely on this property since not all sources provide a description. If you really need a description (or more details) about current location, use a reverse-geocoding API, e.g geocode-glib.
	GetDescription() (string, error)
	//The timestamp when the location was determined, in seconds and microseconds since the Epoch. This is the time of measurement if the backend provided that information, otherwise the time when Geoclue received the new location.
	// Note that Geoclue can't guarantee that the timestamp will always monotonically increase, as a backend may not respect that. Also note that a timestamp can be very old, e.g. because of a cached location.
	GetTimestamp() (time.Time, error)
	MarshalJSON() ([]byte, error)
}

// NewGeoclueLocation returns new NewGeoclueLocation Interface
func NewGeoclueLocation(objectPath dbus.ObjectPath) (GeoclueLocation, error) {
	var gcl geoclueLocation
	return &gcl, gcl.init(GeoclueInterface, objectPath)
}

type geoclueLocation struct {
	dbusBase
	sigChan chan *dbus.Signal
}

func (gcl geoclueLocation) GetLatitude() (float64, error) {
	v, err := gcl.getFloat64Property(GeoclueLocationPropertyLatitude)
	return v, err
}

func (gcl geoclueLocation) GetLongitude() (float64, error) {
	v, err := gcl.getFloat64Property(GeoclueLocationPropertyLongitude)
	return v, err
}

func (gcl geoclueLocation) GetAccuracy() (float64, error) {
	v, err := gcl.getFloat64Property(GeoclueLocationPropertyAccuracy)
	return v, err
}

func (gcl geoclueLocation) GetAltitude() (float64, error) {
	v, err := gcl.getFloat64Property(GeoclueLocationPropertyAltitude)
	return v, err
}

func (gcl geoclueLocation) GetSpeed() (float64, error) {
	v, err := gcl.getFloat64Property(GeoclueLocationPropertySpeed)
	return v, err
}

func (gcl geoclueLocation) GetHeading() (float64, error) {
	v, err := gcl.getFloat64Property(GeoclueLocationPropertyHeading)
	return v, err
}

func (gcl geoclueLocation) GetDescription() (string, error) {
	v, err := gcl.getStringProperty(GeoclueLocationPropertyDescription)
	return v, err
}

func (gcl geoclueLocation) GetTimestamp() (time.Time, error) {
	v, err := gcl.getTimestampProperty(GeoclueLocationPropertyTimestamp)

	return v, err
}

func (gcl geoclueLocation) MarshalJSON() ([]byte, error) {
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
