package go_geoclue2

type GClueAccuracyLevel uint32

const (
	None         GClueAccuracyLevel = 0 // Accuracy level unknown or unset.
	Country      GClueAccuracyLevel = 1 //Country-level accuracy.
	City         GClueAccuracyLevel = 4 //City-level accuracy.
	Neighborhood GClueAccuracyLevel = 5 //neighborhood-level accuracy.
	Street       GClueAccuracyLevel = 6 //Street-level accuracy.
	Exact        GClueAccuracyLevel = 8 //Exact accuracy. Typically requires GPS receiver.

)
