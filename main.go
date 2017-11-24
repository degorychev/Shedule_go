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

// type JSONTime time.Time
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

// func (t JSONTime) MarshalJSON() ([]byte, error) {
// 	//do your serializing here
// 	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RubyDate))
// 	return []byte(stamp), nil
// }

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

		rows, _ := db.Query("SELECT CONCAT(`date`,'T', `timeStart`, '-04:00') AS 'start', CONCAT(`date`,'T', `timeStop`, '-04:00') AS 'end', `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (class = ?) ORDER BY `date` ASC, `timeStart` ASC", group)

		pairs := make([]export, 0)
		for rows.Next() {
			_ = (rows.Scan(&timeStartString, &timeStopString, &discipline, &tip, &teacher, &cabinet, &subgroup))

			// var timeStart, _ = time.Parse("2006-01-02 15:04:00", timeStartString)
			// var timeStop, _ = time.Parse("2006-01-02 15:04:00", timeStopString)
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
Z
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

		rows, _ := db.Query("SELECT CONCAT(`date`,'T', `timeStart`, '-04:00') AS 'start', CONCAT(`date`,'T', `timeStop`, '-04:00') AS 'end', `discipline`, `type`, `class`, `cabinet`, `subgroup`, `teacher` FROM timetable WHERE (teacher LIKE (?)) ORDER BY `date` ASC, `timeStart` ASC", prep+"%")

		pairs := make([]export, 0)
		for rows.Next() {
			_ = (rows.Scan(&timeStartString, &timeStopString, &discipline, &tip, &group, &cabinet, &subgroup, &teacher))

			// var timeStart, _ = time.Parse("2006-01-02 15:04:00", timeStartString)
			// var timeStop, _ = time.Parse("2006-01-02 15:04:00", timeStopString)
			pairs = append(pairs, export{"null", timeStartString, timeStopString, group, discipline, tip, teacher, cabinet, subgroup})

		}

		return c.JSON(http.StatusOK, pairs) //вернуть json
	})

	// e.Logger.Fatal(e.Start(":80"))
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
