package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/trieudh58/ethereum-parser/ethereum_client"
	"github.com/trieudh58/ethereum-parser/parser"
	"github.com/trieudh58/ethereum-parser/store"
	"golang.org/x/sync/errgroup"
)

var (
	apiPort         int
	endpoint        string
	numberOfWorkers int
	interval        int
)

func init() {
	flag.IntVar(&apiPort, "api.port", 3000, "port to run api server")
	flag.StringVar(&endpoint, "rpc.endpoint", "https://cloudflare-eth.com", "ethereum RPC endpoint")
	flag.IntVar(&numberOfWorkers, "workers", 2, "number of workers")
	flag.IntVar(&interval, "interval", 10, "time delay between rpc calls")
}

func main() {
	flag.Parse()

	rootCtx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{}, 2)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	store := store.NewMemoryStore()
	store.SaveLastParsedBlock(18078479) // skip ahead to block 18078480

	client := ethereum_client.NewEthereumClient(endpoint)
	parser := parser.NewEthereumParser(store)

	// run background job
	go func() {
		taskQ := make(chan int, numberOfWorkers)

		g, ctx := errgroup.WithContext(rootCtx)
		// create workers
		for id := 1; id <= numberOfWorkers; id++ {
			id := id
			g.Go(func() error {
				worker := NewWorker(id, store, &client)
				log.Printf("worker id %d spawned", id)
				for {
					select {
					case blockNo := <-taskQ:
						worker.Work(ctx, blockNo)
					case <-ctx.Done():
						log.Printf("worker %d received quit signal", id)
						return nil
					}
				}
			})
		}

		// gradually generate tasks and send to queue
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					parsedBlock := store.GetLastParsedBlock()
					currentBlockHeight, err := client.GetBlockNumber(ctx)
					if err != nil {
						log.Println(err)
						continue
					}

					for i := 0; i < numberOfWorkers; i++ {
						nextBlock := parsedBlock + i + 1
						if nextBlock < currentBlockHeight {
							taskQ <- nextBlock
						}
					}
				case <-quit:
					return
				}
			}
		}()

		// wait for the workers to finish
		if err := g.Wait(); err != nil {
			log.Println(err)
		}
		done <- struct{}{}
	}()

	// run api server
	api := NewApiServer(parser)
	go func() {
		go func() {
			log.Printf("api server is running at port %d", apiPort)
			if err := api.ListenAndServe(apiPort); err != nil {
				log.Println(err)
				done <- struct{}{}
			}
		}()

		go func() {
			<-rootCtx.Done()
			log.Println("api server received quit signal")
			api.Stop(rootCtx)
			done <- struct{}{}
		}()
	}()

	<-quit

	log.Println("cleaning up")
	cancel()

	<-done
}
