package debugmutex

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetOutput(ioutil.Discard)
}

func TestLockUnlock(t *testing.T) {
	l := New(1, true)
	l.Lock()
	l.Unlock()
}

func TestDoubleLockPanic(t *testing.T) {
	l := New(1, true)
	lock := l.(*Mutex)
	lock.errFunc = logrus.Panicf
	l.Lock()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Double lock shoud fatal")
		}
		e, ok := r.(*logrus.Entry)
		if !ok {
			t.Fatalf("Unexpected type from panic: %T", r)
		}
		if !strings.Contains(e.Message, "Possible deadlock") {
			t.Fatalf("Unexpected message from panic: %s", e.Message)
		}
	}()

	l.Lock()
}
