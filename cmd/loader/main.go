package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	ponpon "github.com/PonPonLoader"
	"github.com/PonPonLoader/model"
	"github.com/codegangsta/cli"
)

var (
	cliFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "watch",
			Usage: "watch for new images from thread",
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Flags = cliFlags

	app.Action = func(c *cli.Context) error {
		threadURL := c.Args().Get(0)
		targetDir := c.Args().Get(1)
		watch := c.Bool("watch")

		thread, err := parseThreadURL(threadURL)
		if err != nil {
			return err
		}

		if err := createBaseDir(targetDir); err != nil {
			return err
		}

		posts := make(chan *model.Post)

		// Handle SIGINT and SIGTERM.
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		threadWatcher := ponpon.NewThreadWatcher(posts, thread)
		go func() {
			<-ch
			threadWatcher.Stop()
		}()

		go func() {
			if err := threadWatcher.Run(!watch); err != nil {
				panic(err)
			}

		}()

		imagePosts := mapPosts(posts, func(p *model.Post) *model.Post {
			if !p.HasImage() {
				return nil
			}

			p.BoardName = thread.BoardName
			return p
		})

		downloadTasks := mapPostsToImageDownloadTasks(imagePosts, targetDir)

		processor, err := ponpon.NewTaskProcessor(downloadTasks)
		if err != nil {
			return err
		}

		processor.Run()

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
