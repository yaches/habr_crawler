package content

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"strconv"

// 	"github.com/elastic/go-elasticsearch"
// 	"github.com/yaches/habr_crawler/models"
// 	"github.com/yaches/habr_crawler/models/aggs"
// )

// const (
// 	userIndex    = "users"
// 	postIndex    = "posts"
// 	commentIndex = "comments"
// )

// const histQuery = `{
// 	"size": 0,
// 	"query": {
// 	  "range": {
// 		"%[1]s": {
// 		  "gte": 1
// 		}
// 	  }
// 	},
// 	"aggs": {
// 	  "day": {
// 		"terms": {
// 		  "size": 100,
// 		  "script": {
// 			"lang": "painless",
// 			"source": "doc['%[1]s'].value.withZoneSameInstant(java.time.ZoneId.of('+03:00')).getDayOfWeek().value"
// 		  }
// 		}
// 	  },
// 	  "hour": {
// 		"terms": {
// 		  "size": 100,
// 		  "script": {
// 			"lang": "painless",
// 			"source": "doc['%[1]s'].value.withZoneSameInstant(java.time.ZoneId.of('+03:00')).getHour()"
// 		  }
// 		}
// 	  },
// 	  "month": {
// 		"terms": {
// 		  "size": 100,
// 		  "script": {
// 			"lang": "painless",
// 			"source": "doc['%[1]s'].value.withZoneSameInstant(java.time.ZoneId.of('+03:00')).getMonth().value"
// 		  }
// 		}
// 	  },
// 	  "year": {
// 		"terms": {
// 		  "size": 100,
// 		  "script": {
// 			"lang": "painless",
// 			"source": "doc['%[1]s'].value.withZoneSameInstant(java.time.ZoneId.of('+03:00')).getYear()"
// 		  }
// 		}
// 	  }
// 	}
//   }`

// type StorageES struct {
// 	es *elasticsearch.Client
// }

// func NewStorageES() (*StorageES, error) {
// 	es, err := elasticsearch.NewDefaultClient()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &StorageES{es: es}, nil
// }

// func (s *StorageES) AddUsers(users []models.User) error {
// 	for _, u := range users {
// 		if err := s.insertOne(userIndex, u.Username, u); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (s *StorageES) GetAllUsers() ([]models.User, error) {
// 	return nil, nil
// }

// func (s *StorageES) AddPosts(posts []models.Post) error {
// 	for _, p := range posts {
// 		if p.Text == "" {
// 			return nil
// 		}
// 		if err := s.insertOne(postIndex, p.ID, p); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (s *StorageES) GetAllPosts() ([]models.Post, error) {
// 	return nil, nil
// }

// func (s *StorageES) AddComments(comments []models.Comment) error {
// 	for _, c := range comments {
// 		if c.Text == "" {
// 			return nil
// 		}
// 		if err := s.insertOne(commentIndex, c.ID, c); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (s *StorageES) GetAllComments() ([]models.Comment, error) {
// 	return nil, nil
// }

// func (s *StorageES) insertOne(index, id string, obj interface{}) error {
// 	var buf bytes.Buffer
// 	if err := json.NewEncoder(&buf).Encode(obj); err != nil {
// 		return err
// 	}
// 	res, err := s.es.Create(index, id, &buf)
// 	defer res.Body.Close()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *StorageES) GetCommonInfo() (aggs.CommonInfo, error) {
// 	i := aggs.CommonInfo{}

// 	indexes := map[string]aggs.Info{}
// 	indexes["posts"] = aggs.Info{}
// 	indexes["comments"] = aggs.Info{}
// 	indexes["users"] = aggs.Info{}

// 	for index, _ := range indexes {
// 		timeField := "PubDate"
// 		if index == "users" {
// 			timeField = "RegDate"
// 		}
// 		q := []byte(fmt.Sprintf(histQuery, timeField))
// 		resp, err := s.es.Search(s.es.Search.WithIndex(index), s.es.Search.WithBody(bytes.NewReader(q)))
// 		if err != nil {
// 			return i, err
// 		}
// 		defer resp.Body.Close()

// 		mp := map[string]interface{}{}
// 		err = json.NewDecoder(resp.Body).Decode(&mp)
// 		if err != nil {
// 			return i, err
// 		}
// 		h, err := getTimeHist("hour", mp)
// 		d, err := getTimeHist("day", mp)
// 		m, err := getTimeHist("month", mp)
// 		y, err := getTimeHist("year", mp)
// 		if err != nil {
// 			return i, err
// 		}

// 		cResp, err := s.es.Count(s.es.Count.WithIndex(index))
// 		if err != nil {
// 			return i, err
// 		}
// 		defer cResp.Body.Close()
// 		err = json.NewDecoder(cResp.Body).Decode(&mp)
// 		if err != nil {
// 			return i, err
// 		}

// 		indexes[index] = aggs.Info{
// 			Count:         int(mp["count"].(float64)),
// 			HourHist:      h,
// 			DayOfWeekHist: d,
// 			MonthHist:     m,
// 			YearHist:      y,
// 		}
// 	}

// 	i.PostsInfo = indexes["posts"]
// 	i.CommentsInfo = indexes["comments"]
// 	i.UsersInfo = indexes["users"]

// 	return i, nil
// }

// // func GetCommonTimeHist(index, gran string) (map[int]int, error) {
// // 	timeField := "PubDate"
// // 	if index == "users" {
// // 		timeField = "RegDate"
// // 	}

// // }

// func getTimeHist(agg string, resp map[string]interface{}) (map[int]int, error) {
// 	hist := map[int]int{}
// 	buckets := resp["aggregations"].(map[string]interface{})[agg].(map[string]interface{})["buckets"].([]interface{})
// 	for _, b := range buckets {
// 		key, err := strconv.ParseInt(b.(map[string]interface{})["key"].(string), 10, 32)
// 		val := b.(map[string]interface{})["doc_count"].(float64)
// 		if err != nil {
// 			return hist, err
// 		}
// 		hist[int(key)] = int(val)
// 	}

// 	return hist, nil
// }
