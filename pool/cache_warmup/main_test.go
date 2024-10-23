package main_test

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}

		defer server.Close()

		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
			}
			connectToService()

			fmt.Fprintln(conn, "")
			conn.Close()
		}

	}()
	return &wg
}

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool{
		New: connectToService,
	}

	for i := 0; i < 10; i++ {
		p.Put(p.New)
	}
	return p
}

func startNetworkDaemonWithCache() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			svcConn := connPool.Get()
			fmt.Fprintln(conn, "")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()
	return &wg
}

func init() {
	daemonStarted := startNetworkDaemonWithCache()
	daemonStarted.Wait()
}

func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatalf("cannot dia; jost: %v", err)
		}
		if _, err := io.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}
