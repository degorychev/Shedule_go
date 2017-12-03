package botframe

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
)

type (
	Whoexp struct {
		Error    string `json:"error"`
		Position string `json:"Position"`
		Value    string `json:"Value"`
		Setting  string `json:"Setting"`
	}
)

func getConnetctString() string {
	host := os.Getenv("host")
	database := os.Getenv("database")
	user := os.Getenv("user")
	pass := os.Getenv("pass")
	socket := os.Getenv("socket")

	mm := mysql.NewConfig()
	mm.Net = "tcp(" + host + ":" + socket + ")"
	mm.DBName = database
	mm.User = user
	mm.Passwd = pass

	return mm.FormatDSN()
}

//Who Кто этот клиент для бота
func Who(messenger string, id string) Whoexp {
	db, err := sql.Open("mysql", getConnetctString())
	var Position string
	var Value string
	var Setting string

	if err == nil {
		if messenger == "tele" {
			rows, err := db.Query("select `Position`, `Value`, `Setting` from users where `id`=?", id)

			if err == nil {
				if rows.Next() {
					_ = rows.Scan(&Position, &Value, &Setting)
				} else {
					return Whoexp{"no_rows", "", "", ""}
				}
				return Whoexp{"false", Position, Value, Setting}
			}
			return Whoexp{"no_query", "", "", ""}

		} else if messenger == "vk" {
			rows, err := db.Query("select `Position`, `Value`, `Setting` from usersVK where `id`=?", id)

			if err == nil {
				if rows.Next() {
					_ = rows.Scan(&Position, &Value, &Setting)
				} else {
					return Whoexp{"no_rows", "", "", ""}
				}
				return Whoexp{"false", Position, Value, Setting}
			}
			return Whoexp{"no_query", "", "", ""}
		}
		return Whoexp{"no_mess", "", "", ""}
	}
	return Whoexp{"no_sql", "", "", ""}
}
