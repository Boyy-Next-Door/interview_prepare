package main

import "fmt"

func main() {
	// testChan()
	testNewMake()
}

func testChan() {
	ch := make(chan string)
	ch2 := make(chan string, 10)
	ch2 <- "hello"
	data := <-ch2

	fmt.Println(" <- ch2   returned!", data)

	<- ch
	fmt.Println(" <- ch   returned!")
}


func testNewMake() {
	var u1 = new(*User)
	var u2 = new(User)
	var u3 = User{}
	var u4 User

	fmt.Println(u1)	// 0xc000006028  零值User的指针
	fmt.Println(&u1)// &{0  <nil>}	零值User
	fmt.Println(u2)	// &{0  <nil>}
	fmt.Println(u3)	// {0  <nil>}
	fmt.Println(u4)	// {0  <nil>}
}

type User struct {
	Age int32 
	Name string
	Ch chan string
}