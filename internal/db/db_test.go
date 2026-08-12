package db

import (
	"testing"
)

func TestOpen(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestMigrate_CreatesExpectedTables(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	for _, table := range []string{"schema_migrations", "downloads"} {
		var name string
		err := db.QueryRow(
			`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, table,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %q not found: %v", table, err)
		}
	}
}

func TestMigrate_Idempotent(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("first Migrate: %v", err)
	}
	if err := Migrate(db); err != nil {
		t.Fatalf("second Migrate: %v", err)
	}

	// Count after first run — running Migrate again must not insert duplicates.
	var countAfterFirst, countAfterSecond int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&countAfterFirst); err != nil {
		t.Fatalf("querying schema_migrations: %v", err)
	}
	if err := Migrate(db); err != nil {
		t.Fatalf("third Migrate: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&countAfterSecond); err != nil {
		t.Fatalf("querying schema_migrations: %v", err)
	}
	if countAfterSecond != countAfterFirst {
		t.Errorf("idempotency failed: count went from %d to %d after re-running Migrate", countAfterFirst, countAfterSecond)
	}
}
