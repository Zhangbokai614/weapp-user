package model

import (
	"database/sql"
	"fmt"
	"time"
)

const SkuTableName = "sku"

const (
	mysqlSkuCreateTable = iota
	mysqlSkuInsert
	mysqlSkuInfoBySpecAndSpuID
)

var skuSQLString = []string{
	fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
		id		    	BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		spu_id			BIGINT UNSIGNED NOT NULL,
		spec			VARCHAR(512) NOT NULL,
		price			DOUBLE NOT NULL DEFAULT 9999.99,
		stock			INT UNSIGNED NOT NULL DEFAULT 0,
		created_at  	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		INDEX spu_index (spu_id)
	)  ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, SkuTableName),
	fmt.Sprintf(`INSERT INTO %s.%s (spu_id, spec, price, stock) VALUES (?, ?, ?, ?)`, DBName, SkuTableName),
	fmt.Sprintf(`SELECT id, spu_id, spec, price, stock FROM %s.%s WHERE spu_id = ? AND spec = ?`, DBName, SkuTableName),
}

type Sku struct {
	ID        uint32    `json:"id,omitempty"`
	SpuID     uint32    `json:"spu_id,omitempty"`
	Spec      string    `json:"spec,omitempty"`
	Price     float64   `json:"price,omitempty"`
	Stock     uint32    `json:"stock,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func CreateSkuTable(db *sql.DB) error {
	_, err := db.Exec(skuSQLString[mysqlSkuCreateTable])
	if err != nil {
		return err
	}

	return nil
}

func InsertSku(db *sql.DB, spuID uint32, spec string, price float64, stock uint32) error {
	result, err := db.Exec(skuSQLString[mysqlSkuInsert], spuID, spec, price, stock)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

func InfoSkuBySpecAndSpuID(db *sql.DB, spuID uint32, spec string) (*Sku, error) {
	var sku Sku
	if err := db.QueryRow(skuSQLString[mysqlSkuInfoBySpecAndSpuID]).Scan(&sku.ID,
		&sku.SpuID, &sku.Spec, &sku.Price, &sku.Stock); err != nil {
		return nil, err
	}

	return &sku, nil
}
