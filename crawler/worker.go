package crawler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/yaches/habr_crawler/models"

	"github.com/PuerkitoBio/goquery"

	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/state"
	"github.com/yaches/habr_crawler/tasks"
)

type Worker struct {
	cnt   content.Storage
	state state.Storage
	queue tasks.TaskManager
}

func NewWorker(cnt content.Storage, state state.Storage, queue tasks.TaskManager) *Worker {
	return &Worker{
		cnt:   cnt,
		state: state,
		queue: queue,
	}
}

func (w *Worker) Work() {
	for true {
		// get task from queue
		tasksSlice, err := w.queue.Pop(1)
		if err != nil {
			log.Printf("Can't pop tasks from queue")
		}
		task := tasksSlice[0]

		// check task already complete
		// if task.Type != tasks.UserPostsTask {
		// 	ok, err := state.Exists(task.Body)
		// 	if err != nil {
		// 		log.Print(err.Error)
		// 		continue
		// 	}
		// 	if ok {
		// 		continue
		// 	}
		// }

		// convert task to url
		url, err := URLFromTask(task)
		if err != nil {
			log.Print(err.Error())
			continue
		}

		// download content
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// process task
		newTasks, err := w.process(task, doc.Find("*"))
		if err != nil {
			log.Println(err)
		}

		log.Println(newTasks)
		return

		// parse html

		// save to content storage

		// extract new tasks

		// check tasks exists
		// push tasks to queue
	}
}

func (w *Worker) process(task tasks.Task, sel *goquery.Selection) ([]tasks.Task, error) {
	switch task.Type {
	case tasks.PostTask:
		return w.processPost(task, sel)
	case tasks.UserTask:
		return w.processUser(task, sel)
	case tasks.UserPostsTask:
		return w.processUserPosts(task, sel)
	}
	return nil, errors.New("Undefined task type")
}

// Extract tasks for: post author, comments authors;
// Save content: post, comments
func (w *Worker) processPost(task tasks.Task, sel *goquery.Selection) ([]tasks.Task, error) {
	newTasks := []tasks.Task{}
	post := models.Post{ID: task.Body}

	// Get and parse user URL
	url, ok := sel.Find(".post__meta a").Attr("href")
	if ok {
		// All new tasks takes incrementing deep from PostTask
		authorName, authorTasks, err := tasksFromUserURL(task.Deep+1, url)
		if err != nil {
			log.Println(err)
		}
		newTasks = append(newTasks, authorTasks...)
		post.Author = authorName
	} else {
		log.Println("Can't get author url from post page")
	}

	// Get post pub time
	pubTimeStr, ok := sel.Find(".post__meta .post__time").Attr("data-time_published")
	if ok {
		pubTime, err := parseTime(pubTimeStr)
		if err != nil {
			log.Println(err)
		}
		post.PubDate = pubTime
	} else {
		log.Println("Can't get post pub time")
	}

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

	log.Println(post)

	return newTasks, nil
}

func tasksFromUserURL(deep int, url string) (string, []tasks.Task, error) {
	newTasks := []tasks.Task{}
	// Create 2 tasks: UserTask and UserPostsTask{page=1}
	authorTask, err := TaskFromURL(url)
	if err == nil {
		authorTask.Deep = deep
		authorPostsTask := tasks.Task{
			Type: tasks.UserPostsTask,
			Body: authorTask.Body,
			Deep: deep,
			Page: 1,
		}

		newTasks = append(newTasks, authorTask)
		newTasks = append(newTasks, authorPostsTask)
	}
	return authorTask.Body, newTasks, err
}

func (w *Worker) processUser(task tasks.Task, sel *goquery.Selection) ([]tasks.Task, error) {
	return nil, nil
}

func (w *Worker) processUserPosts(task tasks.Task, sel *goquery.Selection) ([]tasks.Task, error) {
	return nil, nil
}

func parseTime(t string) (time.Time, error) {
	if t == "" {
		return time.Time{}, errors.New("Empty time string")
	}
	return time.Parse("2006-01-02T15:04Z07:00", t)
}
