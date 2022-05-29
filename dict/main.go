package main

import (
	"fmt"

	"example.com/dict/mydict"
)

func main() {
	dictionary := mydict.Dictionary{}
	err := dictionary.Add("first", "First Word")
	if err != nil {
		fmt.Println(err)
	}

	err = dictionary.Update("firat", "Second")
	if err != nil {
		fmt.Println(err)
	}

	err = dictionary.Delete("firat")
	if err != nil {
		fmt.Println(err)
	}

	def, err := dictionary.Search("first")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(def)
}
