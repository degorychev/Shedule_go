package main

import (
	"net/http"
	"os"

	"./botframe"
	"./sqlque"
	"github.com/JonathanMH/goClacks/echo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
		return c.File("html/index.html")
	})

	//Все группы
	e.GET("/groups", func(c echo.Context) error {
		groups := sqlque.GetGroups()
		return c.JSON(http.StatusOK, groups) //вернуть json
	})

	//Все преподаватели
	e.GET("/teachers", func(c echo.Context) error {
		teachers := sqlque.GetTeachers()
		return c.JSON(http.StatusOK, teachers) //вернуть json
	})

	//Кто этот пользователь для мессенджера
	e.GET("/bot/who", func(c echo.Context) error {
		messenger := c.QueryParam("mess")
		id := c.QueryParam("id")

		return c.JSON(http.StatusOK, botframe.Who(messenger, id))
	})

	//Расписание для студента
	e.GET("/shedule/student/:group/today", func(c echo.Context) error {
		group := c.Param("group")
		pairs := sqlque.GetSheduleStudent(group, "today")
		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	//Расписание студента на неделю
	e.GET("/shedule/student/:group/week", func(c echo.Context) error {
		group := c.Param("group")
		pairs := sqlque.GetSheduleStudent(group, "week")
		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	//Расписание преподавателя на неделю
	e.GET("/shedule/teacher/:teacher/week", func(c echo.Context) error {
		prep := c.Param("teacher")
		pairs := sqlque.GetSheduleTeacher(prep, "week")
		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	//Расписание для преподавателя
	e.GET("/shedule/teacher/:teacher/today", func(c echo.Context) error {
		prep := c.Param("teacher")
		pairs := sqlque.GetSheduleTeacher(prep, "today")
		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	//Расписание для студента ВСЕ
	e.GET("/shedule/student/:group/all", func(c echo.Context) error {
		group := c.Param("group")
		pairs := sqlque.GetSheduleAll(group, "student")
		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	//Расписание для преподавателя ВСЕ
	e.GET("/shedule/teacher/:teacher/all", func(c echo.Context) error {
		prep := c.Param("teacher")
		pairs := sqlque.GetSheduleAll(prep, "teacher")
		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	// e.Logger.Fatal(e.Start(":80"))
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
