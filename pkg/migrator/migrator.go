package migrator


type Migrator interface {
	MigrateUP() error
	MigrateDown() error
}