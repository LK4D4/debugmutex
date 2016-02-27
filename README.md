debugmutex
==========

[![Build Status](https://travis-ci.org/LK4D4/debugmutex.svg?branch=master)](https://travis-ci.org/LK4D4/debugmutex)
[![GoDoc](https://godoc.org/github.com/LK4D4/debugmutex?status.svg)](https://godoc.org/github.com/LK4D4/debugmutex)

Mutex for debugging deadlocks. It can find non-obvious deadlocks in systems
with heavy `sync.Mutex` usage. I found many deadlocks in Docker with it.

# Usage

```go
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

For logging `github.com/Sirupsen/logrus` is used. You can set debug logging with

```go
logrus.SetLevel(logrus.DebugLevel)
```
