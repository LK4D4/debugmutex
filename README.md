debugmutex
==========

Mutex for debugging deadlocks

# Usage

```
type Struct struct {
    sync.Locker
}

func New() *Struct {
    locker := &sync.Mutex{}
    if os.Getenv("DEBUG") != "" {
        // will crash program with traceback and file:line where deadlock is
        // occured after five tries to acquire mutex with 1 second gap between.
        locker = debugmutex.New(5, true)
    }
    return &Struct{Locker: locker}
}
```

For logging used `github.com/Sirupsen/logrus`. You can set debug logging with

```
logrus.SetLevel(logrus.DebugLevel)
```
