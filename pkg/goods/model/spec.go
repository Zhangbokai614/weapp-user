package model

import (
	"database/sql"
	"fmt"
)

const SpecTableName = "spec"

const (
	mysqlSpecCreateTable = iota
	mysqlSpecInsert
	mysqlSpecInfoBySpuID
)

var specSQLString = []string{
	fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
		id		    	BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		spu_id			BIGINT UNSIGNED NOT NULL,
		kind 			VARCHAR(100) NOT NULL,
		value 			VARCHAR(100) NOT NULL,
		created_at  	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		INDEX spu_index (spu_id)
	)  ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, SpecTableName),
	fmt.Sprintf(`INSERT INTO %s.%s (spu_id, kind, value) VALUES (?, ?, ?)`, DBName, SpecTableName),
	fmt.Sprintf(`SELECT id, spu_id, kind, value FROM %s.%s WHERE spu_id = ?`, DBName, SpecTableName),
}

func CreateSpecTable(db *sql.DB) error {
	_, err := db.Exec(specSQLString[mysqlSpecCreateTable])
	if err != nil {
		return err
	}

	return nil
}

func InsertSpec(db *sql.DB, spuID uint32, kind string, value string) error {
	result, err := db.Exec(specSQLString[mysqlSpecInsert], spuID, kind, value)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

type Spec struct {
	ID    uint32 `json:"id,omitempty"`
	SpuID uint32 `json:"spu_id,omitempty"`
	Kind  string `json:"kind,omitempty"`
	Value string `json:"value,omitempty"`
}

func TxInfoSpecBySpuID(tx *sql.Tx, spuID uint32) ([]*Spec, error) {
	rows, err := tx.Query(specSQLString[mysqlSpecInfoBySpuID], spuID)
	if err != nil {
		return nil, err
	}

	var result []*Spec
	for rows.Next() {
		var (
			id    uint32
			spuID uint32
			kind  string
			value string
		)
		if err := rows.Scan(&id, &spuID, &kind, &value); err != nil {
			return nil, err
		}

		result = append(result, &Spec{
			ID:    id,
			SpuID: spuID,
			Kind:  kind,
			Value: value,
		})
	}

	return result, nil
}
