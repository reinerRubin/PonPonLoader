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
	thread            *model.Thread
	outPosts          chan<- *model.Post
	done, doneConfirm chan struct{}

	ticker <-chan time.Time

	lastModeifedHeader string
}

// NewThreadWatcher TBD
func NewThreadWatcher(outPosts chan<- *model.Post, thread *model.Thread) *ThreadWatcher {
	return &ThreadWatcher{
		outPosts: outPosts,
		thread:   thread,

		done:        make(chan struct{}),
		doneConfirm: make(chan struct{}),
	}
}

// Run TBD
func (tw *ThreadWatcher) Run(once bool) error {
	tw.ticker = time.After(0)

	defer func() {
		close(tw.outPosts)
		tw.doneConfirm <- struct{}{}
	}()

	for firstIteration := true; ; {
		select {
		case <-tw.ticker:
			jsonThread, err := tw.fetchThread()
			if err != nil {
				return err
			}

			if jsonThread != nil {
				for _, post := range jsonThread.ToPosts() {
					tw.outPosts <- post
				}
			}
		case <-tw.done:
			return nil
		}

		if !once && firstIteration {
			firstIteration = false
			// memory leak
			tw.ticker = time.Tick(definition.UpdateEverySeconds * time.Second)
		}

		if once {
			return nil
		}
	}
}

func (tw *ThreadWatcher) asyncStop() {
	tw.ticker = nil
	tw.done <- struct{}{}
}

// Stop TBD
func (tw *ThreadWatcher) Stop() {
	tw.asyncStop()
	<-tw.doneConfirm
}

func (tw *ThreadWatcher) fetchThread() (*model.JSONThread, error) {
	resp, err := tw.makeReq()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tw.lastModeifedHeader = resp.Header.Get("Last-Modified")

	if resp.StatusCode == 404 {
		tw.asyncStop()
		return nil, nil
	}

	if resp.StatusCode == 304 {
		return nil, nil
	}

	return tw.parseResponse(resp)
}

func (tw *ThreadWatcher) makeReq() (*http.Response, error) {
	URLString := fmt.Sprintf("%s/%s", definition.APIHost, tw.thread.URLPath())

	client := http.DefaultClient
	req, err := http.NewRequest("GET", URLString, nil)
	if err != nil {
		return nil, err
	}

	if tw.lastModeifedHeader != "" {
		req.Header.Add("If-Modified-Since", tw.lastModeifedHeader)
	}

	return client.Do(req)
}

func (tw *ThreadWatcher) parseResponse(resp *http.Response) (*model.JSONThread, error) {
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
