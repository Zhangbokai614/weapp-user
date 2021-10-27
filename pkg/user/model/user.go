package model

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	DBName    = "user"
	TableName = "user"
)

const (
	mysqlUserCreateDatabase = iota
	mysqlUserCreateTable
	mysqlUserInsert
	mysqlUserIsExist
	mysqlUserUpdateSessionKey
	mysqlUserModifyInfo
	mysqlUserGetInfo
	mysqlUserModifyActive
	mysqlUserGetIsActive
)

var (
	errInvalidMysql = errors.New("affected 0 rows")

	userSQLString = []string{
		fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s ;`, DBName),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
			id		    	BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			openid     		VARCHAR(100) UNIQUE NOT NULL,
			session_key 	VARCHAR(100) NOT NULL,
			nick_name 		VARCHAR(100) NOT NULL DEFAULT " ",
			avatar			VARCHAR(512) NOT NULL DEFAULT " ",
			gender			TINYINT NOT NULL DEFAULT 0 COMMENT '0 unknown 1 man 2 woman',
			active   		BOOLEAN DEFAULT TRUE,
			created_at  	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, TableName),
		fmt.Sprintf(`INSERT INTO %s.%s (openid, session_key)  VALUES (?,?)`, DBName, TableName),
		fmt.Sprintf(`SELECT id FROM %s.%s WHERE openid = ? LOCK IN SHARE MODE`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET session_key = ? WHERE id = ? LIMIT 1`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET nick_name = ?, avatar = ?, gender = ? WHERE id = ? LIMIT 1`, DBName, TableName),
		fmt.Sprintf(`SELECT nick_name, avatar, gender FROM %s.%s WHERE id = ? LOCK IN SHARE MODE`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET active = ? WHERE id = ? LIMIT 1`, DBName, TableName),
		fmt.Sprintf(`SELECT active FROM %s.%s WHERE id = ? LOCK IN SHARE MODE`, DBName, TableName),
	}
)

// CreateDatabase create user table.
func CreateDatabase(db *sql.DB) error {
	_, err := db.Exec(userSQLString[mysqlUserCreateDatabase])
	if err != nil {
		return err
	}

	return nil
}

// CreateTable create user table.
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(userSQLString[mysqlUserCreateTable])
	if err != nil {
		return err
	}

	return nil
}

//CreateUser create a user
func CreateUser(db *sql.DB, openid, sessionKey string) (uint32, error) {
	result, err := db.Exec(userSQLString[mysqlUserInsert], openid, sessionKey)
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
	if err := db.QueryRow(userSQLString[mysqlUserIsExist], openid).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateSessionKey update the session for user
func UpdateSessionKey(db *sql.DB, id uint32, sessionKey string) error {
	_, err := db.Exec(userSQLString[mysqlUserUpdateSessionKey], sessionKey, id)
	if err != nil {
		return err
	}

	return nil
}

// ModifyUserActive the user updates active
func ModifyUserActive(db *sql.DB, id uint32, active bool) error {
	result, err := db.Exec(userSQLString[mysqlUserModifyActive], active, id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil

}

// IsActive return user.Active and nil if query success.
func IsActive(db *sql.DB, id uint32) (bool, error) {
	var isActive bool

	err := db.QueryRow(userSQLString[mysqlUserGetIsActive], id).Scan(&isActive)
	return isActive, err
}

// ModifyUserInfo the user updates info
func ModifyUserInfo(db *sql.DB, id uint32, nickName string, avatar string, gender int) error {
	result, err := db.Exec(userSQLString[mysqlUserModifyInfo], nickName, avatar, gender, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil

}

type UserInfo struct {
	NickName string
	Avatar   string
	Gender   int
}

func GetUserInfo(db *sql.DB, id uint32) (*UserInfo, error) {
	var (
		nickName string
		avatar   string
		gender   int
	)
	err := db.QueryRow(userSQLString[mysqlUserGetInfo], id).Scan(&nickName, &avatar, &gender)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		NickName: nickName,
		Avatar:   avatar,
		Gender:   gender,
	}, nil
}
