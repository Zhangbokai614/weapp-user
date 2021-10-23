package mysql

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	mysqlUserCreateDatabase = iota
	mysqlUserCreateTable
	mysqlUserInsert
	mysqlUserIsExist
	mysqlUserUpdateSessionKey
	mysqlUserGetIsActive
	mysqlUserModifyActive
)

const (
	DBName    = "user"
	TableName = "user"
)

var (
	errInvalidMysql = errors.New("affected 0 rows")

	adminSQLString = []string{
		fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s ;`, DBName),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
			id		    	BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			openid     		VARCHAR(512) UNIQUE NOT NULL,
			session_key 	VARCHAR(512) NOT NULL,
			active   		BOOLEAN DEFAULT TRUE,
			created_at  	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, TableName),
		fmt.Sprintf(`INSERT INTO %s.%s (openid, session_key)  VALUES (?,?)`, DBName, TableName),
		fmt.Sprintf(`SELECT id FROM %s.%s WHERE openid = ? LOCK IN SHARE MODE`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET session_key=? WHERE id = ? LIMIT 1`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET active = ? WHERE id = ? LIMIT 1`, DBName, TableName),
		fmt.Sprintf(`SELECT active FROM %s.%s WHERE id = ? LOCK IN SHARE MODE`, DBName, TableName),
	}
)

// CreateDatabase create user table.
func CreateDatabase(db *sql.DB) error {
	_, err := db.Exec(adminSQLString[mysqlUserCreateDatabase])
	if err != nil {
		return err
	}

	return nil
}

// CreateTable create user table.
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(adminSQLString[mysqlUserCreateTable])
	if err != nil {
		return err
	}

	return nil
}

//CreateUser create a user
func CreateUser(db *sql.DB, openid, sessionKey string) (uint32, error) {
	result, err := db.Exec(adminSQLString[mysqlUserInsert], openid, sessionKey)
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

	return uint32(id), nil
}

// IsExist check the user if exist
func IsExist(db *sql.DB, openid string) (uint32, error) {
	var id uint32
	if err := db.QueryRow(adminSQLString[mysqlUserIsExist], openid).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateSessionKey update the session for user
func UpdateSessionKey(db *sql.DB, id uint32, sessionKey string) error {
	result, err := db.Exec(adminSQLString[mysqlUserUpdateSessionKey], sessionKey, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

//ModifyAdminActive the administrative user updates active
func ModifyAdminActive(db *sql.DB, id uint32, active bool) error {
	result, err := db.Exec(adminSQLString[mysqlUserModifyActive], active, id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil

}

//IsActive return user.Active and nil if query success.
func IsActive(db *sql.DB, id uint32) (bool, error) {
	var (
		isActive bool
	)

	db.QueryRow(adminSQLString[mysqlUserGetIsActive], id).Scan(&isActive)
	return isActive, nil
}
