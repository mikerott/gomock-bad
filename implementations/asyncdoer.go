package implementations

import (
	"fmt"
	"sync"

	"github.com/mikerott/gomock-bad/interfaces"
)

const maxGoRoutines = 10

var mutex = &sync.RWMutex{}

type semaphore chan struct{}

func (s semaphore) acquire() {
	s <- struct{}{}
}

func (s semaphore) release() {
	<-s
}

func (s semaphore) close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	close(s)
	return nil
}

type Processor struct {
	ThingConsumer interfaces.ThingConsumer
	Jobs          int
}

func (p *Processor) Process() {

	stringAccumulatorChan := make(chan string) // accumulate the strings here

	accumulatedStrings := []string{}

	go func() { // channel reader
		for {
			select {
			case s := <-stringAccumulatorChan:
				fmt.Printf("MIKE: got: %s\n", s)
				mutex.Lock()
				accumulatedStrings = append(accumulatedStrings, s)
				mutex.Unlock()
				fmt.Printf("MIKE: size: %d\n", len(accumulatedStrings))
			}
		}
	}()

	var sema semaphore = make(chan struct{}, maxGoRoutines)
	var wg sync.WaitGroup

	wg.Add(p.Jobs)

	for i := 0; i < p.Jobs; i++ {
		sema.acquire()
		go func() {
			defer sema.release()
			defer wg.Done()

			things, err := p.ThingConsumer.ConsumeThings()
			if err == nil {
				for _, thing := range things {
					stringAccumulatorChan <- thing
				}
			}
		}()
	}

	wg.Wait()

	mutex.Lock()
	fmt.Printf("MIKE: after wg.Wait(): %d\n", len(accumulatedStrings))
	mutex.Unlock()

}
