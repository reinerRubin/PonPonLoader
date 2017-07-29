package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PonPonLoader/definition"
	"github.com/PonPonLoader/model"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		threadURL := c.Args().Get(0)
		targetDir := c.Args().Get(1)

		thread, err := parseThreadURL(threadURL)
		if err != nil {
			panic(err)
		}

		if err := createBaseDir(targetDir); err != nil {
			panic(err)
		}

		jsonThread, err := fetchThread(thread)
		if err != nil {
			panic(err)
		}

		posts := genPostsFromThread(jsonThread)
		imagePosts := mapPosts(posts, func(p *model.Post) *model.Post {
			if !p.HasImage() {
				return nil
			}

			p.BoardName = thread.BoardName
			return p
		})

		downloadTasks := mapPostsToImageDownloadTasks(imagePosts, targetDir)
		executeDownloadTasks(downloadTasks)

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func createBaseDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func parseThreadURL(threadURLString string) (*model.Thread, error) {
	threadURL, err := url.Parse(threadURLString)
	if err != nil {
		return nil, err
	}

	paths := strings.Split(threadURL.Path, "/")
	// /c/thread/2942063/madotsuki-thread'
	if len(paths) < 3 {
		return nil, fmt.Errorf("cant parse URL: %s", threadURLString)
	}
	boardName := paths[1]
	boardNoStr := paths[3]

	boardNo, err := strconv.ParseInt(boardNoStr, 10, 0)
	if err != nil {
		return nil, err
	}

	return &model.Thread{
		No:        boardNo,
		BoardName: boardName,
	}, nil

}

func executeDownloadTasks(tasks <-chan *model.DownloadTask) {
	const numDigesters = 3

	var wg sync.WaitGroup
	wg.Add(numDigesters)

	for i := 0; i < numDigesters; i++ {
		go func() {
			for task := range tasks {
				file, err := os.Create(task.Target)
				if err != nil {
					panic(err)
				}
				defer file.Close()

				resp, err := http.Get(task.Source.String())
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				if _, err := io.Copy(file, resp.Body); err != nil {
					panic(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func fetchThread(thread *model.Thread) (*model.JSONThread, error) {
	URLString := fmt.Sprintf(
		"%s/%s/thread/%d.json",
		definition.APIHost, thread.BoardName, thread.No,
	)

	resp, err := http.Get(URLString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jsThread := &model.JSONThread{}
	if err := json.Unmarshal(body, jsThread); err != nil {
		return nil, err
	}

	return jsThread, nil
}

func mapPosts(posts <-chan *model.Post, f func(p *model.Post) *model.Post) <-chan *model.Post {
	c := make(chan *model.Post)

	go func() {
		defer close(c)

		for post := range posts {
			if newPost := f(post); newPost != nil {
				c <- post
			}
		}
	}()

	return c
}

func mapPostsToImageDownloadTasks(posts <-chan *model.Post, basePath string) <-chan *model.DownloadTask {
	c := make(chan *model.DownloadTask)
	go func() {
		defer close(c)

		for post := range posts {
			if task, err := model.NewDownloadTask(post, basePath); err == nil {
				c <- task
				continue
			}
			// TODO log this situation
		}
	}()
	return c
}

func genPostsFromThread(jsThread *model.JSONThread) <-chan *model.Post {
	c := make(chan *model.Post)
	go func() {
		defer close(c)

		for _, post := range jsThread.ToPosts() {
			c <- post
		}
	}()
	return c
}
