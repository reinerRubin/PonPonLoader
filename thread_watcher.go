package ponpon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PonPonLoader/definition"
	"github.com/PonPonLoader/model"
)

// ThreadWatcher TBD
type ThreadWatcher struct {
	thread   *model.Thread
	outPosts chan<- *model.Post
	done     chan struct{}

	ticker <-chan time.Time
}

// NewThreadWatcher TBD
func NewThreadWatcher(outPosts chan<- *model.Post, thread *model.Thread) *ThreadWatcher {
	return &ThreadWatcher{
		outPosts: outPosts,
		thread:   thread,
		done:     make(chan struct{}),
	}
}

// Run TBD
func (tw *ThreadWatcher) Run(once bool) error {
	ticker := time.After(0 * time.Second)

	defer func() {
		close(tw.outPosts)
		tw.done <- struct{}{}
	}()

	for firstIteration := true; ; {
		select {
		case <-ticker:
			jsonThread, err := fetchThread(tw.thread)
			if err != nil {
				return err
			}

			// 404 error
			if jsonThread == nil {
				return nil
			}

			for post := range genPostsFromThread(jsonThread) {
				tw.outPosts <- post
			}
		case <-tw.done:
			return nil
		}

		if !once && firstIteration {
			firstIteration = false
			// memory leak
			ticker = time.Tick(definition.UpdateEverySeconds * time.Second)
		}

		if once {
			return nil
		}
	}
}

// Stop TBD
func (tw *ThreadWatcher) Stop() {
	tw.done <- struct{}{}
	<-tw.done
}

func fetchThread(thread *model.Thread) (*model.JSONThread, error) {
	URLString := fmt.Sprintf(
		"%s/%s/thread/%d.json",
		definition.APIHost, thread.BoardName, thread.No,
	)

	// TODO: add If-Modified-Since header
	resp, err := http.Get(URLString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}

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
