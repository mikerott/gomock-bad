package implementations

import (
	"fmt"
	"sync"
	"time"

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

	var sema semaphore = make(chan struct{}, maxGoRoutines)
	var producerWG sync.WaitGroup
	done := make(chan bool)

	stringAccumulatorChan := make(chan string) // accumulate the strings here

	accumulatedStrings := []string{}

	go func() { // channel reader
		for {
			select {
			case s := <-stringAccumulatorChan:
				fmt.Printf("MIKE: got: %s\n", s)
				time.Sleep(1 * time.Second)
				mutex.Lock()
				accumulatedStrings = append(accumulatedStrings, s)
				mutex.Unlock()
				fmt.Printf("MIKE: size: %d\n", len(accumulatedStrings))
			case <-done:
				return
			}
		}
	}()

	producerWG.Add(p.Jobs)

	for i := 0; i < p.Jobs; i++ {
		sema.acquire()
		go func() {
			defer sema.release()
			defer producerWG.Done()

			things, err := p.ThingConsumer.ConsumeThings()
			if err == nil {
				for _, thing := range things {
					stringAccumulatorChan <- thing
				}
			}
		}()
	}

	producerWG.Wait()
	done <- true

	mutex.Lock()
	fmt.Printf("MIKE: after producerWG.Wait(): %d\n", len(accumulatedStrings))
	mutex.Unlock()

}
