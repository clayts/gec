package main

import "fmt"

func main() {
	f := 1.0
	i := 1
	for {
		fmt.Println(f, i)
		if int(f) != i {
			break
		}
		f *= 2
		i *= 2
	}
}
