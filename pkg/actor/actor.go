package actor

type RequestMessage interface{}

type RequestHandlerFunc func(re RequestMessage)

type Actor struct {
	Requests      chan RequestMessage
	done          chan struct{}
	handleRequest RequestHandlerFunc
}

// spawn new actor
func New(reqHandler RequestHandlerFunc) *Actor {
	return &Actor{
		Requests:      make(chan RequestMessage),
		done:          make(chan struct{}),
		handleRequest: reqHandler,
	}
}

func (a *Actor) Start() {
	go func() {
		for {
			select {
			case req := <-a.Requests:
				a.handleRequest(req)
			case <-a.done:
				return
			}
		}
	}()
}

func (a *Actor) Stop() {
	close(a.done)
}
