package main

import (
	"fmt"
	"github.com/r3labs/diff/v3"
	"strings"
)

type ClickhouseDatabaseMigrationSpec struct {
	Name       string   `json:"name"` // 资源名称
	Cluster    string   `json:"cluster"`
	Database   string   `json:"database"`
	Statements []string `json:"statements"`
}

func main() {
	// Create two instances to compare
	oldSpec := ClickhouseDatabaseMigrationSpec{
		Name:       "migration-1",
		Cluster:    "production",
		Database:   "analytics",
		Statements: []string{"CREATE TABLE test (id Int64)", "ALTER TABLE test ADD COLUMN name String"},
	}

	newSpec := ClickhouseDatabaseMigrationSpec{
		Name:       "migration-1",
		Cluster:    "production",
		Database:   "analytics", // Typo changed
		Statements: []string{"ALTER TABLE test ADD COLUMN name String", "CREATE TABLE test (id Int64)", "ALTER TABLE test ADD COLUMN age Int32"},
	}

	// Compare the two structs
	changelog, err := diff.Diff(oldSpec, newSpec)
	if err != nil {
		panic(err)
	}

	// Print the differences
	for _, change := range changelog {
		fmt.Println(change.Path)
		fmt.Printf("Field '%s' changed:\n", strings.Join(change.Path, "."))
		fmt.Printf("  Old: %v\n", change.From)
		fmt.Printf("  New: %v\n", change.To)
		fmt.Println("------")
	}
}
