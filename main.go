package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {

	group, context := errgroup.WithContext(context.Background())

	start := make(chan int)
	stop := make(chan int)

	srv := &http.Server{Addr: ":8080"}

	group.Go(func() error {
		<-start
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello world!"))
		})
		fmt.Printf("server start")
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("ListenAndServe Err")
		}
		fmt.Printf("server over")
		return nil
	})

	group.Go(func() error {
		<-stop
		fmt.Printf("over")
		srv.Shutdown(context)
		return nil
	})

	group.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT)

		for {
			sig := <-sigs
			if sig == syscall.SIGINT {
				start <- 1
			} else if sig == syscall.SIGQUIT {
				stop <- 1
				break
			}
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}

}
