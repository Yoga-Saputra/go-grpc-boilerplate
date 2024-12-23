package boot

import (
	"github.com/Yoga-Saputra/go-grpc-boilerplate/app/entity"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/gormadp"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database Adapter pointer value
var DBA *gormadp.DBAdapter

// Start database connection
func dbUp(args *AppArgs) {
	var loglevel logger.LogLevel
	if config.Of.App.Debug() {
		loglevel = logger.Info
	} else {
		loglevel = logger.Silent
	}

	pkgOptions := &gorm.Config{
		Logger: logger.Default.LogMode(loglevel),
	}

	cfg := gormadp.Config{
		Host:     config.Of.Database.Host,
		Port:     config.Of.Database.Port,
		User:     config.Of.Database.User,
		Password: config.Of.Database.Password,
		DBName:   config.Of.Database.Name,
		Dialect:  gormadp.Dialect(config.Of.Database.Dialect),
		Options:  pkgOptions,
	}
	opts := cfg.Dialect.PgOptions(gormadp.PgConfig{
		SSLMode:  false,
		TimeZone: "Asia/Manila",
	})

	dba := gormadp.Open(cfg, opts)

	// Register other DB connection sources to DB resolver
	var resolverSourcesCfg []gormadp.ResolverConfig
	for _, s := range config.Of.Database.Resolver.Sources {
		mapp := resolverEntityMapper(s.Identifier)

		resolverSourcesCfg = append(resolverSourcesCfg, gormadp.ResolverConfig{
			AdapterConfig: gormadp.Config{
				Host:     s.Host,
				Port:     s.Port,
				User:     s.User,
				Password: s.Password,
				DBName:   s.Name,
				Dialect:  gormadp.Dialect(s.Dialect),
				Options:  pkgOptions,
			},
			Entity: mapp.Ent,
			Name:   mapp.Nm,
		})
	}
	dba.RegisterResolver(resolverSourcesCfg)

	DBA = dba
	printOutUp("New DB connection successfully open")
}

// Stop database connection
func dbDown() {
	printOutDown("Closing current DB connection...")
	DBA.Close()
	DBA = nil
}

func resolverEntityMapper(identifier string) (res struct {
	Ent interface{}
	Nm  string
}) {
	switch identifier {
	// case "stake_log":
	// 	res = struct {
	// 		Ent interface{}
	// 		Nm  string
	// 	}{Ent: &entity.TransactionLogProvider{}, Nm: conDBName}

	case "mcs_log":
		res = struct {
			Ent interface{}
			Nm  string
		}{Ent: &entity.Transfer{}, Nm: identifier}
	}

	return
}
