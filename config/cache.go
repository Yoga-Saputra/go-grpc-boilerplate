package config

// Cache configuration key value
// Using redis as cache driver
type cache struct {
	// Cache driver connection
	Redis redis `json:"redis"`
}
