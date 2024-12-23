package wallet

import "time"

type (
	Meta struct {
		repo Repository
		tz   *time.Location
	}
)

// RegisterUsecase V1 wallet into the app services.
func RegisterUsecase(
	repo Repository,
	tz *time.Location,
	callback ...func(s string),
) *Meta {
	if len(callback) > 0 {
		callback[0]("Registering Wallet V1 usecase...")
	}

	m := &Meta{repo, tz}

	// Acquire public capability
	newPublicAPI(m)

	return m
}
