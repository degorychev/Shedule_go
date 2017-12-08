package sqlque

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
)

type (
	Expname struct {
		Error string `json:"error"`
		ID    string `json:"ID"`
		Name  string `json:"Naimenovanie"`
	}

	whoexp struct {
		Error    string `json:"error"`
		Position string `json:"Position"`
		Value    string `json:"Value"`
		Setting  string `json:"Setting"`
	}

	Shedule struct {
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

//GetGroups Получить список групп
func GetGroups() []Expname {
	db, _ := sql.Open("mysql", getConnetctString())
	var Naimenovanie string
	var ID string
	rows, _ := db.Query("SELECT timetable.class, groups_original.ID FROM timetable LEFT join groups_original on groups_original.Naimenovanie = timetable.class  where (date>DATE_ADD(now(), INTERVAL -31 DAY)) group by class")
	groups := make([]Expname, 0)
	for rows.Next() {
		_ = rows.Scan(&Naimenovanie, &ID)
		groups = append(groups, Expname{"", ID, Naimenovanie})
	}
	return groups
}

//GetTeachers Получить список преподавателей
func GetTeachers() []Expname {
	db, _ := sql.Open("mysql", getConnetctString())
	var Naimenovanie string
	var ID string
	rows, _ := db.Query("SELECT timetable.teacher, prepodavatel_original.ID FROM timetable LEFT join prepodavatel_original on prepodavatel_original.FIO = timetable.teacher where (date>DATE_ADD(now(), INTERVAL -31 DAY)) group by teacher")
	teachers := make([]Expname, 0)
	for rows.Next() {
		_ = rows.Scan(&Naimenovanie, &ID)
		teachers = append(teachers, Expname{"", ID, Naimenovanie})
	}
	return teachers
}

//GetSheduleStudent Получить расписание для группы УСТАРЕЛО!
func GetSheduleStudent(group string, day string) map[string][]Shedule {
	db, _ := sql.Open("mysql", getConnetctString()) //Открыть соединение с БД

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

	pairsh := map[string][]Shedule{
		"Monday":    make([]Shedule, 0),
		"Tuesday":   make([]Shedule, 0),
		"Wednesday": make([]Shedule, 0),
		"Thursday":  make([]Shedule, 0),
		"Friday":    make([]Shedule, 0),
		"Saturday":  make([]Shedule, 0),
	}
	if day == "today" {
		rows, _ := db.Query("SELECT DAYNAME(`date`) as 'week_day', `date`, `class`, `timeStart`, `timeStop`, `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (class = ?)and(date = CURDATE()) ORDER BY `date` ASC, `timeStart` ASC", group)
		for rows.Next() {
			rows.Scan(&weekD, &date, &class, &timeStart, &timeStop, &discipline, &tip, &teacher, &cabinet, &subgroup)
			pairsh[weekD] = append(pairsh[weekD], Shedule{"false", date, class, timeStart, timeStop, discipline, tip, teacher, cabinet, subgroup})
		}
		return pairsh
	} else if day == "week" {
		rows, _ := db.Query("SELECT DAYNAME(`date`) as 'week_day', `date`, `class`, `timeStart`, `timeStop`, `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (class = ?)and(year(`date`) = year(now()) and week(`date`, 0) = week(now(), 0)) ORDER BY `date` ASC, `timeStart` ASC", group)

		for rows.Next() {
			rows.Scan(&weekD, &date, &class, &timeStart, &timeStop, &discipline, &tip, &teacher, &cabinet, &subgroup)
			pairsh[weekD] = append(pairsh[weekD], Shedule{"false", date, class, timeStart, timeStop, discipline, tip, teacher, cabinet, subgroup})
		}
		return pairsh
	}
	return pairsh
}

//GetSheduleTeacher Получить расписание для преподавателя УСТАРЕЛО!
func GetSheduleTeacher(prep string, day string) map[string][]Shedule {
	db, _ := sql.Open("mysql", getConnetctString()) //Открыть соединение с БД

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

	pairsh := map[string][]Shedule{
		"Monday":    make([]Shedule, 0),
		"Tuesday":   make([]Shedule, 0),
		"Wednesday": make([]Shedule, 0),
		"Thursday":  make([]Shedule, 0),
		"Friday":    make([]Shedule, 0),
		"Saturday":  make([]Shedule, 0),
	}
	if day == "today" {
		rows, _ := db.Query("SELECT DAYNAME(`date`) as 'week_day', `date`, `class`, `timeStart`, `timeStop`, `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (teacher LIKE (?))and(date = CURDATE()) ORDER BY `date` ASC, `timeStart` ASC", prep+"%")
		for rows.Next() {
			rows.Scan(&weekD, &date, &class, &timeStart, &timeStop, &discipline, &tip, &teacher, &cabinet, &subgroup)
			pairsh[weekD] = append(pairsh[weekD], Shedule{"false", date, class, timeStart, timeStop, discipline, tip, teacher, cabinet, subgroup})
		}
		return pairsh
	} else if day == "week" {
		rows, _ := db.Query("SELECT DAYNAME(`date`) as 'week_day', `date`, `class`, `timeStart`, `timeStop`, `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (teacher LIKE (?))and(year(`date`) = year(now()) and week(`date`, 0) = week(now(), 0)) ORDER BY `date` ASC, `timeStart` ASC", prep+"%")
		for rows.Next() {
			rows.Scan(&weekD, &date, &class, &timeStart, &timeStop, &discipline, &tip, &teacher, &cabinet, &subgroup)
			pairsh[weekD] = append(pairsh[weekD], Shedule{"false", date, class, timeStart, timeStop, discipline, tip, teacher, cabinet, subgroup})
		}
		return pairsh
	}
	return pairsh
}

//GetSheduleAll Все расписание
func GetSheduleAll(value string, who string) []export {
	db, _ := sql.Open("mysql", getConnetctString())
	var timeStartString string
	var timeStopString string
	var discipline string
	var tip string
	var group string
	var teacher string
	var cabinet string
	var subgroup string
	pairs := make([]export, 0)
	if who == "student" {
		rows, _ := db.Query("SELECT CONCAT(`date`,'T', `timeStart`, '+04:00') AS 'start', CONCAT(`date`,'T', `timeStop`, '+04:00') AS 'end', `discipline`, `type`, `teacher`, `cabinet`, `subgroup` FROM timetable WHERE (class = ?) ORDER BY `date` ASC, `timeStart` ASC", value)
		for rows.Next() {
			_ = (rows.Scan(&timeStartString, &timeStopString, &discipline, &tip, &teacher, &cabinet, &subgroup))
			pairs = append(pairs, export{"null", timeStartString, timeStopString, value, discipline, tip, teacher, cabinet, subgroup}) //value это давно забытый костыль
		}
	} else if who == "teacher" {
		rows, _ := db.Query("SELECT CONCAT(`date`,'T', `timeStart`, '+04:00') AS 'start', CONCAT(`date`,'T', `timeStop`, '+04:00') AS 'end', `discipline`, `type`, `class`, `cabinet`, `subgroup`, `teacher` FROM timetable WHERE (teacher LIKE (?)) ORDER BY `date` ASC, `timeStart` ASC", value+"%")
		for rows.Next() {
			_ = (rows.Scan(&timeStartString, &timeStopString, &discipline, &tip, &group, &cabinet, &subgroup, &teacher))
			pairs = append(pairs, export{"null", timeStartString, timeStopString, group, discipline, tip, teacher, cabinet, subgroup})
		}
	}
	return pairs
}
