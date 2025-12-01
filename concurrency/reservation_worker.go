package concurrency

// Reservation worker that processes reservation requests concurrently using
// a pool of goroutines and a channel-based queue. The worker does not have
// any knowledge of the library internals — the caller provides a processor
// function which will be executed for each request.

import "sync"

// ReservationRequest is a request to reserve a book for a member.
type ReservationRequest struct {
	BookID   int
	MemberID int
	RespChan chan error
}

// ReservationCenter manages worker goroutines that process reservation requests.
type ReservationCenter struct {
	Requests    chan ReservationRequest
	wg          sync.WaitGroup
	workerCount int
	process     func(req ReservationRequest)
	stop        chan struct{}
}

// NewReservationCenter creates a new center with `workerCount` goroutines that
// call the provided process function for each request.
func NewReservationCenter(workerCount int, process func(req ReservationRequest)) *ReservationCenter {
	return &ReservationCenter{
		Requests:    make(chan ReservationRequest, 256),
		workerCount: workerCount,
		process:     process,
		stop:        make(chan struct{}),
	}
}

// Start launches the worker goroutines.
func (r *ReservationCenter) Start() {
	for i := 0; i < r.workerCount; i++ {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			for {
				select {
				case req := <-r.Requests:
					r.process(req)
				case <-r.stop:
					return
				}
			}
		}()
	}
}

// Stop signals workers to shut down and waits for them to finish.
func (r *ReservationCenter) Stop() {
	close(r.stop)
	r.wg.Wait()
}

// Enqueue adds a request to the queue for processing.
func (r *ReservationCenter) Enqueue(req ReservationRequest) {
	r.Requests <- req
}
