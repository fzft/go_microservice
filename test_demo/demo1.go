package main

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Person struct {
	id   int
	name string
}

func main() {
	//fmt.Printf("before %s ", people)
	//person1 := *people[0]
	//person1.name = "a-a"
	//fmt.Printf("after  %s", people)
	var a  *status.Status
	for i := 1; i < 2; i++ {
		a = status.New(codes.AlreadyExists, "Unable to subscribe for currency as subscription already exists")
		fmt.Println(a)
	}
	fmt.Println(a)

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
