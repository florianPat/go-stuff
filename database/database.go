package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"log"
)

type Album struct {
	Id string
	Title string
	Artist string
	Price float32
}

func main() {
	var db *sql.DB

	cfg := mysql.Config{
		User: "mariadb",
		Passwd: "mariadb",
		Net: "tcp",
		Addr: "database:3306",
		DBName: "database",
		AllowNativePasswords: true,
		Collation: "utf8_general_ci",
	}
	dataSourceName := cfg.FormatDSN()
	log.Println(dataSourceName)

	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		log.Fatal("connect", err)
	}

	db = sql.OpenDB(connector)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("ping err", pingErr)
	}

	_, err = db.Exec(`CREATE OR REPLACE TABLE album (
    				id         INT AUTO_INCREMENT NOT NULL,
    				title      VARCHAR(128) NOT NULL,
    				artist     VARCHAR(255) NOT NULL,
    				price      DECIMAL(5,2) NOT NULL,
    				PRIMARY KEY (id)
					) engine=InnoDB collate=utf8_general_ci;`)
	if err != nil {
		log.Fatal("Create table", err)
	}

	_, err = db.Exec(`
		INSERT INTO album
		  (title, artist, price)
		VALUES
		  ('Blue Train', 'John Coltrane', 56.99),
		  ('Giant Steps', 'John Coltrane', 63.99),
		  ('Jeru', 'Gerry Mulligan', 17.99),
		  ('Sarah Vaughan', 'Sarah Vaughan', 34.98);
	`)
	if err != nil {
		log.Fatal("Insert into", err)
	}

	var albums []Album

	_, err = db.Begin()
	if err != nil {
		log.Fatal("Begin tx", err)
	}

	rows, err := db.Query(`SELECT * FROM album`)
	if err != nil {
		log.Fatal("select ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err = rows.Scan(&album.Id, &album.Title, &album.Artist, &album.Price); err != nil {
			log.Fatal(err)
		}
		albums = append(albums, album)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	jsonBytesSlice, err := json.Marshal(albums)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(jsonBytesSlice)
	log.Println(jsonString)
}
