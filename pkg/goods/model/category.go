package model

import (
	"database/sql"
	"fmt"
)

const CatagoryTableName = "catagory"

const (
	mysqlCatagoryCreateTable = iota
	mysqlCatagoryInsert
	mysqlCatagoryInfoAll
)

var catagorySQLString = []string{
	fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
		id		    	BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		name			VARCHAR(100) UNIQUE NOT NULL DEFAULT " ",
		created_at  	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, CatagoryTableName),
	fmt.Sprintf(`INSERT INTO %s.%s(name) VALUES (?)`, DBName, CatagoryTableName),
	fmt.Sprintf(`SELECT id, name FROM %s.%s`, DBName, CatagoryTableName),
}

type Catagory struct {
	ID   uint32 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// CreateTable create catagory table.
func CreateCatagoryTable(db *sql.DB) error {
	_, err := db.Exec(catagorySQLString[mysqlCatagoryCreateTable])
	if err != nil {
		return err
	}

	return nil
}

func InsertCatagory(db *sql.DB, catagoryName string) error {
	result, err := db.Exec(catagorySQLString[mysqlCatagoryInsert], catagoryName)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

func InfoAllCatagory(db *sql.DB) ([]*Catagory, error) {
	rows, err := db.Query(catagorySQLString[mysqlCatagoryInfoAll])
	if err != nil {
		return nil, err
	}

	var result []*Catagory
	for rows.Next() {
		var (
			id   uint32
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		result = append(result, &Catagory{
			ID:   id,
			Name: name,
		})
	}

	return result, nil
}
