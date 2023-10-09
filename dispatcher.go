package hardworker

// Dispatcher holds a double-layer channel of Job.
type Dispatcher struct {
	// JobQueue is a channel that transports Job between inside and outside.
	JobQueue chan Job

	workerPool chan chan Job
	maxWorker  int
	exit       chan bool
}

// NewDispatcher makes a pointer to Dispatcher.
func NewDispatcher(maxWorker int) *Dispatcher {
	return &Dispatcher{
		JobQueue:   make(chan Job, maxWorker),
		workerPool: make(chan chan Job, maxWorker),
		maxWorker:  maxWorker,
		exit:       make(chan bool),
	}
}

// Run makes some worker(s).
// It also makes a goroutine that calls function dispatch.
func (d *Dispatcher) Run() []*worker {
	workers := make([]*worker, 0, d.maxWorker)

	for i := 0; i < d.maxWorker; i++ {
		worker := NewWorker(d.workerPool)
		worker.Start()

		workers = append(workers, worker)
	}

	go d.dispatch()

	return workers
}

// dispatch makes a infinite loop to receive Job from the channel JobQueue.
func (d *Dispatcher) dispatch() {
	for job := range d.JobQueue { // Receive a job from the channel JobQueue.
		// Try to get a worker from workerPool.
		// Block the current goroutine until it gets an available worker.
		worker := <-d.workerPool

		// Send job to this worker.
		worker <- job
	}

	d.exit <- true
}

// Exit passes true to the channel exit so makes a safe exit.
func (d *Dispatcher) Exit() {
	close(d.JobQueue)
	<-d.exit
}
