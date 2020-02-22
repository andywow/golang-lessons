package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	host    = flag.String("host", "127.0.0.1", "host address")
	port    = flag.Int("port", 80, "host port")
	timeout = flag.Int("timeout", 10, "connection timeout in seconds")
)

func readRoutine(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc, conn net.Conn) {
	defer wg.Done()
	defer cancel()

	scanner := bufio.NewScanner(conn)

	for {
		select {

		case <-ctx.Done():
			return

		default:
			if !scanner.Scan() {

				if scanner.Err() != nil {
					select {
					case <-ctx.Done():
						// no need to process error, cause context closed by us
						return
					default:
						log.Printf("Error from server: %s	\n", scanner.Err())
						break
					}

				} else {
					log.Println("Remote server closed connection")
				}
				return
			}
			text := scanner.Text()
			log.Printf("Response from server: %s", text)
		}

	}
}

func writeRoutine(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc, conn net.Conn) {
	defer wg.Done()
	defer cancel()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {

		case <-ctx.Done():
			return

		default:
			fmt.Print(">>>>>>>>>> ")
			if !scanner.Scan() {
				if scanner.Err() == nil {
					fmt.Println("Terminated by CTRL+D")
					conn.Close()
				} else {
					fmt.Printf("Error while reading stdin: %s", scanner.Err())
				}
				return
			}
			str := scanner.Text()
			log.Printf("Sent to server: %v\n", str)

			conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}

	}
}

func main() {
	flag.Parse()

	dialerContext, dialerCancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*timeout))
	defer dialerCancel()

	dialer := &net.Dialer{}
	log.Printf("Connecting to %s ...\n", *host)
	conn, err := dialer.DialContext(dialerContext, "tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatalf("error, while connecting: %s\n", err)
	}
	log.Println("Connected. Type your commands")

	abortContext, abortCancel := context.WithCancel(context.Background())
	defer abortCancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go readRoutine(abortContext, wg, abortCancel, conn)

	wg.Add(1)
	go writeRoutine(abortContext, wg, abortCancel, conn)

	wg.Wait()

}
