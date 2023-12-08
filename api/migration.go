package api

import (
	"fmt"
	"strings"
)

type Migration struct {
	ID     int
	Table  string
	Schema string
}

//////////////////////////
// Type Conversion
// (TODO :: More precise type conversion)
//////////////////////////

var SQL_TYPES map[string]string = map[string]string{
	"bool":    "BOOLEAN",
	"string":  "VARCHAR",
	"uint8":   "INTEGER",
	"uint16":  "INTEGER",
	"uint32":  "INTEGER",
	"uint64":  "INTEGER",
	"int8":    "INTEGER",
	"int16":   "INTEGER",
	"int32":   "INTEGER",
	"int64":   "INTEGER",
	"int":     "INTEGER",
	"float32": "FLOAT",
	"float64": "FLOAT",
	"float":   "FLOAT",
	"Time":    "DATETIME",
}

func MigrationCreate(model Model) Migration {
	var schema strings.Builder

	schema.WriteString("CREATE TABLE ")
	schema.WriteString(model.Name)
	schema.WriteString(" (")

	for idx, field := range model.Fields {
		fieldType := field.Type
		if fieldType == "VARCHAR" {
			fieldType += fmt.Sprintf("(%v)", field.Properties.MaxLength)
		}

		schema.WriteString(fmt.Sprintf("%s %s", field.Name, fieldType))

		if field.Properties.PrimaryKey {
			schema.WriteString(" PRIMARY KEY AUTOINCREMENT")
		} else if field.Properties.BelongsTo != nil {
			belongsTo := field.Properties.BelongsTo

			belongsToID := field.Properties.BelongsTo.GetPrimaryField()
			if belongsToID == nil {
				panic(fmt.Errorf("Model '%v' belongs to '%v' but doesn't have a primary key", model.Name, belongsTo.Name))
			}

			schema.WriteString(fmt.Sprintf(" REFERENCES %v(%v)", belongsTo.Name, belongsToID.Name))
		}

		if len(field.Properties.Default) > 0 {
			schema.WriteString(fmt.Sprintf(" DEFAULT %v", field.Properties.Default))
		}

		if idx < len(model.Fields)-1 {
			schema.WriteString(", ")
		}
	}
	schema.WriteString(")")

	return Migration{
		ID:     model.ID,
		Table:  model.Name,
		Schema: schema.String(),
	}
}

func (migration Migration) QuerySchema() (string, error) {
	stmt, err := Project.DataBase.Prepare("SELECT sql FROM sqlite_schema WHERE name = ?")
	if err != nil {
		return "", err
	}

	row := stmt.QueryRow(migration.Table)
	if err != nil {
		return "", err
	}

	var schema string
	err = row.Scan(&schema)
	if err != nil {
		return "", err
	}

	return schema, nil
}

func (migration Migration) Create() error {
	_, err := Project.DataBase.Exec(migration.Schema)
	if err != nil {
		return err
	}
	return nil
}

func (migration Migration) Drop() error {
	_, err := Project.DataBase.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %v", migration.Table))
	if err != nil {
		return err
	}
	return nil
}

func (migration Migration) GetFile() string {
	return fmt.Sprintf("%s/%v-%s.sql", Project.MigrationsDirectory, migration.ID, migration.Table)
}
