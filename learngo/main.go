package main

import (
	"fmt"
	"strings"
)


func multiply(a , b int) int {
	return a * b 
}

func lenAndUpper(name string) (int, string) {
	return len(name), strings.ToUpper(name)
}

func repeatMe(words ...string) {
	fmt.Println(words)
}

func lenAndLower(name string) (length int, lowercase string) {
	defer fmt.Println("I'm done")
	length = len(name)
	lowercase = strings.ToLower(name)
	return
}	

func superAdd(numbers ...int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}

func canIDrink(age int) bool {
	if koreanAge := age + 2; koreanAge < 20 {
		return false
	}
	return true

}

func canIGo(time int) bool {
	switch {
	case time > 21:
		return true
	case time <=21:
		return false
	}

	return false
}

//struct
type person struct {
	name string
	age int
	favFood []string
}

func main() {
	//var name string = "Shin"
	name := "Shin"
	name = "Young Hoon"
	fmt.Println(name)
	fmt.Println(multiply(123, 123))
	fmt.Println(lenAndUpper(name))
	fmt.Println(lenAndLower(name))
	repeatMe("asdf", "dfa", "fdasdf")
	total := superAdd(1, 2, 3, 4, 5, 6, 123)
	println(total)
	fmt.Println(canIDrink(16))
	fmt.Println(canIGo(22))

	//pointer
	a := 20
	b := &a
	fmt.Println(b)
	fmt.Println(*b)
	*b = 22
	fmt.Println(a)

	//array: 크기 정해짐
	names := [5]string{"Shine", "Baek", "Hong"}
	names[3] = "blah"
	fmt.Println(names)
	//slice: 무한. append 사용
	slice := [] string {"blah"}
	slice = append(slice, "blah")
	fmt.Println(slice)

	//map [key] value
	YH := map[string] string {"name": "YoungHoon", "age": "22" }

	for key, value := range YH{
		fmt.Println(key)
		fmt.Println(value)
	}
	
	favFood := [] string {"chicken", "pizza"}
	Shin := person{name: "YoungHoon", age: 21, favFood: favFood}
	fmt.Println(Shin.name)

}