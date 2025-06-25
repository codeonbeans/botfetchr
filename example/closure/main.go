package main

import "fmt"

func VariableLoop() {
	f := make([]func(), 3)
	for i := 0; i < 3; i++ {
		// closure over variable i
		f[i] = func() {
			fmt.Println(i)
		}
	}
	fmt.Println("VariableLoop")
	for _, f := range f {
		f()
	}
}

func ValueLoop() {
	f := make([]func(), 3)
	for i := 0; i < 3; i++ {
		i := i
		// closure over value of i
		f[i] = func() {
			fmt.Println(i)
		}
	}
	fmt.Println("ValueLoop")
	for _, f := range f {
		f()
	}
}

func VariableRange() {
	f := make([]func(), 3)
	for i := range f {
		// closure over variable i
		f[i] = func() {
			fmt.Println(i)
		}
	}
	fmt.Println("VariableRange")
	for _, f := range f {
		f()
	}
}

func ValueRange() {
	f := make([]func(), 3)
	for i := range f {
		i := i
		// closure over value of i
		f[i] = func() {
			fmt.Println(i)
		}
	}
	fmt.Println("ValueRange")
	for _, f := range f {
		f()
	}
}

func main() {
	VariableLoop()
	ValueLoop()
	VariableRange()
	ValueRange()
}
