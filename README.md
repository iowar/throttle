## Golang Dynamic Throttle
[![GoDoc](https://godoc.org/github.com/iowar/throttle?status.svg)](https://godoc.org/github.com/iowar/throttle)

In some cases it may be necessary to change the throttle at run time.
The 'Update' channel in throttle can be used to provide triggering like time.Ticker.

## Installation

```sh
$ go get -u github.com/iowar/throttle
```

## Usage

Simple usage example:
```go
package main

import (
        "fmt"
        "time"

        "github.com/iowar/throttle"
)

var (
        kill = make(chan struct{})
)

func main() {
        th := throttle.NewThrottle(time.Second)
        th.Start()

        go func() {
                time.Sleep(time.Second * 6)
                th.ChangeInterval(time.Second * 3)
                fmt.Println("Changing interval.")
        }()

        go func() {
                time.Sleep(time.Second * 15)
                th.Stop()
                fmt.Println("Stopping throttle.")
        }()

        go func() {
                time.Sleep(time.Second * 24)
                th.Start()
                fmt.Println("Starting throttle.")
        }()

        go func() {
                time.Sleep(time.Second * 33)
                kill <- struct{}{}
                fmt.Println("Finish.")
        }()

        go func() {
                for {
                        <-th.Update
                        fmt.Println(time.Now())
                }
        }()

        <-kill
}

```
