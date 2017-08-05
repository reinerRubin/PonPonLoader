package ponpon

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/PonPonLoader/model"
	"github.com/codegangsta/cli"
)

// App TBD
type App struct {
	thread    *model.Thread
	watch     bool
	targetDir string
}

// NewApp TBD
func NewApp(c *cli.Context) (*App, error) {
	threadURL := c.Args().Get(0)
	targetDir := c.Args().Get(1)
	watch := c.Bool("watch")

	thread, err := model.NewThreadFromURL(threadURL)
	if err != nil {
		return nil, err
	}

	return &App{
		thread:    thread,
		watch:     watch,
		targetDir: targetDir,
	}, nil
}

// Run TBD
func (a *App) Run() error {
	if err := a.prepareToRun(); err != nil {
		return err
	}

	posts := make(chan *model.Post)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	threadWatcher := NewThreadWatcher(posts, a.thread)
	go func() {
		<-ch
		threadWatcher.Stop()
	}()

	go func() {
		if err := threadWatcher.Run(!a.watch); err != nil {
			panic(err)
		}

	}()

	imagePosts := mapPosts(posts, func(p *model.Post) *model.Post {
		if !p.HasImage() {
			return nil
		}

		p.BoardName = a.thread.BoardName
		return p
	})

	downloadTasks := mapPostsToImageDownloadTasks(imagePosts, a.targetDir)

	processor, err := NewTaskProcessor(downloadTasks)
	if err != nil {
		return err
	}

	processor.Run()

	return nil
}

// PrepareToRun TBD
func (a *App) prepareToRun() error {
	if err := a.createBaseDir(); err != nil {
		return err
	}

	return nil
}

func (a *App) createBaseDir() error {
	return os.MkdirAll(a.targetDir, os.ModePerm)
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
