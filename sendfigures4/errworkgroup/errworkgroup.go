package errworkgroup

import "sync"

type ErrorWorkGroup struct {
	routinesLimiter chan struct{}
	wg              sync.WaitGroup
	errMutex        sync.RWMutex
	err             error
	skipWhenError   bool
}

func NewErrorWorkGroup(numWorkers int, skipWhenError bool) ErrorWorkGroup {
	if numWorkers < 1 {
		numWorkers = 1
	}
	return ErrorWorkGroup{
		routinesLimiter: make(chan struct{}, numWorkers),
		skipWhenError:   skipWhenError,
	}
}

func (w *ErrorWorkGroup) Wait() error {
	w.wg.Wait()

	return w.err
}

func (w *ErrorWorkGroup) Go(work func() error) {
	w.wg.Add(1)
	go func(fn func() error) {
		w.routinesLimiter <- struct{}{}
		if w.skipWhenError { // whether we should skip works execution when there is an error
			w.errMutex.RLock()
			if err := w.err; err == nil {
				w.errMutex.Unlock()
				w.execute(fn)
			} else {
				w.errMutex.Unlock()
			}
		} else {
			w.execute(fn)
		}

		w.wg.Done()
		<-w.routinesLimiter // free up one space to unblock one goroutine
	}(work)
}

func (w *ErrorWorkGroup) execute(work func() error) {
	if err := work(); err != nil {
		w.errMutex.Lock()
		w.err = err
		w.errMutex.Unlock()
	}
}
