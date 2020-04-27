package geoclue2

// GClueAccuracyLevel is used to specify level of accuracy requested by, or allowed for a client.
type GClueAccuracyLevel uint32

//go:generate stringer -type=GClueAccuracyLevel -trimprefix=GClueAccuracyLevel
const (
	GClueAccuracyLevelNone         GClueAccuracyLevel = 0 // Accuracy level unknown or unset.
	GClueAccuracyLevelCountry      GClueAccuracyLevel = 1 //Country-level accuracy.
	GClueAccuracyLevelCity         GClueAccuracyLevel = 4 //City-level accuracy.
	GClueAccuracyLevelNeighborhood GClueAccuracyLevel = 5 //neighborhood-level accuracy.
	GClueAccuracyLevelStreet       GClueAccuracyLevel = 6 //Street-level accuracy.
	GClueAccuracyLevelExact        GClueAccuracyLevel = 8 //Exact accuracy. Typically requires GPS receiver.

)
