package mcslog

import "time"

type Meta struct {
	repo Reppository
	tz   *time.Location
}

// RegisterUsecase v1 mcslog into the app services.
func RegisterUseCase(
	repo Reppository,
	tz *time.Location,
	callback ...func(s string),
) *Meta {
	if len(callback) > 0 {
		callback[0]("Registering McsLog V1 usecase...")
	}

	m := &Meta{repo: repo, tz: tz}
	public = m
	return m
}
