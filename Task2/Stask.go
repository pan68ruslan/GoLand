package main

// Stack структура
type Stack struct {
	items []int
}

func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (int, bool) {
	if len(s.items) == 0 {
		return 0, false // якщо стек порожній
	}
	lastIndex := len(s.items) - 1
	element := s.items[lastIndex]
	s.items = s.items[:lastIndex]
	return element, true
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

//func main2() {
//	stack := Stack{}
//
//	stack.Push(10)
//	stack.Push(20)
//	stack.Push(30)
//
//	x, _ := stack.Peek()
//	fmt.Println("Peek:", x) // 30
//	x, _ = stack.Pop()
//	fmt.Println("Pop:", x) // 30
//	//fmt.Println("Pop:", stack.Pop())         // 20
//	fmt.Println("IsEmpty:", stack.IsEmpty()) // false
//}
