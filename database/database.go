package database

import (
	"database/sql"
	"log"

	"scienzabot/consts"

	_ "github.com/mattn/go-sqlite3" //apk add --update gcc musl-dev
)

//InitDatabaseConnection occupied of establishing the first connection to the database
func InitDatabaseConnection(dbPath string) (sqlitedbStruct *SQLiteDB, err error) {
	//Initializes the connection to the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	//Check for errors
	if err != nil {
		log.Fatal(err)
	}

	//Execute the initialization query
	_, err = db.Exec(consts.InitSQLString)
	if err != nil {
		return nil, err
	}

	dbs := &SQLiteDB{db}

	return dbs, nil

	//SAMPLE QUERIES

	/*
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		for i := 0; i < 100; i++ {
			_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
			if err != nil {
				log.Fatal(err)
			}
		}
		tx.Commit()

		rows, err := db.Query("select id, name from foo")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			err = rows.Scan(&id, &name)		if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, name)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		stmt, err = db.Prepare("select name from foo where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		var name string
		err = stmt.QueryRow("3").Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(name)

		_, err = db.Exec("delete from foo")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
		if err != nil {
			log.Fatal(err)
		}

		rows, err = db.Query("select id, name from foo")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			err = rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, name)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}*/
}

/*
func (db *SQLiteDB) getSettingText(key string) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	stmt, err := tx.Prepare("select name from foo where id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var value string

	row := stmt.QueryRow()

	/*If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.*/
/*
	err = row.Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}
*/
