package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/JonathanMH/goClacks/echo"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
) // import "time"

// "os"

type (
	groupname struct {
		Error string `json:"error"`
		ID    string `json:"ID"`
		Name  string `json:"Naimenovanie"`
	}

	teachername struct {
		Error string `json:"error"`
		ID    string `json:"ID"`
		Name  string `json:"FIO"`
	}

	whoexp struct {
		Error    string `json:"error"`
		Position string `json:"Position"`
		Value    string `json:"Value"`
		Setting  string `json:"Setting"`
	}

	shedule struct {
		Error      string `json:"error"`
		Date       string `json:"date"`
		Class      string `json:"group"`
		TimeStart  string `json:"start"`
		TimeStop   string `json:"stop"`
		Discipline string `json:"disc"`
		Tip        string `json:"type"`
		Teacher    string `json:"teacher"`
		Cabinet    string `json:"kabinet"`
		Subgroup   string `json:"subgr"`
	}

	export struct {
		Error      string `json:"error"`
		Start      string `json:"start"`
		End        string `json:"end"`
		Class      string `json:"group"`
		Discipline string `json:"title"`
		Tip        string `json:"type"`
		Teacher    string `json:"teacher"`
		Cabinet    string `json:"kabinet"`
		Subgroup   string `json:"subgr"`
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

	database := getConnetctString()

	//Все группы
	e.GET("/groups", func(c echo.Context) error {
		db, _ := sql.Open("mysql", database) //Открыть соединение с БД

		var Naimenovanie string
		var ID string

		rows, _ := db.Query("SELECT timetable.class, groups_original.ID FROM timetable LEFT join groups_original on groups_original.Naimenovanie = timetable.class  where (date>DATE_ADD(now(), INTERVAL -31 DAY)) group by class")

		groups := make([]groupname, 0)
		for rows.Next() {
			_ = rows.Scan(&Naimenovanie, &ID)
			groups = append(groups, groupname{"", ID, Naimenovanie})
		}

		return c.JSON(http.StatusOK, groups) //вернуть json
	})

	//Все преподаватели
	e.GET("/teachers", func(c echo.Context) error {
		db, _ := sql.Open("mysql", database) //Открыть соединение с БД

		var Naimenovanie string
		var ID string

		rows, _ := db.Query("SELECT timetable.teacher, prepodavatel_original.ID FROM timetable LEFT join prepodavatel_original on prepodavatel_original.FIO = timetable.teacher where (date>DATE_ADD(now(), INTERVAL -31 DAY)) group by teacher")

		teachers := make([]teachername, 0)
		for rows.Next() {
			_ = rows.Scan(&Naimenovanie, &ID)
			teachers = append(teachers, teachername{"", ID, Naimenovanie})
		}

		return c.JSON(http.StatusOK, teachers) //вернуть json
	})

	//Кто этот пользователь для мессенджера
	e.GET("/bot/who", func(c echo.Context) error {
		db, err := sql.Open("mysql", database) //Открыть соединение с БД
		messenger := c.QueryParam("mess")
		id := c.QueryParam("id")

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
						return c.JSON(http.StatusOK, whoexp{"no_rows", "", "", ""})
					}
					return c.JSON(http.StatusOK, whoexp{"false", Position, Value, Setting})
				}
				return c.JSON(http.StatusOK, whoexp{"no_query", "", "", ""})

			} else if messenger == "vk" {
				rows, err := db.Query("select `Position`, `Value`, `Setting` from usersVK where `id`=?", id)

				if err == nil {
					if rows.Next() {
						_ = rows.Scan(&Position, &Value, &Setting)
					} else {
						return c.JSON(http.StatusOK, whoexp{"no_rows", "", "", ""})
					}
					return c.JSON(http.StatusOK, whoexp{"false", Position, Value, Setting})
				}
				return c.JSON(http.StatusOK, whoexp{"no_query", "", "", ""})
			}
			return c.JSON(http.StatusOK, whoexp{"no_mess", "", "", ""})
		}
		return c.JSON(http.StatusOK, whoexp{"no_sql", "", "", ""})
	})

	//Расписание пользователя
	e.GET("/bot/shedulefor", func(c echo.Context) error {
		db, err := sql.Open("mysql", database) //Открыть соединение с БД
		messenger := c.QueryParam("mess")
		id := c.QueryParam("id")

		var Position string
		var Value string
		var Setting string

		if err == nil {
			if messenger == "tele" {
				rows, err := db.Query("select `Position`, `Value`, `Setting` from users where `id`=?", id)
				if err == nil {
					if rows.Next() {
						err = rows.Scan(&Position, &Value, &Setting)
						if err != nil {
							return c.JSON(http.StatusOK, whoexp{"no_scan", "", "", ""})
						}
					} else {
						return c.JSON(http.StatusOK, whoexp{"no_rows", "", "", ""})
					}
					if Position == "Студент" {
						rows, err := db.Query("SELECT CONCAT(`date`,'T', `timeStart`, '+04:00') AS 'start', CONCAT(`date`,'T', `timeStop`, '+04:00') AS 'end', `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (class = ?) ORDER BY `date` ASC, `timeStart` ASC", Value)
						if err == nil {
							var timeStartString string
							var timeStopString string
							var discipline string
							var tip string
							var teacher string
							var cabinet string
							var subgroup string
							pairs := make([]export, 0)
							for rows.Next() {
								err = (rows.Scan(&timeStartString, &timeStopString, &discipline, &tip, &teacher, &cabinet, &subgroup))
								if err == nil {
									pairs = append(pairs, export{"null", timeStartString, timeStopString, Value, discipline, tip, teacher, cabinet, subgroup})
								} else {
									pairs = append(pairs, export{"Scan", "", "", "", "", "", "", "", ""})
								}
							}
							return c.JSON(http.StatusOK, pairs) //вернуть json
						}
						return c.JSON(http.StatusOK, whoexp{"no_query_2", "", "", ""})
					} else if Position == "Преподаватель" {
						rows, err := db.Query("SELECT CONCAT(`date`,'T', `timeStart`, '+04:00') AS 'start', CONCAT(`date`,'T', `timeStop`, '+04:00') AS 'end', `discipline`, `type`, `class`, `cabinet`, `subgroup`, `teacher` FROM timetable WHERE (teacher LIKE (?)) ORDER BY `date` ASC, `timeStart` ASC", Value+"%")
						if err == nil {
							var timeStartString string
							var timeStopString string
							var discipline string
							var tip string
							var group string
							var cabinet string
							var subgroup string
							var teacher string
							pairs := make([]export, 0)
							for rows.Next() {
								err = (rows.Scan(&timeStartString, &timeStopString, &discipline, &tip, &group, &cabinet, &subgroup, &teacher))
								if err == nil {
									pairs = append(pairs, export{"null", timeStartString, timeStopString, Value, discipline, tip, teacher, cabinet, subgroup})
								} else {
									pairs = append(pairs, export{"Scan", "", "", "", "", "", "", "", ""})
								}
							}
							return c.JSON(http.StatusOK, pairs) //вернуть json
						}
						return c.JSON(http.StatusOK, whoexp{"no_query_2", "", "", ""})
					}
					return c.JSON(http.StatusOK, whoexp{"No_position", Position, Value, Setting})
				}
				return c.JSON(http.StatusOK, whoexp{"no_query", "", "", ""})

			} else if messenger == "vk" {
				rows, err := db.Query("select `Position`, `Value`, `Setting` from usersvk where `id`=?", id)

				if err == nil {
					if rows.Next() {
						_ = rows.Scan(&Position, &Value, &Setting)
					} else {
						return c.JSON(http.StatusOK, whoexp{"no_rows", "", "", ""})
					}
					if Position == "Студент" {
						rows, err := db.Query("SELECT CONCAT(`date`,'T', `timeStart`, '+04:00') AS 'start', CONCAT(`date`,'T', `timeStop`, '+04:00') AS 'end', `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (class = ?) ORDER BY `date` ASC, `timeStart` ASC", Value)
						if err == nil {
							var timeStartString string
							var timeStopString string
							var discipline string
							var tip string
							var teacher string
							var cabinet string
							var subgroup string
							pairs := make([]export, 0)
							for rows.Next() {
								err = (rows.Scan(&timeStartString, &timeStopString, &discipline, &tip, &teacher, &cabinet, &subgroup))
								if err == nil {
									pairs = append(pairs, export{"null", timeStartString, timeStopString, Value, discipline, tip, teacher, cabinet, subgroup})
								} else {
									pairs = append(pairs, export{"Scan", "", "", "", "", "", "", "", ""})
								}
							}
							return c.JSON(http.StatusOK, pairs) //вернуть json
						}
						return c.JSON(http.StatusOK, whoexp{"no_query_2", "", "", ""})
					} else if Position == "Преподаватель" {
						rows, err := db.Query("SELECT CONCAT(`date`,'T', `timeStart`, '+04:00') AS 'start', CONCAT(`date`,'T', `timeStop`, '+04:00') AS 'end', `discipline`, `type`, `class`, `cabinet`, `subgroup`, `teacher` FROM timetable WHERE (teacher LIKE (?)) ORDER BY `date` ASC, `timeStart` ASC", Value+"%")
						if err == nil {
							var timeStartString string
							var timeStopString string
							var discipline string
							var tip string
							var group string
							var cabinet string
							var subgroup string
							var teacher string
							pairs := make([]export, 0)
							for rows.Next() {
								err = (rows.Scan(&timeStartString, &timeStopString, &discipline, &tip, &group, &cabinet, &subgroup, &teacher))
								if err == nil {
									pairs = append(pairs, export{"null", timeStartString, timeStopString, Value, discipline, tip, teacher, cabinet, subgroup})
								} else {
									pairs = append(pairs, export{"Scan", "", "", "", "", "", "", "", ""})
								}
							}
							return c.JSON(http.StatusOK, pairs) //вернуть json
						}
						return c.JSON(http.StatusOK, whoexp{"no_query_2", "", "", ""})
					}
					return c.JSON(http.StatusOK, whoexp{"No_position", Position, Value, Setting})
				}
				return c.JSON(http.StatusOK, whoexp{"no_query", "", "", ""})
			}
			return c.JSON(http.StatusOK, whoexp{"no_mess", "", "", ""})
		}
		return c.JSON(http.StatusOK, whoexp{"no_sql", "", "", ""})
	})

	//Установить параметры пользователю
	e.POST("/bot/setsetting", func(c echo.Context) error {
		db, err := sql.Open("mysql", database) //Открыть соединение с БД
		secret := os.Getenv("secretWord")
		token := c.FormValue("token")
		messenger := c.FormValue("mess")
		id := c.FormValue("id")
		position := c.FormValue("position")
		value := c.FormValue("value")
		table := "none"
		if messenger == "tele" {
			table = "users"
		} else if messenger == "vk" {
			table = "usersvk"
		}
		if token == secret {
			if err == nil {
				_, err := db.Exec("INSERT INTO " + table + " SET `ID`='" + id + "', `Position`='" + position + "', `Value`='" + value + "' ON DUPLICATE KEY UPDATE `Position`='" + position + "', `Value`='" + value + "'")
				//_, err := db.Exec("INSERT INTO ? SET `ID`='?', `Position`='?', `Value`='?' ON DUPLICATE KEY UPDATE `Position`='?', `Value`='?'", table, id, position, value, position, value)
				if err != nil {
					return c.String(http.StatusOK, "bad_sql_insert")
				}
				return c.String(http.StatusOK, "insert_ok")
			}
			return c.String(http.StatusOK, "no_database")
		}
		return c.String(http.StatusOK, "token is not ok")

	})

	//Расписание для студента
	e.GET("/shedule/student/:group/today", func(c echo.Context) error {
		db, _ := sql.Open("mysql", database) //Открыть соединение с БД
		group := c.Param("group")

		var date string
		var class string
		var timeStart string
		var timeStop string
		var discipline string
		var tip string
		var teacher string
		var cabinet string
		var subgroup string

		rows, _ := db.Query("SELECT `date`, `class`, `timeStart`, `timeStop`, `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (class = ?)and(date = CURDATE()) ORDER BY `date` ASC, `timeStart` ASC", group)

		pairs := make([]shedule, 0)
		for rows.Next() {
			_ = rows.Scan(&date, &class, &timeStart, &timeStop, &discipline, &tip, &teacher, &cabinet, &subgroup)
			pairs = append(pairs, shedule{"false", date, class, timeStart, timeStop, discipline, tip, teacher, cabinet, subgroup})

		}

		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	//Расписание для студента ВСЕ
	e.GET("/shedule/student/:group/all", func(c echo.Context) error {
		db, _ := sql.Open("mysql", database) //Открыть соединение с БД
		group := c.Param("group")

		var timeStartString string

		var timeStopString string
		var discipline string
		var tip string
		var teacher string
		var cabinet string
		var subgroup string

		rows, _ := db.Query("SELECT CONCAT(`date`,'T', `timeStart`, '+04:00') AS 'start', CONCAT(`date`,'T', `timeStop`, '+04:00') AS 'end', `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (class = ?) ORDER BY `date` ASC, `timeStart` ASC", group)

		pairs := make([]export, 0)
		for rows.Next() {
			_ = (rows.Scan(&timeStartString, &timeStopString, &discipline, &tip, &teacher, &cabinet, &subgroup))
			pairs = append(pairs, export{"null", timeStartString, timeStopString, group, discipline, tip, teacher, cabinet, subgroup})
		}

		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	//Расписание студента на неделю
	e.GET("/shedule/student/:group/week", func(c echo.Context) error {
		db, _ := sql.Open("mysql", database) //Открыть соединение с БД
		group := c.Param("group")

		// var sheduleweek [6]shedule
		var weekD string
		var date string
		var class string
		var timeStart string
		var timeStop string
		var discipline string
		var tip string
		var teacher string
		var cabinet string
		var subgroup string

		rows, _ := db.Query("SELECT DAYNAME(`date`) as 'week_day', `date`, `class`, `timeStart`, `timeStop`, `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (class = ?)and(year(`date`) = year(now()) and week(`date`, 0) = week(now(), 0)) ORDER BY `date` ASC, `timeStart` ASC", group)

		// pairs := make([]shedule, 0)
		pairsh := map[string][]shedule{
			"Monday":    make([]shedule, 0),
			"Tuesday":   make([]shedule, 0),
			"Wednesday": make([]shedule, 0),
			"Thursday":  make([]shedule, 0),
			"Friday":    make([]shedule, 0),
			"Saturday":  make([]shedule, 0),
		}
		for rows.Next() {
			_ = rows.Scan(&weekD, &date, &class, &timeStart, &timeStop, &discipline, &tip, &teacher, &cabinet, &subgroup)
			pairsh[weekD] = append(pairsh[weekD], shedule{"false", date, class, timeStart, timeStop, discipline, tip, teacher, cabinet, subgroup})

		}

		return c.JSON(http.StatusOK, pairsh) //вернуть json
	})

	//Расписание преподавателя на неделю
	e.GET("/shedule/teacher/:teacher/week", func(c echo.Context) error {
		db, _ := sql.Open("mysql", database) //Открыть соединение с БД
		prep := c.Param("teacher")

		// var sheduleweek [6]shedule
		var weekD string
		var date string
		var class string
		var timeStart string
		var timeStop string
		var discipline string
		var tip string
		var teacher string
		var cabinet string
		var subgroup string

		rows, _ := db.Query("SELECT DAYNAME(`date`) as 'week_day', `date`, `class`, `timeStart`, `timeStop`, `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (teacher LIKE (?))and(year(`date`) = year(now()) and week(`date`, 0) = week(now(), 0)) ORDER BY `date` ASC, `timeStart` ASC", prep+"%")

		// pairs := make([]shedule, 0)
		pairsh := make(map[string][]shedule)
		for rows.Next() {
			_ = rows.Scan(&weekD, &date, &class, &timeStart, &timeStop, &discipline, &tip, &teacher, &cabinet, &subgroup)
			pairsh[weekD] = append(pairsh[weekD], shedule{"false", date, class, timeStart, timeStop, discipline, tip, teacher, cabinet, subgroup})

		}

		return c.JSON(http.StatusOK, pairsh) //вернуть json
	})

	//Расписание для преподавателя
	e.GET("/shedule/teacher/:teacher/today", func(c echo.Context) error {
		db, _ := sql.Open("mysql", database) //Открыть соединение с БД
		prep := c.Param("teacher")

		var date string
		var class string
		var timeStart string
		var timeStop string
		var discipline string
		var tip string
		var teacher string
		var cabinet string
		var subgroup string

		rows, _ := db.Query("SELECT `date`, `class`, `timeStart`, `timeStop`, `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (teacher LIKE (?))and(date = CURDATE()) ORDER BY `date` ASC, `timeStart` ASC", prep+"%")

		pairs := make([]shedule, 0)
		for rows.Next() {
			_ = rows.Scan(&date, &class, &timeStart, &timeStop, &discipline, &tip, &teacher, &cabinet, &subgroup)
			pairs = append(pairs, shedule{"false", date, class, timeStart, timeStop, discipline, tip, teacher, cabinet, subgroup})

		}

		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	//Расписание для студента ВСЕ
	e.GET("/shedule/teacher/:teacher/all", func(c echo.Context) error {
		db, _ := sql.Open("mysql", database) //Открыть соединение с БД
		prep := c.Param("teacher")

		var timeStartString string

		var timeStopString string
		var discipline string
		var tip string
		var group string
		var cabinet string
		var subgroup string
		var teacher string

		rows, _ := db.Query("SELECT CONCAT(`date`,'T', `timeStart`, '+04:00') AS 'start', CONCAT(`date`,'T', `timeStop`, '+04:00') AS 'end', `discipline`, `type`, `class`, `cabinet`, `subgroup`, `teacher` FROM timetable WHERE (teacher LIKE (?)) ORDER BY `date` ASC, `timeStart` ASC", prep+"%")

		pairs := make([]export, 0)
		for rows.Next() {
			_ = (rows.Scan(&timeStartString, &timeStopString, &discipline, &tip, &group, &cabinet, &subgroup, &teacher))
			pairs = append(pairs, export{"null", timeStartString, timeStopString, group, discipline, tip, teacher, cabinet, subgroup})

		}

		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	// e.Logger.Fatal(e.Start(":80"))
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
