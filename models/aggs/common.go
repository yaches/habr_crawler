package aggs

type CommonInfo struct {
	PostsInfo    Info
	CommentsInfo Info
	UsersInfo    Info
}

type Info struct {
	Count         int
	HourHist      map[int]int
	DayOfWeekHist map[int]int
	MonthHist     map[int]int
	YearHist      map[int]int
}
