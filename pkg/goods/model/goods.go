package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	DBName    = "goods"
	TableName = "goods"
)

const (
	mysqlGoodsCreateDatabase = iota
	mysqlGoodsCreateTable
	mysqlGoodsInsert
	mysqlGoodsInfoByKind
	mysqlGoodsInfoRecommend
)

var (
	errInvalidMysql = errors.New("affected 0 rows")

	goodsSQLString = []string{
		fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s ;`, DBName),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
			id		    	BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			kind_id			BIGINT NOT NULL DEFAULT 0,
			title			VARCHAR(100) NOT NULL DEFAULT " ",
			production_code VARCHAR(100) NOT NULL DEFAULT " ",
			standard_code	VARCHAR(100) NOT NULL DEFAULT " ",
			inventory		INT NOT NULL DEFAULT 0,
			price			DOUBLE NOT NULL DEFAULT 999999,
			shelf_life		JSON,
			images			JSON,
			detailImages	JSON,

			recommend		BOOLEAN DEFAULT FALSE,
			active   		BOOLEAN DEFAULT TRUE,
			created_at  	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			INDEX kind_index (kind_id)
		)  ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, TableName),
		fmt.Sprintf(`INSERT INTO %s.%s (kind_id, title, production_code, standard_code, inventory, price, shelf_life, images, detailImages, recommend) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, DBName, TableName),
		`SELECT goods.id, kind.name as kind, goods.title, goods.images, goods.price 
		FROM goods.goods LEFT JOIN goods.kind ON kind.id = goods.kind_id WHERE goods.kind_id = ? AND goods.active = true`,
		`SELECT goods.id, kind.name as kind, goods.title, goods.images, goods.price 
		FROM goods.goods LEFT JOIN goods.kind ON kind.id = goods.kind_id WHERE goods.recommend = true AND goods.active = true`,
	}
)

type Goods struct {
	ID       uint32 `json:"id,omitempty"`
	KindID   uint32 `json:"kind_id,omitempty"`
	KindName string `json:"kind_name,omitempty"`

	Title          string      `json:"title,omitempty"`
	ProductionCode string      `json:"production_code,omitempty"`
	StandardCode   string      `json:"standard_code,omitempty"`
	Inventory      uint32      `json:"inventory,omitempty"`
	Price          float64     `json:"price,omitempty"`
	ShelfLife      interface{} `json:"shelf_life,omitempty"`
	Images         interface{} `json:"images,omitempty"`
	DetailImages   interface{} `json:"detail_images,omitempty"`

	Recommend bool      `json:"recommend,omitempty"`
	Active    bool      `json:"active,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// CreateDatabase create user table.
func CreateDatabase(db *sql.DB) error {
	_, err := db.Exec(goodsSQLString[mysqlGoodsCreateDatabase])
	if err != nil {
		return err
	}

	return nil
}

// CreateTable create goods table.
func CreateGoodsTable(db *sql.DB) error {
	_, err := db.Exec(goodsSQLString[mysqlGoodsCreateTable])
	if err != nil {
		return err
	}

	return nil
}

// InsertGoods add a goods
func InsertGoods(db *sql.DB, goods Goods) error {
	shelfLife, err := json.Marshal(goods.ShelfLife)
	if err != nil {
		return err
	}

	images, err := json.Marshal(goods.Images)
	if err != nil {
		return err
	}

	detailImages, err := json.Marshal(goods.DetailImages)
	if err != nil {
		return err
	}

	result, err := db.Exec(goodsSQLString[mysqlGoodsInsert], goods.KindID, goods.Title, goods.ProductionCode,
		goods.StandardCode, goods.Inventory, goods.Price, shelfLife, images, detailImages, goods.Recommend)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// GetGoodsByKind returns the goods belong to kind
func GetGoodsByKind(db *sql.DB, kindID uint32) ([]*Goods, error) {
	rows, err := db.Query(goodsSQLString[mysqlGoodsInfoByKind], kindID)
	if err != nil {
		return nil, err
	}

	var result []*Goods
	for rows.Next() {
		var (
			id       uint32
			kindName string
			title    string
			images   string
			price    float64
		)
		if err := rows.Scan(&id, &kindName, &title, &images, &price); err != nil {
			return nil, err
		}

		result = append(result, &Goods{
			ID:       id,
			KindName: kindName,
			Title:    title,
			Images:   images,
			Price:    price,
		})
	}

	return result, nil
}

// GetGoodsByKind returns the goods belong to kind
func GetRecommendGoods(db *sql.DB) ([]*Goods, error) {
	rows, err := db.Query(goodsSQLString[mysqlGoodsInfoRecommend])
	if err != nil {
		return nil, err
	}

	var result []*Goods
	for rows.Next() {
		var (
			id       uint32
			kindName string
			title    string
			images   string
			price    float64
		)
		if err := rows.Scan(&id, &kindName, &title, &images, &price); err != nil {
			return nil, err
		}

		result = append(result, &Goods{
			ID:       id,
			KindName: kindName,
			Title:    title,
			Images:   images,
			Price:    price,
		})
	}

	return result, nil
}
