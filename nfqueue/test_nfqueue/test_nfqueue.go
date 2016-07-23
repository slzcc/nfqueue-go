package main

import (
    "encoding/hex"
    "fmt"
    "github.com/chifflier/nfqueue-go/nfqueue"
    "os"
    "os/signal"
    "syscall"
)

func real_callback(payload *nfqueue.Payload) int {
    fmt.Println("Real callback")
    fmt.Printf("  id: %d\n", payload.Id)
    fmt.Printf("  mark: %d\n", payload.GetNFMark())
    fmt.Printf("  in  %d      out  %d\n", payload.GetInDev(), payload.GetOutDev())
    fmt.Printf("  Φin %d      Φout %d\n", payload.GetPhysInDev(), payload.GetPhysOutDev())
    fmt.Println(hex.Dump(payload.Data))
    fmt.Println("-- ")
    return nfqueue.NF_ACCEPT
}

func main() {
    q := new(nfqueue.Queue)

    q.SetCallback(real_callback)

    q.Init()
    defer q.Close()

    q.Unbind(syscall.AF_INET)
    q.Bind(syscall.AF_INET)

    q.CreateQueue(0)
    q.SetMode(nfqueue.NFQNL_COPY_PACKET)

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    go func(){
        for sig := range c {
            // sig is a ^C, handle it
            _ = sig
            q.Close()
            os.Exit(0)
            // XXX we should break gracefully from loop
        }
    }()

    // XXX Drop privileges here

    // XXX this should be the loop
    q.TryRun()

    fmt.Printf("hello, world\n")
}
