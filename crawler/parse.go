package crawler

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/yaches/habr_crawler/models"
	"go.uber.org/zap"

	"github.com/PuerkitoBio/goquery"
)

var months map[string]string = map[string]string{
	"января":   "January",
	"февраля":  "February",
	"марта":    "March",
	"апреля":   "April",
	"мая":      "May",
	"июня":     "June",
	"июля":     "July",
	"августа":  "August",
	"сентября": "September",
	"октября":  "October",
	"ноября":   "November",
	"декабря":  "December",
}

func parseRusTime(t string) (time.Time, error) {
	sl := strings.Split(t, " ")

	if len(sl) != 5 && len(sl) != 3 {
		return time.Time{}, errors.New("Invalid date format: " + t)
	}

	if len(sl) == 3 {
		date := time.Time{}
		switch sl[0] {
		case "сегодня":
			date = time.Now()
		case "вчера":
			date = time.Now().AddDate(0, 0, -1)
		default:
			return date, errors.New("Invalid date format: " + t)
		}
		y, m, d := date.Date()
		t = fmt.Sprintf("%v-%v-%vT%s MSK", y, m, d, sl[2])
	} else {
		t = fmt.Sprintf("%s-%s-%sT%s MSK", sl[2], months[sl[1]], sl[0], sl[4])
	}
	return time.Parse("2006-January-2T15:04 MST", t)
}

func parseRusDay(t string) (time.Time, error) {
	sl := strings.Split(t, " ")
	if len(sl) != 4 {
		return time.Time{}, errors.New("Invalid birthday format: " + t)
	}
	t = fmt.Sprintf("%s-%s-%s", sl[2], months[sl[1]], sl[0])
	return time.Parse("2006-January-2", t)
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
		zap.L().Warn("Can't get post author")
	}

	// Get post pub time
	t, err := parseRusTime(sel.Find(".post__meta .post__time").Text())
	if err != nil {
		zap.L().Warn("Can't get post pub date", zap.Error(err))
	}
	post.PubDate = t

	// Get post title
	post.Title = sel.Find(".post__title-text").Text()
	if post.Title == "" {
		zap.L().Warn("Can't get post title")
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
		zap.L().Warn("Can't get post text")
	}

	// Get rating
	rate, err := strconv.ParseInt(strings.ReplaceAll(sel.Find(".post-additionals").Find(".voting-wjt__counter").Text(), "–", "-"), 10, 32)
	if err != nil {
		zap.L().Warn("Can't parse post rating", zap.Error(err))
	}
	post.Rating = int(rate)

	// Get comments
	sel.Find(".content-list__item_comment").Each(func(i int, s *goquery.Selection) {
		comment := models.Comment{PostID: post.ID}
		id, ok := s.First().Attr("rel")
		comment.ID = id
		if !ok {
			zap.L().Warn("Can't get comment id")
		}

		parent, ok := s.Find(".parent_id").First().Attr("data-parent_id")
		comment.ParentID = parent
		if !ok {
			zap.L().Warn("Can't get comment parent id")
		}

		comment.Author = s.Find(".user-info__nickname_comment").First().Text()
		if comment.Author == "" {
			zap.L().Warn("Can't get comment author")
		}

		t, err := parseRusTime(s.Find(".comment__date-time_published").First().Text())
		if err != nil {
			zap.L().Warn("Can't get comment pub date", zap.Error(err))
		}
		comment.PubDate = t

		rate, err := strconv.ParseInt(strings.ReplaceAll(sel.Find(".voting-wjt_comments .voting-wjt__counter").First().Text(), "–", "-"), 10, 32)
		if err != nil {
			zap.L().Warn("Can't parse comment rating", zap.Error(err))
		}
		comment.Rating = int(rate)

		comment.Text = s.Find(".comment__message").First().Text()
		if comment.Text == "" {
			zap.L().Warn("Can't get comment text")
		}

		comments = append(comments, comment)
	})

	post.CommentsCount = len(comments)
	return post, comments, nil
}

