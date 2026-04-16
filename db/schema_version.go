// schema_version.go — Query helpers for SchemaVersion table.
package db

// SchemaVersionEntry represents a single migration record.
type SchemaVersionEntry struct {
	Description string
	AppliedAt   string
	Version     int
}

// ListSchemaVersions returns all applied migrations ordered by version.
func (d *DB) ListSchemaVersions() ([]SchemaVersionEntry, error) {
	rows, err := d.Query("SELECT Version, Description, AppliedAt FROM SchemaVersion ORDER BY Version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []SchemaVersionEntry
	for rows.Next() {
		var e SchemaVersionEntry
		if err := rows.Scan(&e.Version, &e.Description, &e.AppliedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
