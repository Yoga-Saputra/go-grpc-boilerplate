// This package is all related stuff about database using gorm.
package dbgorm

import (
	"errors"

	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/cli/dbgorm/seeder"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type (
	// AdvanceDBMigrate hold type of DB migrator using advance mode.
	AdvanceDBMigrate struct {
		Entity interface{}
		Schema interface{}
		Tx     *gorm.DB
		Source string

		AutoMigrate bool

		fn func(a *AdvanceDBMigrate) error
	}

	// Seeding hold type of DB seeding.
	Seeding struct {
		fn func(refresh bool, tx *gorm.DB) error
	}

	// Defines data struct for manual migrating
	manuallyMigrate struct {
		Entity interface{}
		Schema interface{}

		AutoMigrate bool

		source string
	}

	// Defines type of migration using raw SQL.
	migrateWithTableOpt struct {
		Entity interface{}
		Schema interface{}

		opts map[string]interface{}
	}

	// Defines type of entity or process that want to be migrating.
	migratable struct {
		key    string
		entity interface{}
	}
)

// Static map definition of table entity
var mappingEntity = map[uint]migratable{
	// Table from main DB connection
	0: {key: "wallet", entity: entity.Wallet{}},
	1: {key: "registered_service", entity: entity.RegisteredService{}},
	2: {key: "wallet_promo", entity: entity.WalletPromo{}},

	// DB function or other stuff from main DB connection
	3: {key: "create_cron_partman", entity: AdvanceDBMigrate{
		AutoMigrate: true,
		fn:          CreateCronPartman,
	}},
	4: {key: "create_setup_partition_proc", entity: AdvanceDBMigrate{
		AutoMigrate: true,
		fn:          CreateSetupPartitionProc,
	}},

	5: {key: "alter_all_proviver_log", entity: AdvanceDBMigrate{
		Entity: &entity.TransactionLogProvider{},
		fn:     AlterAllTableProviderLog,
	}},

	// Database seeding
	6: {key: "seed_registered_service", entity: Seeding{fn: (seeder.RegisteredServices{}).Run}},
}

func prepareInterface(t ...string) (dest map[uint]migratable, runAll bool, err error) {
	if len(mappingEntity) <= 0 {
		err = errors.New("given table name not found in entity mapping")
		return
	}

	// Set slice of entities from given table name(s)
	dest = map[uint]migratable{}
	for _, v := range t {
		for k, m := range mappingEntity {
			if v == m.key {
				dest[k] = m
			}
		}
	}

	// Return custom
	if len(dest) > 0 {
		return
	}

	// Return default
	dest = mappingEntity
	runAll = true
	return
}

func migrateManually(tx *gorm.DB, m manuallyMigrate, resolver ...string) error {
	if len(resolver) > 0 {
		tx = tx.Clauses(dbresolver.Use(resolver[0]))
	}

	tx = tx.Session(&gorm.Session{})
	if !tx.Migrator().HasTable(m.Entity) {
		if err := tx.Migrator().CreateTable(m.Entity); err != nil {
			return err
		}
	}

	return nil
}

func RunDBMigrate(tx *gorm.DB, seedRefresh bool, table ...string) error {
	dest, runAll, err := prepareInterface(table...)
	if err != nil {
		return err
	}

	// Keyed migrateable key and sort it
	// to guarantee the ordering
	keys := make([]int, 0, len(dest))
	for k := range dest {
		keys = append(keys, int(k))
	}

	for _, v := range keys {
		switch v := dest[uint(v)].entity.(type) {
		case manuallyMigrate:
			if !v.AutoMigrate && runAll {
				continue
			}
			if err := migrateManually(tx, v, v.source); err != nil {
				return err
			}

		case migrateWithTableOpt:
			for optK, optV := range v.opts {
				tx = tx.Set(optK, optV)
			}
			if err := tx.AutoMigrate(v.Entity); err != nil {
				return err
			}

		case AdvanceDBMigrate:
			if !v.AutoMigrate && runAll {
				continue
			}
			if v.Source != "" {
				tx = tx.Clauses(dbresolver.Use(v.Source))
			}

			v.Tx = tx
			if err := v.fn(&v); err != nil {
				return err
			}

		case Seeding:
			if err := v.fn(seedRefresh, tx); err != nil {
				return err
			}

		default:
			if err := tx.AutoMigrate(v); err != nil {
				return err
			}
		}
	}

	return nil
}
