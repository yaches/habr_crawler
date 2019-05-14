package crawler

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/yaches/habr_crawler/contentstore"
	"github.com/yaches/habr_crawler/statestore"
	"github.com/yaches/habr_crawler/tasks"
)

func Work(cnt contentstore.Storage, state statestore.Storage, queue tasks.TaskManager) {
	for true {
		// get task from queue
		tasksSlice, err := queue.Pop(1)
		if err != nil {
			log.Printf("Can't pop tasks from queue")
		}
		task := tasksSlice[0]

		// check task already complete
		if task.Type != tasks.UserPostsTask {
			ok, err := state.Exists(task.Body)
			if err != nil {
				log.Print(err.Error)
				continue
			}
			if ok {
				continue
			}
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
			log.Print(err.Error())
			continue
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err.Error())
		}
		dataStr := string(data)

		log.Print(dataStr)
		return

		// parse html

		// save to content storage

		// extract new tasks

		// push tasks to queue
	}
}
