package main

import (
	"fmt"
	"strconv"
)

func FibonacciIterative(n int) int {
	n1, n2 := 0, 1
	for i := 1; i < n; i++ {
		n2 += n1
		n1 = n2 - n1
	}
	return n2
}

func FibonacciRecursive(n int) int {
	res := 1
	if n <= 1 {
		return n
	} else {
		res = FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
	}
	return res
}

func IsPrime(n int) bool {
	for i := 2; i <= n; i++ {
		if n != i && n%i == 0 {
			return false
		}
	}
	return true
}

func IsBinaryPalindrome(n int) bool {
	ss := []rune(strconv.FormatInt(int64(n), 2))
	l := len(ss)
	for i := 0; i < l/2; i++ {
		if ss[i] != ss[l-i-1] {
			return false
		}
	}
	return true
}

func ValidParentheses(s string) bool {
	if len(s) == 0 {
		return false
	}
	ss := []rune(s)
	stack := make([]rune, len(ss))
	idx := 0
	for _, c := range ss {
		isOpenBracket := c == '(' || c == '[' || c == '{'
		isCloseBracket := c == ')' || c == ']' || c == '}'
		if !(isOpenBracket || isCloseBracket) { // check if a current symbol is a bracket
			return false
		}
		if isOpenBracket { // add (push) a current open bracket
			stack[idx] = c
			idx++
		} else { // check type and remove (pop) a previous open bracket
			if idx > 0 {
				idx--
				switch c {
				case ')':
					if stack[idx] != 40 { //'('
						return false
					}
				case ']':
					if stack[idx] != 91 { //'['
						return false
					}
				case '}':
					if stack[idx] != 123 { //'{'
						return false
					}
				}
			} else {
				return false
			}
		}
	}
	return idx == 0
}

func Increment(num string) int {
	digit := 0
	for _, n := range num {
		if n == '0' || n == '1' {
			digit <<= 1
			if n == '1' {
				digit++
			}
		} else {
			return 0
		}
	}
	digit++
	return digit
}

func main() {
	fmt.Println("FibonacciIterative(10):", FibonacciIterative(10))
	fmt.Println("FibonacciRecursive(10):", FibonacciRecursive(10))

	fmt.Println("IsPrime(2):", IsPrime(2))                         // true
	fmt.Println("IsPrime(15):", IsPrime(15))                       // false
	fmt.Println("IsPrime(29):", IsPrime(29))                       // true
	fmt.Println("IsPrime(822):", IsPrime(822))                     // false
	fmt.Println("IsPrime(823):", IsPrime(823))                     // true
	fmt.Println("IsBinaryPalindrome(7):", IsBinaryPalindrome(7))   // true  (111)
	fmt.Println("IsBinaryPalindrome(6):", IsBinaryPalindrome(6))   // false (110)
	fmt.Println("IsBinaryPalindrome(9):", IsBinaryPalindrome(9))   // true  (1001)
	fmt.Println("IsBinaryPalindrome(13):", IsBinaryPalindrome(13)) // false (1101)

	fmt.Println(`ValidParentheses("[]{}()"):`, ValidParentheses("[]{}()")) // true
	fmt.Println(`ValidParentheses("[{]}"):`, ValidParentheses("[{]}"))     // false
	fmt.Println(`ValidParentheses("[[[]]]"):`, ValidParentheses("[[[]]]")) // true
	fmt.Println(`ValidParentheses("[[[]]"):`, ValidParentheses("[[[]]"))   // false

	fmt.Println(`Increment("101") ->`, Increment("101"))   // 6
	fmt.Println(`Increment("10a") ->`, Increment("10a"))   // 0
	fmt.Println(`Increment("1000") ->`, Increment("1000")) // 9
	fmt.Println(`Increment("1111") ->`, Increment("1111")) // 16
}
