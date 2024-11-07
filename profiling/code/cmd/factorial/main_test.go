package main

import "testing"

const factorialInput = 100000

func BenchmarkFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		factorial(factorialInput)
	}
}

func BenchmarkFactorialOptimised(b *testing.B) {
	for i := 0; i < b.N; i++ {
		factorialOptimised(factorialInput)
	}
}
