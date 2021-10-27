package model

import (
	"database/sql"
	"fmt"
)

const KindTableName = "kind"

const (
	mysqlKindCreateTable = iota
	mysqlKindInsert
	mysqlKindGetAll
)

var kindSQLString = []string{
	fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
		id		    	BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		name			VARCHAR(100) UNIQUE NOT NULL DEFAULT " ",
		PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, KindTableName),
	fmt.Sprintf(`INSERT INTO %s.%s(name) VALUES (?)`, DBName, KindTableName),
	fmt.Sprintf(`SELECT id, name FROM %s.%s`, DBName, KindTableName),
}

type Kind struct {
	ID   uint32 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// CreateTable create kind table.
func CreateKindTable(db *sql.DB) error {
	_, err := db.Exec(kindSQLString[mysqlKindCreateTable])
	if err != nil {
		return err
	}

	return nil
}

func InsertKind(db *sql.DB, kindName string) error {
	result, err := db.Exec(kindSQLString[mysqlKindInsert], kindName)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

func GetAllKind(db *sql.DB) ([]*Kind, error) {
	rows, err := db.Query(kindSQLString[mysqlKindGetAll])
	if err != nil {
		return nil, err
	}

	var result []*Kind
	for rows.Next() {
		var (
			id   uint32
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		result = append(result, &Kind{
			ID:   id,
			Name: name,
		})
	}

	return result, nil
}
