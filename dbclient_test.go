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
	dbs, err := client.ListDatabases()
	if err != nil {
		panic(err)
	}
	fmt.Println(dbs)
	// if err := client.CreateDatabase("database1"); err != nil {
	// 	panic(err)
	// }
	// dbs, err = client.ListDatabases()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(dbs)

}
