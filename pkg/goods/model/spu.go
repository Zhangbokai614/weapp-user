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
	TableName = "spu"
)

const (
	mysqlSpuCreateDatabase = iota
	mysqlSpuCreateTable
	mysqlSpuInsert
	mysqlSpuInfoByCatagory
	mysqlSpuInfoRecommend
	mysqlSpuInfoByID
)

var (
	errInvalidMysql = errors.New("affected 0 rows")

	spuSQLString = []string{
		fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s ;`, DBName),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
			id		    	BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			catagory_id			BIGINT NOT NULL DEFAULT 0,
			title			VARCHAR(100) NOT NULL DEFAULT " ",
			production_code VARCHAR(100) NOT NULL DEFAULT " ",
			standard_code	VARCHAR(100) NOT NULL DEFAULT " ",
			inventory		INT NOT NULL DEFAULT 0,
			price			DOUBLE NOT NULL DEFAULT 9999.99,
			shelf_life		JSON,
			images			JSON,
			detail_images	JSON,

			recommend		BOOLEAN DEFAULT FALSE,
			active   		BOOLEAN DEFAULT TRUE,
			created_at  	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			INDEX catagory_index (catagory_id)
		)  ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, TableName),
		fmt.Sprintf(`INSERT INTO %s.%s (catagory_id, title, production_code, standard_code, inventory, 
			price, shelf_life, images, detail_images, recommend) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, DBName, TableName),
		fmt.Sprintf(`SELECT spu.id, catagory.name as catagory, spu.title, spu.images, spu.price 
			FROM %s.%s LEFT JOIN goods.catagory ON catagory.id = spu.catagory_id 
			WHERE spu.catagory_id = ? AND spu.active = true`, DBName, TableName),
		fmt.Sprintf(`SELECT spu.id, catagory.name as catagory, spu.title, spu.images, spu.price 
			FROM %s.%s LEFT JOIN goods.catagory ON catagory.id = spu.catagory_id 
			WHERE spu.recommend = true AND spu.active = true`, DBName, TableName),
		fmt.Sprintf(`SELECT id, catagory_id, title, production_code, standard_code, inventory, 
		shelf_life, images, detail_images, created_at FROM %s.%s WHERE id = ?`, DBName, TableName),
	}
)

type Spu struct {
	ID           uint32 `json:"id,omitempty"`
	CatagoryID   uint32 `json:"catagory_id,omitempty"`
	CatagoryName string `json:"catagory_name,omitempty"`

	Title          string      `json:"title,omitempty"`
	ProductionCode string      `json:"production_code,omitempty"`
	StandardCode   string      `json:"standard_code,omitempty"`
	Inventory      uint32      `json:"inventory,omitempty"`
	Price          float64     `json:"price,omitempty"`
	ShelfLife      interface{} `json:"shelf_life,omitempty"`
	Images         interface{} `json:"images,omitempty"`
	DetailImages   interface{} `json:"detail_images,omitempty"`
	Spec           []*Spec     `json:"spec,omitempty"`
	Sku            []*Sku      `json:"sku,omitempty"`
	Recommend      bool        `json:"recommend,omitempty"`
	Active         bool        `json:"active,omitempty"`
	CreatedAt      time.Time   `json:"created_at,omitempty"`
}

// CreateDatabase create user table.
func CreateDatabase(db *sql.DB) error {
	_, err := db.Exec(spuSQLString[mysqlSpuCreateDatabase])
	if err != nil {
		return err
	}

	return nil
}

// CreateTable create spu table.
func CreateSpuTable(db *sql.DB) error {
	_, err := db.Exec(spuSQLString[mysqlSpuCreateTable])
	if err != nil {
		return err
	}

	return nil
}

// TxInsertSpu add a spu
func TxInsertSpu(tx *sql.Tx, spu Spu) (uint32, error) {
	shelfLife, err := json.Marshal(spu.ShelfLife)
	if err != nil {
		return 0, err
	}

	images, err := json.Marshal(spu.Images)
	if err != nil {
		return 0, err
	}

	detailImages, err := json.Marshal(spu.DetailImages)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(spuSQLString[mysqlSpuInsert], spu.CatagoryID, spu.Title, spu.ProductionCode,
		spu.StandardCode, spu.Inventory, spu.Price, shelfLife, images, detailImages, spu.Recommend)
	if err != nil {
		return 0, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return 0, errInvalidMysql
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), err
}

// GetSpuByCatagory returns the spu belong to catagory
func GetSpuByCatagory(db *sql.DB, catagoryID uint32) ([]*Spu, error) {
	rows, err := db.Query(spuSQLString[mysqlSpuInfoByCatagory], catagoryID)
	if err != nil {
		return nil, err
	}

	var result []*Spu
	for rows.Next() {
		var (
			id           uint32
			catagoryName string
			title        string
			images       string
			price        float64
		)
		if err := rows.Scan(&id, &catagoryName, &title, &images, &price); err != nil {
			return nil, err
		}

		result = append(result, &Spu{
			ID:           id,
			CatagoryName: catagoryName,
			Title:        title,
			Images:       images,
			Price:        price,
		})
	}

	return result, nil
}

// GetSpuByCatagory returns the spu belong to catagory
func GetRecommendSpu(db *sql.DB) ([]*Spu, error) {
	rows, err := db.Query(spuSQLString[mysqlSpuInfoRecommend])
	if err != nil {
		return nil, err
	}

	var result []*Spu
	for rows.Next() {
		var (
			id           uint32
			catagoryName string
			title        string
			images       string
			price        float64
		)
		if err := rows.Scan(&id, &catagoryName, &title, &images, &price); err != nil {
			return nil, err
		}

		result = append(result, &Spu{
			ID:           id,
			CatagoryName: catagoryName,
			Title:        title,
			Images:       images,
			Price:        price,
		})
	}

	return result, nil
}

func TxInfoSpuByID(tx *sql.Tx, spuID uint32) (*Spu, error) {
	var (
		id             uint32
		catagoryID     uint32
		title          string
		productionCode string
		standardCode   string
		inventory      uint32
		shelfLife      string
		images         string
		detailImages   string
		createdAt      time.Time
	)
	if err := tx.QueryRow(spuSQLString[mysqlSpuInfoByID], spuID).Scan(&id, &catagoryID, &title, &productionCode,
		&standardCode, &inventory, &shelfLife, &images, &detailImages, &createdAt); err != nil {
		return nil, err
	}

	return &Spu{
		ID:             id,
		CatagoryID:     catagoryID,
		Title:          title,
		ProductionCode: productionCode,
		StandardCode:   standardCode,
		Inventory:      inventory,
		ShelfLife:      shelfLife,
		Images:         images,
		DetailImages:   detailImages,
		CreatedAt:      createdAt,
	}, nil
}
