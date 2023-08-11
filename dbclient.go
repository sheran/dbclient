package dbclient

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type (
	DBClient struct {
		Type   string
		Conn   *sql.DB
		DBName string
		DBHost string
	}

	DBOption func(c *DBClient) error
)

func NewDBClient(options ...DBOption) (*DBClient, error) {
	client := &DBClient{}
	for _, option := range options {
		if err := option(client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

func UseSQLite(filename string) DBOption {
	return func(c *DBClient) error {
		db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", filename))
		if err != nil {
			return err
		}
		c.Conn = db
		c.Type = "sqlite3"
		return nil
	}
}

func UsePostgres(user, pass, host, name string, useSSL bool) DBOption {
	return func(c *DBClient) error {
		var dsn strings.Builder
		dsn.WriteString(fmt.Sprintf("user=%s ", user))
		dsn.WriteString(fmt.Sprintf("password=%s ", pass))
		dsn.WriteString(fmt.Sprintf("host=%s ", host))
		if name != "" {
			dsn.WriteString(fmt.Sprintf("dbname=%s ", name))
			c.DBName = name
		}
		if !useSSL {
			dsn.WriteString("sslmode=disable")
		}

		c.DBHost = host

		c.Type = "postgres"
		db, err := sql.Open("postgres", dsn.String())
		if err != nil {
			return err
		}
		c.Conn = db
		return nil
	}
}

func (c *DBClient) Ping() error {
	return c.Conn.Ping()
}

func (c *DBClient) ListDatabases() ([]string, error) {
	var sql string
	switch c.Type {
	case "postgres":
		sql = "SELECT datname FROM pg_database WHERE datistemplate = false;"
	}
	rows, err := c.Conn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	databases := make([]string, 0)
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}
	return databases, nil
}

func (c *DBClient) CreateDatabase(name string) error {
	var sql string
	switch c.Type {
	case "postgres":
		sql = fmt.Sprintf("CREATE DATABASE %s;", name)
	}
	_, err := c.Conn.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) CreateTable(name string, fields ...map[string]string) error {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ", name)
	sql += "("
	for i, field := range fields {
		for k, v := range field {
			sql += fmt.Sprintf("%s %s", k, v)
			if i+1 < len(fields) {
				sql += ","
			}
		}
	}
	sql += ")"
	_, err := c.Conn.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) Close() error {
	return c.Conn.Close()
}
