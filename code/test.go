package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	// testChan()
	// testNewMake()
	// testReflect()
	// result := isMatch("aa","a*")
	// fmt.Println(result)
	testChan3()
}

func testChan3() {
	ch := make(chan int, 1)
	ch <- 1

	// for val := range ch {
	// 	fmt.Println(val)
	// 	if val == 5 {
	// 		close(ch)
	// 	} else {
	// 		ch <- val + 1
	// 	}
	// }

	for {
		select {
			case val, ok := <- ch: {
				if !ok {
					fmt.Println("结束")
					goto finish
				} else {
					fmt.Println(val)
					if val < 5 {
						ch <- val + 1
					} else {
						close(ch)
					}
				}
			}
		}
	}

	finish:
	fmt.Println("结束2")
}

func testChan2() {
	ch1 := make(chan int, 10)
	ticker := time.NewTicker(1000 * time.Second)


	for {
		select {
			case val := <- ch1:{
				fmt.Println(val)
			}
			case <- ticker.C: {
				ch1 <- 666
			}
		}
	}
}

func isMatch(s string, p string) bool {
	// dp[i][j] ---  s的前i个字符 能否匹配p的前j个字符
	// dp[0][0] --- 1
	// dp[0][j] --- 如果 p[j] 是 * 且dp[i][j-2]为1 
	// dp[i][0] --- 0 
	// dp[i][j] --- 
	lenS, lenP := len(s), len(p)
	dp := make([][]uint8, lenS + 1)
	for i := range dp {
		dp[i] = make([]uint8, lenP + 1)
	}

	// fmt.Println(dp)
	dp[0][0] = 1
	// dp[0][j] --- 如果 p[j] 是 * 且dp[i][j-2]为1 
	for j:= 2; j <= lenP; j ++ {
		if p[j - 1] == '*' {
			dp[0][j] = dp[0][j - 2]
		}
	}

	for i := 1; i <= lenS; i ++ {
		for j := 1; j <= lenP; j ++ {
			// 当前两个字符相同 
			if s[i - 1] == p[j - 1] {
				dp[i][j] = dp[i - 1][j - 1]
			} else {
				switch p[j - 1] {
				case '*': {
					// fmt.Printf("case '*': i:%d j:%d p[j-2]:%s s[i-1]:%s\n", i, j, string(p[j - 2]),  string(s[i - 1]))
					if j >= 2 {
                        if p[j - 2] == '.' || p[j - 2] == s[i - 1] {
							// fmt.Println(" if p[j - 2] == '.' || p[j - 2] == s[i - 1] ")
							val1 := dp[i-1][j] // p最后的 x*可以顶替s末尾的x  如果 这之前的s和p都匹配上了 那肯定没问题
							val2 := dp[i][j-2]	// p最后的 x*可以当作不存在  如果此时s已经可以和p末尾出去x*的部分匹配 那也没问题
						
							if val1 == 1 || val2 == 1 {
								dp[i][j] = 1
							}
                        } else {
                            dp[i][j] = dp[i][j - 2]
                        }
					} else {
                        dp[i][j] = 0
                    }
				}
				case '.': {
					dp[i][j] = dp[i - 1][j - 1]
				}
				default: {
					dp[i][j] = 0
				}
				}
			}

		}
	}

	for i:=0; i <=lenS; i ++ {
		fmt.Println(dp[i])
	}
	return dp[lenS][lenP] == 1
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

	// val := reflect.ValueOf(animal)
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