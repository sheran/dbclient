package dbclient

import (
	"fmt"
	"testing"
)

func TestSQLite(t *testing.T) {
	client, err := NewDBClient(UseSQLite("test.db"))
	if err != nil {
		panic(err)
	}
	defer client.Close()
	if err := client.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("connected")
}

func TestPostgres(t *testing.T) {
	client, err := NewDBClient(UsePostgres(
		"postgres",
		"postgres",
		"localhost",
		"",
		false,
	))
	if err != nil {
		panic(err)
	}
	defer client.Close()
	if err := client.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("connected")
	fields := []map[string]string{
		{"ID": "SERIAL PRIMARY KEY"},
		{"NAME": "varchar"}}

	if err := client.CreateTable("dertable", fields...); err != nil {
		panic(err)
	}

}