func parseUser(r io.Reader) (models.User, error) {
	user := models.User{}
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return user, err
	}
	sel := doc.Find("*")

	user.Name = sel.Find(".user-info__fullname").Text()
	user.Spec = sel.Find(".user-info__specialization").Text()
	user.About = sel.Find(".profile-section__about-text").Text()

	sel.Find(".defination-list__item_profile-summary").Each(func(i int, s *goquery.Selection) {
		label := s.Find(".defination-list__label_profile-summary").Text()
		switch label {
		case "Откуда":
			user.From = []string{}
			s.Find(".defination-list__value a").Each(func(i int, s *goquery.Selection) {
				user.From = append(user.From, s.Text())
			})
		case "Дата рождения":
			t, err := parseRusDay(s.Find(".defination-list__value").Text())
			if err != nil {
				zap.L().Warn("Can't parse user birthday")
			}
			user.Birthday = t
		case "Зарегистрирован":
			t, err := parseRusDay(s.Find(".defination-list__value").Text())
			if err != nil {
				zap.L().Warn("Can't parse user reg date")
			}
			user.RegDate = t
		case "Работает в":
			user.Works = []string{}
			s.Find(".defination-list__value a").Each(func(i int, s *goquery.Selection) {
				user.Works = append(user.Works, s.Text())
			})
		}
	})

	karma, err := strconv.ParseFloat(strings.ReplaceAll(sel.Find(".stacked-counter__value_green").Text(), ",", "."), 32)
	if err != nil {
		zap.L().Warn("Can't parse user karma", zap.Error(err))
	}
	user.Karma = float32(karma)

	rate, err := strconv.ParseFloat(strings.ReplaceAll(sel.Find(".stacked-counter_rating .stacked-counter__value").Text(), ",", "."), 32)
	if err != nil {
		zap.L().Warn("Can't parse user rating", zap.Error(err))
	}
	user.Rating = float32(rate)

	subs, err := strconv.ParseInt(sel.Find(".stacked-counter_subscribers .stacked-counter__value").Text(), 10, 32)
	if err != nil {
		zap.L().Warn("Can't parse user subscribers count", zap.Error(err))
	}
	user.Subscribers = int(subs)

	sel.Find(".tabs-level_top .tabs-menu__item-counter_total").Each(func(i int, s *goquery.Selection) {
		title, ok := s.Attr("title")
		if !ok {
			zap.L().Warn("Can't get user pubs or comments count")
			return
		}
		if strings.Contains(title, "Публикации") {
			pubs := strings.TrimPrefix(title, "Публикации: ")
			pubsCnt, err := strconv.ParseInt(pubs, 10, 32)
			if err != nil {
				zap.L().Warn("Can't parse user publications count", zap.Error(err))
				return
			}
			user.PostsCount = int(pubsCnt)
		}
		if strings.Contains(title, "Комментарии") {
			comms := strings.TrimPrefix(title, "Комментарии: ")
			commsCnt, err := strconv.ParseInt(comms, 10, 32)
			if err != nil {
				zap.L().Warn("Can't parse user comments count", zap.Error(err))
				return
			}
			user.CommentsCount = int(commsCnt)

		}
	})

	user.Badges = []string{}
	sel.Find(".profile-section__user-badge").Each(func(i int, s *goquery.Selection) {
		b := s.Text()
		if b != "" {
			user.Badges = append(user.Badges, s.Text())
		}
	})

	user.Hubs = []string{}
	sel.Find(".profile-section__user-hub").Each(func(i int, s *goquery.Selection) {
		h := s.Text()
		if h != "" {
			user.Hubs = append(user.Hubs, h)
		}
	})

	user.Invites = []string{}
	sel.Find(".content-list__item_invited-users .list-snippet__nickname").Each(func(i int, s *goquery.Selection) {
		u := s.Text()
		if u != "" {
			user.Invites = append(user.Invites, u)
		}
	})

	user.SubscribeCompanies = []string{}
	sel.Find(".content-list__item_companies .list-snippet__title-link").Each(func(i int, s *goquery.Selection) {
		c := s.Text()
		if c != "" {
			user.SubscribeCompanies = append(user.SubscribeCompanies, c)
		}
	})

	user.InvitedBy = sel.Find(".profile-section__invited a").Text()
	if user.InvitedBy == "" {
		zap.L().Warn("Can't get user inviter")
	}

	return user, nil
}

func parseUserPosts(r io.Reader) ([]string, error) {
	postsIDs := []string{}
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return postsIDs, err
	}
	sel := doc.Find("*")

	sel.Find(".content-list__item_post").Each(func(i int, s *goquery.Selection) {
		id, ok := s.Attr("id")
		if !ok {
			zap.L().Warn("Can't get post id from user posts page")
			return
		}
		id = strings.TrimPrefix(id, "post_")
		postsIDs = append(postsIDs, id)
	})

	return postsIDs, nil
}
