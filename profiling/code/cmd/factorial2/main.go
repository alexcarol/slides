package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime/pprof"
)

func factorial(n int64) *big.Int {
	result := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}
	return result
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var n int64 = 100000
	result := factorial(n)
	fmt.Printf("Factorial of %d: %s\n", n, result.String())
}
