package main

import (
        "fmt"
        "time"

        "github.com/iowar/throttle"
)

var (
        interval = time.Second
)

func watcher(throttle *throttle.Throttle) {
        for {
                <-throttle.Update
                fmt.Println(time.Now())
        }
}

func main() {
        th := throttle.NewThrottle(interval)
        th.Start()
        go watcher(th)

        time.Sleep(time.Second * 5)

        // change interval
        th.ChangeInterval(time.Second * 2)
        time.Sleep(time.Second * 5)
}

// sample output:
// 2019-04-03 15:32:05.272257431 +0300 +03 m=+1.001780388
// 2019-04-03 15:32:06.272450811 +0300 +03 m=+2.001973724
// 2019-04-03 15:32:07.272291136 +0300 +03 m=+3.001814122
// 2019-04-03 15:32:08.272286699 +0300 +03 m=+4.001809752
// 2019-04-03 15:32:09.272253546 +0300 +03 m=+5.001776560
// 2019-04-03 15:32:11.272973758 +0300 +03 m=+7.002496712
// 2019-04-03 15:32:13.273010406 +0300 +03 m=+9.002533361
