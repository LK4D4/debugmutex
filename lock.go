package debugmutex

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// Mutex is debug mutex, which prints logs and traces for locks
// implements sync.Locker interface
type Mutex struct {
	retries  int
	mu       sync.Mutex
	myMu     sync.RWMutex
	lockedAt string
}

// New returns new debug mutex
// retries parameter specify number of retries acquiring Mutex before fatal
// exit, if <=0, then wait forever
func New(retries int) sync.Locker {
	return &Mutex{retries: retries}
}

// Lock tries to lock mutex. Interval between attempts is 1 second.
// On each attempt stack trace and file:lino of previous Lock will be printed.
// Lock does os.Exit(1) after last attempt.
func (m *Mutex) Lock() {
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)
	log.Printf("Trying to acquire Lock at %s", caller)
	wait := make(chan struct{})
	go func() {
		m.mu.Lock()
		close(wait)
	}()
	seconds := 0
loop:
	for {
		select {
		case <-wait:
			break loop
		case <-time.After(1 * time.Second):
			m.myMu.RLock()
			seconds++
			if m.retries > 0 && seconds > m.retries {
				debug.PrintStack()
				log.Fatalf("Possible deadlock - can't acquire lock at %s for 5 second, locked by %s", caller, m.lockedAt)
			}
			log.Printf("Lock is stuck at %s, wait for lock from %s", caller, m.lockedAt)
			m.myMu.RUnlock()
		}
	}
	log.Printf("Lock acquired at %s", caller)
	debug.PrintStack()
	m.myMu.Lock()
	m.lockedAt = caller
	m.myMu.Unlock()
}

// Unlock unlocks mutex. It prints place in code where it was called and where
// mutex was locked.
func (m *Mutex) Unlock() {
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)
	m.myMu.RLock()
	log.Printf("Release Lock locked at %s, at %s", m.lockedAt, caller)
	m.myMu.RUnlock()
	m.mu.Unlock()
}
