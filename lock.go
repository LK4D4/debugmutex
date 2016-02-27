package debugmutex

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

// Mutex is debugging mutex, which prints logs and traces for locks.
// Mutex implements sync.Locker interface.
type Mutex struct {
	retries  int
	mu       sync.Mutex
	myMu     sync.RWMutex
	lockedAt string
	errFunc  func(format string, args ...interface{})
}

// New returns new debug mutex.
// retries specify the number of retries acquiring Mutex before error
// message, if <=0, then retry forever.
// If fatal is true, then program will exit by os.Exit(1).
func New(retries int, fatal bool) sync.Locker {
	var errF func(format string, args ...interface{})
	errF = logrus.Errorf
	if fatal {
		errF = logrus.Fatalf
	}
	return &Mutex{
		retries: retries,
		errFunc: errF,
	}
}

// Lock tries to lock mutex. Interval between attempts is 1 second.
// On each attempt stack trace and file:line of a previous Lock printed.
// Lock does os.Exit(1) after last attempt if fatal parameter was true.
func (m *Mutex) Lock() {
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)
	logrus.Debugf("Trying to acquire Lock at %s", caller)
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
				logrus.Errorf("Stack:\n%s", stack())
				m.errFunc("Possible deadlock - can't acquire lock at %s for 5 second, locked by %s", caller, m.lockedAt)
			}
			logrus.Debugf("Lock is stuck at %s, wait for lock from %s", caller, m.lockedAt)
			m.myMu.RUnlock()
		}
	}
	logrus.Debugf("Lock acquired at %s", caller)
	logrus.Debugf("Stack:\n%s", stack())
	m.myMu.Lock()
	m.lockedAt = caller
	m.myMu.Unlock()
}

// Unlock unlocks the mutex. It prints place in the code where it was called
// and where mutex was locked.
func (m *Mutex) Unlock() {
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)
	m.myMu.RLock()
	logrus.Debugf("Release Lock locked at %s, at %s", m.lockedAt, caller)
	m.myMu.RUnlock()
	m.mu.Unlock()
}

func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], false)])
}
