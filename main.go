package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"pipe-network/internal/api"
	"pipe-network/internal/hydraulics"
	"pipe-network/internal/leak"
	"pipe-network/internal/network"
)

//go:embed web/*
var webFS embed.FS

//go:embed example/*
var exampleFS embed.FS

func main() {
	httpAddr := flag.String("http", "", "serve web UI and API on this address, e.g. :8080")
	flag.Parse()
	if *httpAddr != "" {
		wsub, err := fs.Sub(webFS, "web")
		if err != nil {
			log.Fatal(err)
		}
		esub, err := fs.Sub(exampleFS, "example")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("pipe-network listening on %s", *httpAddr)
		log.Fatal(http.ListenAndServe(*httpAddr, api.New(wsub, esub)))
	}
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: pipe-network [-http :8080] | solve <file.json> | leak <file.json> <node> <extraDemand>")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "solve":
		cmdSolve(os.Args[2:])
	case "leak":
		cmdLeak(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command %s\n", os.Args[1])
		os.Exit(2)
	}
}

func cmdSolve(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: solve <file.json>")
		os.Exit(2)
	}
	n, err := loadNet(args[0])
	if err != nil {
		log.Fatal(err)
	}
	ig, err := network.Prepare(n)
	if err != nil {
		log.Fatal(err)
	}
	res, err := hydraulics.Solve(ig, hydraulics.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	printResult(res, ig)
}

func cmdLeak(args []string) {
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: leak <file.json> <node> <extraDemand>")
		os.Exit(2)
	}
	n, err := loadNet(args[0])
	if err != nil {
		log.Fatal(err)
	}
	extra := parseFloatArg(args[2])
	spec := leak.Spec{NodeID: args[1], Mode: leak.ModeDemand, ExtraDemand: extra}
	ig, err := network.Prepare(n)
	if err != nil {
		log.Fatal(err)
	}
	res, err := leak.SolveWithLeak(n, spec, hydraulics.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	printResult(res, ig)
}

func loadNet(path string) (*network.Network, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return network.ParseJSON(data)
}

func parseFloatArg(s string) float64 {
	var v float64
	fmt.Sscanf(s, "%g", &v)
	return v
}

func printResult(res *hydraulics.Result, ig *network.Indexed) {
	fmt.Println("flow (m3/s):")
	for _, p := range ig.Pipes {
		fmt.Printf("  %s = %.6g\n", p.ID, res.Flow[p.ID])
	}
	fmt.Println("head (m):")
	for _, id := range ig.IDs {
		fmt.Printf("  %s = %.6g\n", id, res.Head[id])
	}
	fmt.Printf("method=%s iterations=%d converged=%v\n", res.Method, res.Iterations, res.Converged)
	fmt.Printf("max_mass_residual=%.3e loop_closure=%.3e\n",
		hydraulics.MaxMassResidual(ig, res), hydraulics.LoopClosure(ig, res))
}
