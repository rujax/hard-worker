package hardworker

// Job is an interface that declares function Work.
// All jobs must implement this function.
type Job interface {
	Work()
}

// worker has a channel jobChannel that receives Job from JobQueue and sends it.
type worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
}

// NewWorker makes a pointer to worker with the same workerPool from Dispatcher.
func NewWorker(workerPool chan chan Job) *worker {
	return &worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

// Start makes a infinite loop that receives Job from jobChannel and calls Job.Work.
func (w *worker) Start() {
	go func() {
		for {
			// Register jobChannel to workerPool.
			w.workerPool <- w.jobChannel

			// select blocks current goroutine until channels which receives a value.
			select {
			case job := <-w.jobChannel:
				job.Work()

			case <-w.quit:
				return
			}
		}
	}()
}

// Stop passes true to the channel quit so makes a safe exit.
func (w *worker) Stop() {
	w.quit <- true
}
