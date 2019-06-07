package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch"
	"github.com/yaches/habr_crawler/models"
)

const (
	userIndex    = "users"
	postIndex    = "posts"
	commentIndex = "comments"
)

var granMap = map[string]string{
	"hour": "getHour()",
	"day":  "getDayOfWeek().value",
	"year": "getYear()",
}

const searchQuery = `
	"query": {
		"term": {
			"%s": {
				"value": "%s"
			}
		}	
	}`

const filterQuery = `,{
	"term": {
		"%s": {
			"value": "%s"
		}
	}
}`

const histQuery = `{
	"size": 0, 
	"query": {
		"bool": {
			"must": [
				{
					"range": {
						"%[1]s": {
							  "gte": 1
						}
					}
				}
				%[3]s
			]
		}
	}, 
	"aggs": {
		"agg": {
			"terms": {
				"size": 100, 
				"script": {
					"lang": "painless",
					"source": "doc['%[1]s'].value.withZoneSameInstant(java.time.ZoneId.of('+03:00')).%[2]s"
				}
			}
		}
	}
}`

const hubHistQuery = `{
	"size": 0, 
	%s
	"aggs": {
		"agg": {
			"terms": {
				"size": 100, 
				"field": "%s.keyword"
			}
		}
	}
}`

type StorageES struct {
	es *elasticsearch.Client
}

func NewStorageES() (*StorageES, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	return &StorageES{es: es}, nil
}

func (s *StorageES) AddUsers(users []models.User) error {
	for _, u := range users {
		if err := s.insertOne(userIndex, u.Username, u); err != nil {
			return err
		}
	}
	return nil
}

func (s *StorageES) AddPosts(posts []models.Post) error {
	for _, p := range posts {
		if p.Text == "" {
			return nil
		}
		if err := s.insertOne(postIndex, p.ID, p); err != nil {
			return err
		}
	}
	return nil
}

func (s *StorageES) AddComments(comments []models.Comment) error {
	for _, c := range comments {
		if c.Text == "" {
			return nil
		}
		if err := s.insertOne(commentIndex, c.ID, c); err != nil {
			return err
		}
	}
	return nil
}

func (s *StorageES) insertOne(index, id string, obj interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(obj); err != nil {
		return err
	}
	res, err := s.es.Create(index, id, &buf)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageES) GetCommonCount(index string) (int, error) {
	mp := map[string]interface{}{}
	cResp, err := s.es.Count(s.es.Count.WithIndex(index))
	if err != nil {
		return 0, err
	}
	defer cResp.Body.Close()
	err = json.NewDecoder(cResp.Body).Decode(&mp)
	if err != nil {
		return 0, err
	}

	return int(mp["count"].(float64)), nil
}

func (s *StorageES) getHist(index, gran, filter string) (map[int]int, error) {
	timeField := "PubDate"
	if index == "users" {
		timeField = "RegDate"
	}
	q := []byte(fmt.Sprintf(histQuery, timeField, granMap[gran], filter))
	resp, err := s.es.Search(s.es.Search.WithIndex(index), s.es.Search.WithBody(bytes.NewReader(q)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&mp)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(q))
	// fmt.Println()
	// fmt.Println(mp)

	return parseTimeHist("agg", mp)
}

func (s *StorageES) get(index, id string) ([]byte, error) {
	res, err := s.es.Get(index, id, s.es.Get.WithDocumentType("_doc"))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var mp map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&mp)
	if err != nil {
		return nil, err
	}

	str, err := json.Marshal(mp["_source"])
	if err != nil {
		return nil, err
	}

	return str, nil
}

func (s *StorageES) GetUser(name string) (models.User, error) {
	var u models.User

	str, err := s.get("users", name)
	err = json.Unmarshal(str, &u)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (s StorageES) GetPost(id string) (models.Post, error) {
	var p models.Post

	str, err := s.get("posts", id)
	err = json.Unmarshal(str, &p)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s *StorageES) GetCommonHist(index, gran string) (map[int]int, error) {
	return s.getHist(index, gran, "")
}

func (s *StorageES) GetTermFilteredHist(index, gran, field, value string) (map[int]int, error) {
	f := fmt.Sprintf(filterQuery, field, value)
	return s.getHist(index, gran, f)
}

func parseTimeHist(agg string, resp map[string]interface{}) (map[int]int, error) {
	hist := map[int]int{}
	buckets := resp["aggregations"].(map[string]interface{})[agg].(map[string]interface{})["buckets"].([]interface{})
	for _, b := range buckets {
		key, err := strconv.ParseInt(b.(map[string]interface{})["key"].(string), 10, 32)
		val := b.(map[string]interface{})["doc_count"].(float64)
		if err != nil {
			return hist, err
		}
		hist[int(key)] = int(val)
	}

	return hist, nil
}

func parseHubHist(agg string, resp map[string]interface{}) (map[string]int, error) {
	hist := map[string]int{}
	buckets := resp["aggregations"].(map[string]interface{})[agg].(map[string]interface{})["buckets"].([]interface{})
	for _, b := range buckets {
		key := b.(map[string]interface{})["key"].(string)
		val := b.(map[string]interface{})["doc_count"].(float64)
		hist[key] = int(val)
	}

	return hist, nil
}

func (s *StorageES) GetHubCommonHist(field string) (map[string]int, error) {
	return s.getHubHist("", field)
}

func (s *StorageES) GetHubUserHist(user, field string) (map[string]int, error) {
	f := fmt.Sprintf(searchQuery, "Author", user)
	return s.getHubHist(f, field)
}

func (s *StorageES) getHubHist(filter, field string) (map[string]int, error) {
	if filter != "" {
		filter = filter + ","
	}
	q := []byte(fmt.Sprintf(hubHistQuery, filter, field))
	resp, err := s.es.Search(s.es.Search.WithIndex("posts"), s.es.Search.WithBody(bytes.NewReader(q)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&mp)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(q))
	// fmt.Println()
	// fmt.Println(mp)

	return parseHubHist("agg", mp)
}
