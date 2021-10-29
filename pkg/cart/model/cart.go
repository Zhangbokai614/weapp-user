package model

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	DBName    = "cart"
	TableName = "cart"
)

const (
	mysqlCartCreateDatabase = iota
	mysqlCartCreateTable
	mysqlCartInsert
	mysqlCartInfoByUserID
)

var (
	errInvalidMysql = errors.New("affected 0 rows")

	cartSQLString = []string{
		fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s ;`, DBName),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
			id		    	BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			user_id			BIGINT NOT NULL,
			sku_id			BIGINT NOT NULL,
			spu_id			BIGINT NOT NULL,
			count			INT NOT NULL DEFAULT 1,

			active   		BOOLEAN DEFAULT TRUE,
			created_at  	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			INDEX user_index (user_id)
		)  ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, TableName),
		fmt.Sprintf(`INSERT INTO %s.%s (user_id, sku_id, spu_id, count) VALUES (?, ?, ?, ?)`, DBName, TableName),
		fmt.Sprintf(`SELECT id, sku_id, count, active FROM %s.%s WHERE user_id = ?`, DBName, TableName),
	}
)

// CreateDatabase create cart table.
func CreateDatabase(db *sql.DB) error {
	_, err := db.Exec(cartSQLString[mysqlCartCreateDatabase])
	if err != nil {
		return err
	}

	return nil
}

// CreateTable create cart table.
func CreateCartTable(db *sql.DB) error {
	_, err := db.Exec(cartSQLString[mysqlCartCreateTable])
	if err != nil {
		return err
	}

	return nil
}

func InsertCart(db *sql.DB, userID uint32, skuID uint32, spuID uint32, count uint32) error {
	result, err := db.Exec(cartSQLString[mysqlCartInsert], userID, skuID, spuID, count)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

type CartGoods struct {
	ID     uint32
	SkuID  uint32
	Count  uint32
	Active bool
}

func InfoByUserID(db *sql.DB, userID uint32) ([]*CartGoods, error) {
	rows, err := db.Query(cartSQLString[mysqlCartInfoByUserID], userID)
	if err != nil {
		return nil, err
	}

	var result []*CartGoods
	for rows.Next() {
		var (
			id     uint32
			skuID  uint32
			count  uint32
			active bool
		)
		if err := rows.Scan(&id, &skuID, &count, &active); err != nil {
			return nil, err
		}

		result = append(result, &CartGoods{
			ID:     id,
			SkuID:  skuID,
			Count:  count,
			Active: active,
		})
	}

	return result, nil
}
