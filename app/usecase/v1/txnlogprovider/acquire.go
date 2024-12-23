package txnlogprovider

type (
	Meta struct {
		repo Repository
	}
)

// RegisterUsecase V1 txnlogprovider into the app services.
func RegisterUsecase(
	repo Repository,
	callback ...func(s string),
) *Meta {
	if len(callback) > 0 {
		callback[0]("Registering TransactionLogProvider V1 usecase...")
	}

	m := &Meta{repo}

	// Acquire public capability
	newPublicAPI(m)

	return m
}
