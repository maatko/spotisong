package api

import (
	"fmt"
	"strings"
)

type Migration struct {
	Index int
	Table string
	Query string
}

// used for managning the database, every migration
// is tied with a model so changes to the model affect
// the migration and there for will affect the table in the database
var MigrationRegistry map [string] Migration = map [string] Migration {}

func CreateMigration(model *ModelInformation) error {
	var builder strings.Builder

	// TODO :: Dropping the table every single time the
	// migration is ran is not a good idea, changes need
	// to be detect based on the state of the database
	// and if any changes were detected the the table
	// needs to be dropped
	builder.WriteString("DROP TABLE IF EXISTS ")
	builder.WriteString(model.Name)
	builder.WriteString("; ")

	builder.WriteString("CREATE TABLE IF NOT EXISTS ")
	builder.WriteString(model.Name)
	builder.WriteString(" (")

	fieldCount := len(model.Fields)
	for i := 0; i < fieldCount; i++ {
		field := model.Fields[i]
		properties := field.Properties

		builder.WriteString(field.Name)
		builder.WriteString(" ")
		builder.WriteString(field.Type)

		if field.Type == "VARCHAR" {
			builder.WriteString(fmt.Sprintf("(%v)", properties.MaxLength))
		}

		if properties.AutoIncrement {
			builder.WriteString(" PRIMARY KEY AUTOINCREMENT")
		} else if properties.BelongsTo != nil {
			info := properties.BelongsTo
			for i := 0; i < len(info.Fields); i++ {
				field := info.Fields[i]
				if field.Properties.AutoIncrement {
					builder.WriteString(fmt.Sprintf(" REFERENCES %v(%v)", info.Name, field.Name))
					break
				}
			}
		}

		if len(properties.Default) > 0 {
			builder.WriteString(" DEFAULT ")
			builder.WriteString(properties.Default)
		}

		if i < fieldCount - 1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(")")

	MigrationRegistry[model.Name] = Migration {
		Index: model.Index,
		Table: model.Name,
		Query: builder.String(),
	}

	return nil
}