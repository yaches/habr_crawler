package printer

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yaches/habr_crawler/models"
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

func PrintUser(user models.User) {
	fmt.Printf("Имя пользователя: %s\n", user.Username)
	fmt.Printf("Полное имя: %s\n", user.Name)
	fmt.Printf("Специализация: %s\n", user.Spec)
	fmt.Printf("О пользователе: %s\n", user.About)
	fmt.Printf("Дата рождения: %v\n", user.Birthday)
	fmt.Printf("Значки: %v\n", user.Badges)
	fmt.Printf("Хабы: %v\n", user.Hubs)
	fmt.Printf("Работа: %v\n", user.Works)
	fmt.Printf("Подписки на компании: %v\n", user.SubscribeCompanies)
	fmt.Printf("Приглашенные: %v\n", user.Invites)
	fmt.Printf("Приглашен пользователем: %v\n", user.InvitedBy)
	fmt.Printf("Карма: %v\n", user.Karma)
	fmt.Printf("Рейтинг: %v\n", user.Rating)
	fmt.Printf("Подписчики: %v\n", user.Subscribers)
	fmt.Printf("Откуда: %v\n", user.From)
	fmt.Printf("Дата регистрации: %v\n", user.RegDate)
	fmt.Printf("Количество постов: %v\n", user.PostsCount)
	fmt.Printf("Количество комментариев: %v\n", user.CommentsCount)
}

func PrintPost(post models.Post) {

}

func PrintHubHist(hist map[string]int) {
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range hist {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	fmt.Printf("\nРаспределение публикаций по интересам (Хабам):\n")
	max, dig := maxValS(hist)
	fstr := "%50s\t%" + strconv.Itoa(dig) + "d\t%s\n"
	for _, v := range ss {
		column := strings.Repeat("#", int(100.0/float64(max)*float64(v.Value)))
		fmt.Printf(fstr, v.Key, v.Value, column)
	}
}

func PrintHist(index, gran string, hist map[int]int) {
	var comment string
	switch index {
	case "posts":
		comment = "Распределение времени публикации"
	case "comments":
		comment = "Распределение времени комментирования"
	case "users":
		comment = "Распределение времени регистрации"
	}

	switch gran {
	case "hour":
		fmt.Printf("\n%s по часам в сутки:\n", comment)
		max, dig := maxVal(hist)
		fstr := "%2d%s\t%" + strconv.Itoa(dig) + "d\t%s\n"
		for i := 0; i < 24; i++ {
			column := strings.Repeat("#", int(100.0/float64(max)*float64(hist[i])))
			fmt.Printf(fstr, i, "ч", hist[i], column)
		}
	case "day":
		fmt.Printf("\n%s по дням недели:\n", comment)
		max, dig := maxVal(hist)
		fstr := "%-11s\t%" + strconv.Itoa(dig) + "d\t%s\n"
		for i := 1; i < 8; i++ {
			column := strings.Repeat("#", int(100.0/float64(max)*float64(hist[i])))
			fmt.Printf(fstr, weekdays[i], hist[i], column)
		}
	case "year":
		fmt.Printf("\n%s по году:\n", comment)
		max, dig := maxVal(hist)
		fstr := "%4d%s\t%" + strconv.Itoa(dig) + "d\t%s\n"
		for i := 2006; i < time.Now().Year()+1; i++ {
			v, ok := hist[i]
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

func maxValS(hist map[string]int) (int, int) {
	m := 0
	for _, v := range hist {
		if v > m {
			m = v
		}
	}
	return m, len(strconv.Itoa(m))
}
