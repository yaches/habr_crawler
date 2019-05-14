package crawler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yaches/habr_crawler/tasks"
)

const (
	userPageTemplate      = "https://habr.com/ru/users/%s/"
	userPostsPageTemplate = "https://habr.com/ru/users/%s/posts/page%d/"
	postPageTemplate      = "https://habr.com/ru/post/%s/"
)

func URLFromTask(task tasks.Task) (string, error) {
	switch task.Type {
	case tasks.PostTask:
		return fmt.Sprintf(postPageTemplate, task.Body), nil
	case tasks.UserTask:
		return fmt.Sprintf(userPageTemplate, task.Body), nil
	case tasks.UserPostsTask:
		return fmt.Sprintf(userPostsPageTemplate, task.Body, task.Page), nil
	}

	return "", errors.New("Undefined task type")
}

func TaskFromURL(url string) (tasks.Task, error) {
	userPrefix := "https://habr.com/ru/users/"
	postPrefix := "https://habr.com/ru/post/"

	if strings.Contains(url, userPrefix) {
		body := strings.TrimPrefix(url, userPrefix)
		body = strings.Trim(body, "/ ")
		return tasks.Task{
			Type: tasks.UserTask,
			Body: body,
		}, nil
	}

	if strings.Contains(url, postPrefix) {
		body := strings.TrimPrefix(url, postPrefix)
		body = strings.Trim(body, "/ ")
		return tasks.Task{
			Type: tasks.PostTask,
			Body: body,
		}, nil
	}

	return tasks.Task{}, fmt.Errorf("Can't extract task from url: %v", url)
}
