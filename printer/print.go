package printer

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yaches/habr_crawler/models/aggs"
)

var weekdays = map[int]string{
	1: "Понедельник",
	2: "Вторник",
	3: "Среда",
	4: "Четверг",
	5: "Пятница",
	6: "Суббота",
	7: "Воскресенье",
}

var months = map[int]string{
	1:  "Январь",
	2:  "Февраль",
	3:  "Март",
	4:  "Апрель",
	5:  "Май",
	6:  "Июнь",
	7:  "Июль",
	8:  "Август",
	9:  "Сентябрь",
	10: "Октябрь",
	11: "Ноябрь",
	12: "Декабрь",
}

func PrintCommons(common aggs.CommonInfo) {
	m := map[string]aggs.Info{}
	m["p"] = common.PostsInfo
	m["c"] = common.CommentsInfo
	m["u"] = common.UsersInfo

	for k, info := range m {
		var comment string
		switch k {
		case "p":
			comment = "Распределение времени публикации"
			fmt.Println("\n\nПОСТЫ:")
		case "c":
			comment = "Распределение времени комментирования"
			fmt.Println("\n\nКОММЕНТАРИИ:")
		case "u":
			comment = "Распределение времени регистрации"
			fmt.Println("\n\nПОЛЬЗОВАТЕЛИ:")
		}

		fmt.Printf("\nВсего: %d\n", info.Count)

		var max, dig int
		var fstr string

		if k != "u" {
			fmt.Printf("\n%s по часам в сутки:\n", comment)
			max, dig := maxVal(info.HourHist)
			fstr := "%2d%s\t%" + strconv.Itoa(dig) + "d\t%s\n"
			for i := 0; i < 24; i++ {
				column := strings.Repeat("#", int(100.0/float64(max)*float64(info.HourHist[i])))
				fmt.Printf(fstr, i, "ч", info.HourHist[i], column)
			}
		}

		fmt.Printf("\n%s по дням недели:\n", comment)
		max, dig = maxVal(info.DayOfWeekHist)
		fstr = "%-11s\t%" + strconv.Itoa(dig) + "d\t%s\n"
		for i := 1; i < 8; i++ {
			column := strings.Repeat("#", int(100.0/float64(max)*float64(info.DayOfWeekHist[i])))
			fmt.Printf(fstr, weekdays[i], info.DayOfWeekHist[i], column)
		}

		fmt.Printf("\n%s по месяцам:\n", comment)
		max, dig = maxVal(info.MonthHist)
		fstr = "%-8s\t%" + strconv.Itoa(dig) + "d\t%s\n"
		for i := 1; i < 13; i++ {
			column := strings.Repeat("#", int(100.0/float64(max)*float64(info.MonthHist[i])))
			fmt.Printf(fstr, months[i], info.MonthHist[i], column)
		}

		fmt.Printf("\n%s по году:\n", comment)
		max, dig = maxVal(info.YearHist)
		fstr = "%4d%s\t%" + strconv.Itoa(dig) + "d\t%s\n"
		for i := 2006; i < time.Now().Year()+1; i++ {
			v, ok := info.YearHist[i]
			if ok {
				column := strings.Repeat("#", int(100.0/float64(max)*float64(v)))
				fmt.Printf(fstr, i, "г", v, column)
			}
		}
	}
}

func maxVal(hist map[int]int) (int, int) {
	m := 0
	for _, v := range hist {
		if v > m {
			m = v
		}
	}
	return m, len(strconv.Itoa(m))
}
