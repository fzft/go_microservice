package main

import "fmt"

type Person struct {
	id   int
	name string
}

func main() {
	fmt.Printf("before %s ", people)
	person1 := *people[0]
	person1.name = "a-a"
	fmt.Printf("after  %s", people)

}

var people = []*Person{
	&Person{
		id:   0,
		name: "a",
	},
	&Person{
		id:   1,
		name: "b",
	},
}
