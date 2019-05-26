package crawler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/yaches/habr_crawler/models"

	"github.com/PuerkitoBio/goquery"
)

var months map[string]string = map[string]string{
	"января":   "01",
	"февраля":  "02",
	"марта":    "03",
	"апреля":   "04",
	"мая":      "05",
	"июня":     "06",
	"июля":     "07",
	"августа":  "08",
	"сентября": "09",
	"октября":  "10",
	"ноября":   "11",
	"декабря":  "12",
}

func parseTime(t string) (time.Time, error) {
	return time.Parse("2006-01-2T15:04 MST", t)
}

func parseRusTime(t string) (time.Time, error) {
	sl := strings.Split(t, " ")
	if len(sl) != 5 {
		return time.Time{}, errors.New("Invalid date format: " + t)
	}
	t = fmt.Sprintf("%s-%s-%sT%s MSK", sl[2], months[sl[1]], sl[0], sl[4])
	return parseTime(t)
}

func parsePost(r io.Reader) (models.Post, []models.Comment, error) {
	post := models.Post{}
	comments := []models.Comment{}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return post, comments, err
	}
	sel := doc.Find("*")

	// Get author
	post.Author = sel.Find(".post__meta .user-info__nickname").Text()
	if post.Author == "" {
		log.Println("Can't get post author")
	}

	// Get post pub time
	t, err := parseRusTime(sel.Find(".post__meta .post__time").Text())
	if err != nil {
		log.Println("Can't get post pub date", err)
	}
	post.PubDate = t

	// Get post title
	post.Title = sel.Find(".post__title-text").Text()
	if post.Title == "" {
		log.Println("Can't get post title")
	}

	// Get post hubs
	post.Hubs = []string{}
	sel.Find(".inline-list__item_hub a").Each(func(i int, s *goquery.Selection) {
		if t := s.Text(); t != "" {
			post.Hubs = append(post.Hubs, t)
		}
	})

	// Get post tags
	post.Tags = []string{}
	sel.Find(".inline-list_fav-tags a").Each(func(i int, s *goquery.Selection) {
		if t := s.Text(); t != "" {
			post.Tags = append(post.Tags, t)
		}
	})

	// Get post body
	post.Text = sel.Find(".post__text").Text()
	if post.Text == "" {
		log.Println("Can't get post text")
	}

	// Get comments
	sel.Find(".content-list__item_comment").Each(func(i int, s *goquery.Selection) {
		comment := models.Comment{PostID: post.ID}
		id, ok := s.First().Attr("rel")
		comment.ID = id
		if !ok {
			log.Println("Can't get comment id")
		}

		parent, ok := s.Find(".parent_id").First().Attr("data-parent_id")
		comment.ParentID = parent
		if !ok {
			log.Println("Can't get comment parent id")
		}

		comment.Author = s.Find(".user-info__nickname_comment").First().Text()
		if comment.Author == "" {
			log.Println("Can't get comment author")
		}

		t, err := parseRusTime(s.Find(".comment__date-time_published").First().Text())
		if err != nil {
			log.Println("Can't get comment pub date", err)
		}
		comment.PubDate = t

		comment.Text = s.Find(".comment__message").First().Text()
		if comment.Text == "" {
			log.Println("Can't get comment text")
		}

		comments = append(comments, comment)
	})

	post.CommentsCount = len(comments)
	return post, comments, nil
}
