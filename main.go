package main

import (
  "database/sql"
  "fmt"
  "net/http"
  "os"

  _ "github.com/go-sql-driver/mysql"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
  "github.com/JonathanMH/goClacks/echo"
)

type (
	//json Структура для отправки
	groupname struct {
		Error string `json:"error"`
    	ID string `json:"ID"`
		Name string `json:"Naimenovanie"`
	}
)

func main() {
	// Echo instance
	e := echo.New()
	e.Use(goClacks.Terrify)
	
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	
	// Route => handler
	e.GET("/", func(c echo.Context) error {
		response := groupname{ID: "noID", Error: "true", Name: "Добро пожаловать в api"}
		return c.JSON(http.StatusOK, response)
	})
	
	e.GET("/groups/name/:id", func(c echo.Context) error {
		requestedID := c.Param("id") //вытащить id из запроса
		db, err := sql.Open("mysql", "egor:egor@tcp(95.104.192.212:3306)/raspisanie") //Открыть соединение с БД
		if err != nil { //в случае ошибки
			fmt.Println(err.Error())
			response := groupname{ID: "", Error: "true", Name: ""}
			return c.JSON(http.StatusInternalServerError, response)
		}
		defer db.Close() //В случае ошибки (?)
		var Naimenovanie string;
		var ID string;
		err = db.QueryRow("SELECT ID, Naimenovanie FROM groups_original WHERE ID = ?", requestedID).Scan(&ID, &Naimenovanie) //Запрос, вернет ошибку, если не удалось просканировать
		if err != nil {
			fmt.Println(err)
		}
		response := groupname{ID: ID, Error: "false", Name: Naimenovanie} //Создание нового json объекта
		return c.JSON(http.StatusOK, response)//вернуть json
	})
	
	//e.Logger.Fatal(e.Start(":4000"))
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}