package txnlog

type (
	Meta struct {
		repo Repository
	}
)

// RegisterUsecase V1 txnlog into the app services.
func RegisterUsecase(
	repo Repository,
	callback ...func(s string),
) *Meta {
	if len(callback) > 0 {
		callback[0]("Registering TransactionLogInternal V1 usecase...")
	}

	m := &Meta{repo}

	// Acquire public capability
	newPublicAPI(m)

	return m
}
