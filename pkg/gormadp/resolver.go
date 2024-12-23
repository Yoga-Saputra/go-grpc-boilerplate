package gormadp

import (
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type ResolverConfig struct {
	AdapterConfig Config
	Entity        interface{}
	Name          string
}

// RegisterResolver register multiple DB connection
func (dba *DBAdapter) RegisterResolver(cfgs []ResolverConfig) {
	// Configure the resolver
	var dialectors []gorm.Dialector
	var resolver *dbresolver.DBResolver
	for _, c := range cfgs {
		cfg := configDefault(c.AdapterConfig)
		dialectors = append(dialectors, cfg.dialector)

		if resolver == nil {
			resolver = dbresolver.Register(dbresolver.Config{
				Sources: dialectors,
			}, c.Entity, c.Name)
		} else {
			resolver = resolver.Register(dbresolver.Config{
				Sources: dialectors,
			}, c.Entity, c.Name)
		}
	}

	// Use configured resolver
	dba.DB.Use(resolver)
}
