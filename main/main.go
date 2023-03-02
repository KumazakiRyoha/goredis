package main

import "fmt"

func main() {
	str := "*3\r\ndsfjisdf"
	str1 := string(str[1 : len(str)-2])
	fmt.Println(str1)
}
