package main

import (
	"fmt"
	"reflect"
)

func main() {
	// testChan()
	// testNewMake()
	testReflect()
}

type  Animal interface {
	Say()
	Eat()
}

type Dog struct {
	Name string `myTag:"wuhu~"`
}
type Cat struct {
}

func (d Dog) Say() {
	fmt.Println("汪汪汪")
}

func (d Dog) Eat() {
	fmt.Println("吃骨头")
}

func (c Cat) Say() {
	fmt.Println("喵喵喵")
}

func (c Cat) Eat() {
	fmt.Println("吃小虾子")
}

func testReflect() {
	someAnimal := Dog{}
	handle(someAnimal)
}

func handle(animal Animal) {
	// 如果是狗  就让它叫
	// 如果是猫  就让它吃

	val := reflect.ValueOf(animal)
	typ := reflect.TypeOf(animal)

	nameOfType := typ.Name()
	switch nameOfType {
	case "Cat" :{
		break;
	}
	case "Dog" :{
		break
	}
	}
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