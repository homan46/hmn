package data

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"codeberg.org/rchan/hmn/config"
	"codeberg.org/rchan/hmn/constant"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(conf *config.Config) *sqlx.DB {

	db, err := sqlx.Connect("sqlite3", conf.Storage.Path)
	if err != nil {
		log.Fatal("open db fail")
	}

	version := getDBVersion(db)
	updateDB(version, db)

	return db
}

func updateDB(currentVersion int, db *sqlx.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("update database fail")
	}

	if currentVersion < 1 {
		_, err := tx.Exec(
			`
			create table system_config (
				id integer primary key autoincrement,
				key text unique,
				value json
			);
			create table user (
				id integer primary key autoincrement,
				created_time datetime not null,
				created_by	integer not null,
				modified_time datetime not null,
				modified_by integer not null,
			
				user_name text not null,
				password text not null
			);
			insert into user (
				id,
				created_time,created_by,
				modified_time,modified_by,
				user_name, password
			) values (
				` + fmt.Sprint(constant.SystemUserID) + `,
				datetime(),` + fmt.Sprint(constant.SystemUserID) + `,
				datetime(),` + fmt.Sprint(constant.SystemUserID) + `,
				"system",""
			),(
				` + fmt.Sprint(constant.AdminUserID) + `,
				datetime(),` + fmt.Sprint(constant.SystemUserID) + `,
				datetime(),` + fmt.Sprint(constant.SystemUserID) + `,
				"admin",""
			);
			create table note (
				id integer primary key autoincrement,
				created_time datetime not null,
				created_by	integer not null REFERENCES user(id),
				modified_time datetime not null,
				modified_by integer not null REFERENCES user(id),
			
				title text not null,
				content text not null,
				parent_id integer not null,
				idx integer  not null
			);
			insert into note (
				id,
				created_time,created_by,
				modified_time,modified_by,
				title ,content,
				parent_id,idx
			)values(
				1,
				datetime(),` + fmt.Sprint(constant.SystemUserID) + `,
				datetime(),` + fmt.Sprint(constant.SystemUserID) + `,
				"root","",
				0,0
			);
			
			insert into system_config (key,value) values (
				"schema.version",
				'1'
			);
			`)

		if err != nil {
			log.Fatal(err)
		}

	}
	tx.Commit()
}

// type SystemInfo struct {
// 	version int
// }

func getDBVersion(db *sqlx.DB) (version int) {
	var versionText = ""
	err := db.Get(&versionText, "select value from system_config where key = ?", "schema.version")

	log.Println(versionText)
	if err != nil {
		if strings.Contains(err.Error(), "no such table") { //TODO:
			return 0
		}
		log.Fatal(err)
	}

	ver, _ := strconv.Atoi(versionText)

	return ver
}
