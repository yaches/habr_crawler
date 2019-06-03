package crawler

import (
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/yaches/habr_crawler/content"
	"github.com/yaches/habr_crawler/models"
	"github.com/yaches/habr_crawler/tasks"
	"go.uber.org/zap"
)

type Worker struct {
	cnt     content.Storage
	queue   tasks.Manager
	counter counter
	done    chan struct{}
	mx      *sync.Mutex
	httpCl  *http.Client
}

func NewWorker(cnt content.Storage, queue tasks.Manager, client *http.Client) *Worker {
	return &Worker{
		cnt:    cnt,
		queue:  queue,
		httpCl: client,
	}
}

func (w *Worker) Start(threads, maxDeep int) {
	w.counter = NewCounter(threads)
	w.done = make(chan struct{}, threads)
	w.mx = &sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(threads)

	for i := 0; i < threads; i++ {
		go w.Work(maxDeep, &wg, i)
	}

	wg.Wait()
}

func (w *Worker) Work(maxDeep int, wg *sync.WaitGroup, n int) {
	zap.L().Info("[WORKER STARTED]", zap.Int("n", n))
	defer wg.Done()
	defer zap.L().Info("[WORKER STOPPED]", zap.Int("n", n))
	// Handle double-closing done channel
	// defer func() {
	// 	recover()
	// }()

	for {
		var task tasks.Task

		w.counter.Dec()

		// if w.counter.Zero() && len(w.queue.Channel()) == 0 {
		// 	close(w.done)
		// }
		if w.counter.Zero() && len(w.queue.Channel()) == 0 {
			zap.L().Info("Closing signal", zap.Int("n", n))
			w.mx.Lock()
			select {
			// Смогли прочитать, значит канал уже закрыт
			case <-w.done:
			// Не смогли прочитать, ушли в default, значит канал не закрыт, и надо закрыть
			default:
				close(w.done)
			}
			w.mx.Unlock()
		}

		select {
		case task = <-w.queue.Channel():
			w.counter.Inc()
		case <-w.done:
			return
		}

		zap.L().Debug("[GET task from queue]", zap.Int("n", n))

		if task.Deep > maxDeep {
			zap.L().Debug("[DROP task] max deep is reached", zap.Int("n", n))
			continue
		}

		// convert task to url
		url, err := URLFromTask(task)
		if err != nil {
			zap.L().Warn("URL creating from task error", zap.Error(err))
			continue
		}

		// download content
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			zap.L().Warn("Can't create http request", zap.Error(err))
		}
		resp, err := w.httpCl.Do(req)
		if err != nil {
			zap.L().Warn("HTTP GET error", zap.Error(err))
			continue
		}

		// process task
		newTasks, err := w.process(task, resp.Body, n)
		if err != nil {
			zap.L().Warn("Task processing error", zap.Error(err))
			continue
		}

		resp.Body.Close()

		err = w.queue.Done(task)
		if err != nil {
			zap.L().Warn("Can't mark task as done")
		}
		for t := range newTasks {
			err = w.queue.Push(t)
			if err != nil {
				zap.L().Warn("Can't push task to queue", zap.Error(err))
				continue
			}
		}
	}
}

func (w *Worker) process(task tasks.Task, r io.Reader, n int) (map[tasks.Task]struct{}, error) {
	zap.L().Info("[TASK PROCESSING]", zap.Int("n", n), zap.Any("task", task))

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
