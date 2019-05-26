package crawler

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/models"
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

func (w *Worker) Work(maxDeep int) {
	for true {
		// get task from queue
		tasksSlice, err := w.queue.Pop(1)
		if err != nil {
			log.Printf("Can't pop tasks from queue: %v", err)
			continue
		}
		task := tasksSlice[0]
		if task.Deep > maxDeep {
			return
		}

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

		// process task
		newTasks, err := w.process(task, resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		tasksSl := []tasks.Task{}
		for t := range newTasks {
			exists, err := w.state.Exists(t)
			if err != nil {
				log.Println("Can't check task exists:", err)
				continue
			}
			if !exists {
				err = w.state.Add(t)
				if err != nil {
					log.Println("Can't add task to state store:", err)
					continue
				}
				tasksSl = append(tasksSl, t)
			}
		}

		_, err = w.queue.Push(tasksSl)
		if err != nil {
			log.Println("Can't push tasks to queue:", err)
			continue
		}
	}
}

func (w *Worker) process(task tasks.Task, r io.Reader) (map[tasks.Task]struct{}, error) {
	switch task.Type {
	case tasks.PostTask:
		return w.processPost(task, r)
	case tasks.UserTask:
		return w.processUser(task, r)
	case tasks.UserPostsTask:
		return w.processUserPosts(task, r)
	}
	return nil, errors.New("Undefined task type")
}

// Extract tasks for: post author, comments authors;
// Save content: post, comments
func (w *Worker) processPost(task tasks.Task, r io.Reader) (map[tasks.Task]struct{}, error) {
	newTasks := map[tasks.Task]struct{}{}

	post, comments, err := parsePost(r)
	post.ID = task.Body
	if err != nil {
		return newTasks, err
	}

	newTasks[tasks.Task{
		Type: tasks.UserTask,
		Body: post.Author,
		Deep: task.Deep + 1,
	}] = struct{}{}
	newTasks[tasks.Task{
		Type: tasks.UserPostsTask,
		Body: post.Author,
		Deep: task.Deep + 1,
		Page: 1,
	}] = struct{}{}
	for _, com := range comments {
		newTasks[tasks.Task{
			Type: tasks.UserTask,
			Body: com.Author,
			Deep: task.Deep + 1,
		}] = struct{}{}
		newTasks[tasks.Task{
			Type: tasks.UserPostsTask,
			Body: com.Author,
			Deep: task.Deep + 1,
			Page: 1,
		}] = struct{}{}
	}

	err = w.cnt.AddPosts([]models.Post{post})
	if err != nil {
		log.Println(err)
	}
	err = w.cnt.AddComments(comments)
	if err != nil {
		log.Println(err)
	}

	return newTasks, nil
}

func (w *Worker) processUser(task tasks.Task, r io.Reader) (map[tasks.Task]struct{}, error) {
	return nil, nil
}

func (w *Worker) processUserPosts(task tasks.Task, r io.Reader) (map[tasks.Task]struct{}, error) {

	return nil, nil
}
