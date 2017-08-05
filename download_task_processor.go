package ponpon

import (
	"fmt"

	"github.com/PonPonLoader/model"
)

const (
	maxAttempts = 3
	capacity    = 10000
	workersNum  = 3
)

// TaskProcessor TBD
type TaskProcessor struct {
	incomeTasks <-chan *model.DownloadTask

	plannedTasks chan *taskToProcess
	failedTasks  chan *taskToProcess
	readyTasks   chan *taskToProcess

	taskRegistor map[string]struct{}
	taskInWork   int
}

// NewTaskProcessor TBD
func NewTaskProcessor(incomeTasks <-chan *model.DownloadTask) (*TaskProcessor, error) {
	taskProcessor := &TaskProcessor{
		incomeTasks:  incomeTasks,
		taskRegistor: make(map[string]struct{}),

		plannedTasks: make(chan *taskToProcess, capacity),
		failedTasks:  make(chan *taskToProcess, capacity),
		readyTasks:   make(chan *taskToProcess, capacity),

		taskInWork: 0,
	}
	return taskProcessor, nil
}

type taskToProcess struct {
	*model.DownloadTask
	attempt int
}

func newTaskToProcess(task *model.DownloadTask) *taskToProcess {
	return &taskToProcess{
		DownloadTask: task,
		attempt:      0,
	}
}

func (tp *taskToProcess) incAttempts() {
	tp.attempt++
}

func (tp *taskToProcess) attempts() int {
	return tp.attempt
}

// Run TBD
func (tp *TaskProcessor) Run() {
	defer tp.closeAll()

	tp.runTaskConsumers()

	for incomeTasks, closed := tp.incomeTasks, false; ; {
		select {
		case task, moreIncomeTasks := <-incomeTasks:
			if !moreIncomeTasks {
				incomeTasks = nil
				closed = true
				break
			}

			tp.processTask(newTaskToProcess(task))
		case failedTask := <-tp.failedTasks:
			tp.processFailedTask(failedTask)
		case _ = <-tp.readyTasks:
			tp.taskInWork--
		}

		if closed && tp.noTask() {
			return
		}
	}
}

func (tp *TaskProcessor) noTask() bool {
	return tp.taskInWork == 0
}

func (tp *TaskProcessor) closeAll() {
	close(tp.plannedTasks)
	close(tp.failedTasks)
	close(tp.readyTasks)
}

func (tp *TaskProcessor) processFailedTask(task *taskToProcess) {
	if task.attempts() >= maxAttempts {
		fmt.Printf("task (%s) failed with %d attempts\n",
			task.Source.String(), task.attempts())

		delete(tp.taskRegistor, tp.taskKey(task))
		tp.taskInWork--

		return
	}

	task.incAttempts()
	tp.plannedTasks <- task
}

func (tp *TaskProcessor) processTask(task *taskToProcess) {
	taskKey := tp.taskKey(task)
	if _, found := tp.taskRegistor[taskKey]; found {
		return
	}

	tp.taskInWork++
	tp.taskRegistor[taskKey] = struct{}{}
	tp.plannedTasks <- task
}

func (tp *TaskProcessor) runTaskConsumers() {
	worker := func() {
		for task := range tp.plannedTasks {
			if err := task.Run(); err != nil {
				tp.failedTasks <- task
				continue
			}

			tp.readyTasks <- task
		}
	}

	for i := 0; i < workersNum; i++ {
		go worker()
	}
}

func (tp *TaskProcessor) taskKey(task *taskToProcess) string {
	return task.Source.String() + "~8~" + task.Target
}
