package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type MigrationColumn struct {
	Name string
	Type string
	Properties string
}

type Migration struct {
	Table string
	Columns [] MigrationColumn
}

// used for storing all the migrations
// so the can later be accessed
var Migrations map[string] Migration = make(map [string] Migration)

func Create(table string, model any) {
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)

	var columns [] MigrationColumn

	for i := 0; i < modelType.NumField(); i++ {
		fieldType := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		var typeBuilder strings.Builder
		if (fieldType.Type.Kind() == reflect.String) {
			typeBuilder.WriteString(fmt.Sprintf("VARCHAR(%v)", fieldValue.Interface()))
		} else {
			typeBuilder.WriteString("INTEGER")
		}

		column := MigrationColumn {
			Name: strings.ToLower(fieldType.Name),
			Type: typeBuilder.String(),
		}

		if properties, ok := fieldType.Tag.Lookup("properties"); ok {
			column.Properties = properties
		}
		
		columns = append(columns, column)
	}

	// append the migration
	Migrations[table] = Migration {
		Table: table,
		Columns: columns,
	}
}

func Migrate(database *sql.DB) {
	for table, migration := range Migrations {
		var existingColumns [] MigrationColumn

		rows, err := database.Query(fmt.Sprintf("PRAGMA table_info(%v)", table))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var cid int
			var name string
			var dataType string
			var notNull int
			var defaultValue sql.NullString
			var primaryKey int
	
			err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey)
			if err != nil {
				log.Fatal(err)
			}
			
			existingColumns = append(existingColumns, MigrationColumn {
				name,
				dataType,
				"",
			})
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		if len(migration.Columns) != len(existingColumns) {
			migration.Migrate(database)
			return
		}

		for index, column := range migration.Columns {
			existingColumn := existingColumns[index]
			if column.Name != existingColumn.Name || column.Type != existingColumn.Type {
				migration.Migrate(database)
				return
			}
		}
	}
}

func (migration Migration) Migrate(database *sql.DB) {
	// recreate the table
	database.Exec(migration.GenerateDeletionSQL())
	database.Exec(migration.GenerateCreationSQL())

	log.Printf("Migrating the '%v' table...\n", migration.Table)
}

func (migration Migration) GenerateDeletionSQL() string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %v", migration.Table)
}

func (migration Migration) GenerateCreationSQL() string {
	var builder strings.Builder

	builder.WriteString("CREATE TABLE IF NOT EXISTS ")
	builder.WriteString(migration.Table)

	if len(migration.Columns) == 0 {
		return builder.String()
	}

	builder.WriteString(" (")
	for index, column := range migration.Columns {
		builder.WriteString(fmt.Sprintf("%v %v", column.Name, column.Type))
		if len(column.Properties) > 0 {
			builder.WriteString(" ")
			builder.WriteString(column.Properties)
		}
		if index < len(migration.Columns) - 1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(")")
	return builder.String()
}