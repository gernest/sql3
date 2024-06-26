// Copyright 2021 Molecula Corp. All rights reserved.
package planner_test

import (
	"fmt"
	"net"
	"os"
	"testing"

	"net/http"
	_ "net/http/pprof"
)

func TestMain(m *testing.M) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	fmt.Printf("sql3/planner/ TestMain: online stack-traces: curl http://localhost:%v/debug/pprof/goroutine?debug=2\n", port)
	go func() {
		err := http.Serve(l, nil)
		if err != nil {
			panic(err)
		}
	}()
	os.Exit(m.Run())
}
