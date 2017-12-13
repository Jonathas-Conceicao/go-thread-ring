package main

import "fmt"
import "os"
import "strconv"
import "runtime"
import "sync"

const nThreads = 503

type Thread struct {
  id   int
  Op   int
  mtx sync.Mutex
  next *Thread
}

func (t *Thread) pass() {
  if t.Op == 0 {
    done <- t.id
  } else {
    t.next.Op = t.Op -1
    t.next.mtx.Unlock()
  }
}

func (t *Thread) run() {
  for {
    t.mtx.Lock()
    t.pass()
    runtime.Gosched()
  }
}

func (t *Thread) Init(id int, n *Thread) {
  t.id   = id
  t.next = n
  t.mtx.Lock()
  go t.run()
}

func (t *Thread) Start(ops int) {
  t.Op = ops
  t.mtx.Unlock()
}

var done chan int = make(chan int)

func main() {
  ops, _ := strconv.Atoi(os.Args[1]) // os.Args[0] is program's executable, os.Args[1] is arg1
  runtime.GOMAXPROCS(1)
  var threads [nThreads]Thread
  for i := range(threads) {
    threads[i].Init(i+1, &threads[(i+1) % nThreads])
  }
  threads[0].Start(ops)
  fmt.Println(<-done)
}
