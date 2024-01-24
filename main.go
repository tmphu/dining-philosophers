package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

const (
	NUM_PHILOS             = 5
	NUM_CHOPSTICKS         = 5
	EAT_TIMES              = 3
	PHILO_EAT_CONCURRENTLY = 2
)

type ChopS struct {
	sync.Mutex
}

type Philo struct {
	id                    int
	leftCS, rightCS       *ChopS
	leftCsIdx, rightCsIdx int
}

type Host struct {
	permissionChan chan struct{}
}

func NewHost(numEatConcurrent int) *Host {
	return &Host{
		permissionChan: make(chan struct{}, numEatConcurrent),
	}
}

func (h *Host) AskPermission() {
	<-h.permissionChan
}

func (h *Host) ReleasePermission() {
	h.permissionChan <- struct{}{}
}

func (p Philo) eat(host *Host) {
	// when done eating, it should send data to channel to allow other goroutine eating
	defer host.ReleasePermission()
	defer wg.Done()

	// block operation eating here, waiting for permission which is triggered by sending data to channel
	host.AskPermission()

	for i := 0; i < EAT_TIMES; i++ {
		p.leftCS.Lock()
		p.rightCS.Lock()

		fmt.Println("starting to eat", p.id)

		time.Sleep(2 * time.Second)

		fmt.Println("finishing eating", p.id)

		p.leftCS.Unlock()
		p.rightCS.Unlock()
	}
}

func main() {
	CSticks := make([]*ChopS, NUM_CHOPSTICKS)
	for i := 0; i < NUM_CHOPSTICKS; i++ {
		CSticks[i] = new(ChopS)
	}

	philos := make([]*Philo, NUM_PHILOS)
	for i := 0; i < NUM_PHILOS; i++ {
		philos[i] = &Philo{
			id:         i + 1,
			leftCS:     CSticks[i],
			rightCS:    CSticks[(i+1)%NUM_CHOPSTICKS],
			leftCsIdx:  i,
			rightCsIdx: (i + 1) % NUM_CHOPSTICKS,
		}
	}

	host := NewHost(PHILO_EAT_CONCURRENTLY)

	wg.Add(NUM_PHILOS)

	// start 5 goroutines, but they will stand by waiting for permission due to host.AskPermission()
	for i := 0; i < NUM_PHILOS; i++ {
		go philos[i].eat(host)
	}

	// send 2 signals to channel, meaning only 2 philos have permission to eat concurrently
	host.permissionChan <- struct{}{}
	host.permissionChan <- struct{}{}

	wg.Wait()
}
