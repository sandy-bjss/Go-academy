package tasks

var requests chan request
var done chan struct{} = make(chan struct{})

const (
	getRequest    = "get"
	createRequest = "create"
	updateRequest = "update"
	deleteRequest = "delete"
	shutdown      = "shutdown"
)

type request struct {
	reqType  string
	response chan []Task
	task     Task
}

func Start(taskJSONfilestore string) {
	Actor(taskJSONfilestore)
}

func Stop() {
	shutdown := request{reqType: shutdown, response: nil, task: Task{}}
	requests <- shutdown
	<-done
}

func Actor(taskJSONfilestore string) {
	requests = make(chan request, 100)
	go func() {
		// actor stuff
		for req := range requests {
			switch req.reqType {
			case getRequest:
				// use load tasks from json file to memory
				// send back in response
				req.response <- GetTasks(taskJSONfilestore)

			case createRequest:
				// load tasks from json file
				// add the task to tasks list
				// save to updated list to json file
				// send back updated task list
				req.response <- CreateTask(req.task, taskJSONfilestore)

			case updateRequest:
				//load tasks from json file to memory
				// find supplied task & update
				// save updated task list
				// send back updated task list
				req.response <- UpdateTask(req.task, taskJSONfilestore)

			case deleteRequest:
				// load tasks from json file to memory
				// find supplied task and delete
				// save updated task list to json file
				// send back updated task list
				req.response <- DeleteTask(req.task.Id, taskJSONfilestore)

			case shutdown:
				close(req.response)
			}
		}
		close(done)
	}()
}

func Get() []Task {
	response := make(chan []Task, 1)
	req := request{
		reqType:  getRequest,
		response: response,
		task:     Task{},
	}
	requests <- req
	return <-response
}

func Create(task Task) []Task {
	response := make(chan []Task, 1)
	req := request{
		reqType:  createRequest,
		response: response,
		task:     task,
	}
	requests <- req
	return <-response
}

func Update(task Task) []Task {
	response := make(chan []Task, 1)
	req := request{
		reqType:  updateRequest,
		response: response,
		task:     task,
	}
	requests <- req
	return <-response
}

func Delete(task Task) []Task {
	response := make(chan []Task, 1)
	req := request{
		reqType:  deleteRequest,
		response: response,
		task:     task,
	}
	requests <- req
	return <-response
}
