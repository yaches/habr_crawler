package crawler

import (
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/models"
	"github.com/yaches/habr_crawler/state"
	"github.com/yaches/habr_crawler/tasks"
	"go.uber.org/zap"
)

type Worker struct {
	cnt     content.Storage
	state   state.Storage
	queue   tasks.Manager
	counter counter
	done    chan struct{}
	mx      *sync.Mutex
}

func NewWorker(cnt content.Storage, state state.Storage, queue tasks.Manager) *Worker {
	return &Worker{
		cnt:   cnt,
		state: state,
		queue: queue,
	}
}

func (w *Worker) Start(threads, maxDeep int) {
	w.counter = NewCounter(threads)
	w.done = make(chan struct{}, threads)
	w.mx = &sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(threads)

	// go func(wg *sync.WaitGroup) {
	// 	defer wg.Done()
	// 	for {
	// 		select {
	// 		case <-w.counter.Zero():
	// 			if len(w.queue.Channel()) == 0 {
	// 				close(w.done)
	// 				return
	// 			}
	// 		}
	// 	}
	// }(&wg)

	for i := 0; i < threads; i++ {
		go w.Work(maxDeep, &wg)
	}

	zap.L().Warn("11")

	wg.Wait()

	zap.L().Warn("22")

}

func (w *Worker) Work(maxDeep int, wg *sync.WaitGroup) {
	zap.L().Info("[WORKER STARTED]")
	defer wg.Done()
	defer zap.L().Info("[WORKER STOPPED]")
	defer func() {
		recover()
	}()

	for {
		// get task from queue

		var task tasks.Task

		w.counter.Dec()

		if w.counter.Zero() && len(w.queue.Channel()) == 0 {
			close(w.done)
		}

		select {
		case task = <-w.queue.Channel():
			w.counter.Inc()
		case <-w.done:
			return
		}

		zap.L().Debug("[GET task from queue]")

		if task.Deep > maxDeep {
			zap.L().Debug("[DROP task] max deep is reached")
			continue
		}

		// convert task to url
		url, err := URLFromTask(task)
		if err != nil {
			zap.L().Warn("URL creating from task error", zap.Error(err))
			continue
		}

		// download content
		resp, err := http.Get(url)
		if err != nil {
			zap.L().Warn("HTTP GET error", zap.Error(err))
			continue
		}

		// process task
		newTasks, err := w.process(task, resp.Body)
		if err != nil {
			zap.L().Warn("Task processing error", zap.Error(err))
			continue
		}

		resp.Body.Close()

		for t := range newTasks {
			exists, err := w.state.Exists(t)
			if err != nil {
				zap.L().Warn("Check task existing error", zap.Error(err))
				continue
			}
			if !exists {
				err = w.state.Add(t)
				if err != nil {
					zap.L().Warn("Can't add task to state store", zap.Error(err))
					continue
				}
				err = w.queue.Push(t)
				if err != nil {
					zap.L().Warn("Can't push task to queue", zap.Error(err))
					continue
				}
			}
		}
	}
}

func (w *Worker) process(task tasks.Task, r io.Reader) (map[tasks.Task]struct{}, error) {
	zap.L().Info("[TASK PROCESSING]", zap.Any("task", task))

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

	if post.Author != "" {
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
	}
	for _, com := range comments {
		if com.Author != "" {
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
	}

	err = w.cnt.AddPosts([]models.Post{post})
	if err != nil {
		return newTasks, err
	}
	err = w.cnt.AddComments(comments)
	if err != nil {
		return newTasks, err
	}

	return newTasks, nil
}

func (w *Worker) processUser(task tasks.Task, r io.Reader) (map[tasks.Task]struct{}, error) {
	newTasks := map[tasks.Task]struct{}{}
	user, err := parseUser(r)
	user.Username = task.Body
	if err != nil {
		return newTasks, err
	}

	err = w.cnt.AddUsers([]models.User{user})
	if err != nil {
		return newTasks, err
	}

	return newTasks, nil
}

func (w *Worker) processUserPosts(task tasks.Task, r io.Reader) (map[tasks.Task]struct{}, error) {
	newTasks := map[tasks.Task]struct{}{}
	postsIDs, err := parseUserPosts(r)
	if err != nil {
		return newTasks, err
	}

	// Add posts tasks, increment deep
	for _, id := range postsIDs {
		if id != "" {
			newTasks[tasks.Task{
				Type: tasks.PostTask,
				Body: id,
				Deep: task.Deep + 1,
			}] = struct{}{}
		}
	}

	// Add next posts page, not increment deep
	if len(postsIDs) > 0 {
		nextTask := task
		nextTask.Page = task.Page + 1
		newTasks[nextTask] = struct{}{}
	}

	return newTasks, nil
}
