package main

import "fmt"

func FibonacciIterative(n int) int {
	n1, n2, res := 0, 1, 1
	for i := 1; i < n; i++ {
		res = n1 + n2
		n1 = n2
		n2 = res
		fmt.Println(res)
	}
	return res
}

func FibonacciRecursive(n int) int {
	res := 1
	if n <= 1 {
		return n
	} else {
		res = FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
	}
	fmt.Println(res)
	return res
}

func IsPrime(n int) bool {
	res := true
	for i := 2; i <= n*n; i++ {
		if n%i == 0 {
			res = false
			break
		}
	}
	return res
}

func main() {
	fmt.Println("FibonacciIterative(10):", FibonacciIterative(10))
	fmt.Println("FibonacciRecursive(10):", FibonacciRecursive(10))
}
