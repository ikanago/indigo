package main

import "fmt"

func main() {
	fmt.Println(".arch armv8-a")
	fmt.Println(".text")
	fmt.Println(".align	2")
	fmt.Println(".globl _main")
	fmt.Println("_main:")
	fmt.Println("\tmov x0, #0")
	fmt.Println("\tret")
}
