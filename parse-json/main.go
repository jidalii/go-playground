package main

import (
	"encoding/json"
	"fmt"
)

type Employee struct {
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Age      int     `json:"age,omitempty"`
	Salary   *uint64 `json:"salary,omitempty"`
}

type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
}

func ObjectToJSON() {
	employee1 := Employee{
		Name:     "Alice",
		Position: "Manager",
		Age:      30,
		Salary:   new(uint64),
	}
	*employee1.Salary = 5000000
	jsonData, err := json.Marshal(employee1)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func JSONToObject() {
	jsonStr := `{"name":"Bob","position":"Developer","age":25,"salary":4000000}`
	var employee2 Employee
	err := json.Unmarshal([]byte(jsonStr), &employee2)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	fmt.Println(employee2)
}

func MapToJSON() {
	jsonStr := `{"name":"Bob","position":"Developer","age":25,"salary":4000000}`
	var m map[string]any
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	// lots of codes for maintaining the map
	name, ok := m["name"].(string)
	if !ok {
		fmt.Println("Error type assertion")
		return
	}
	fmt.Println(name)
}

func ZeroValueConfusion() {
	str := `{"name":"Bob","position":"Developer","age":25,"salary":100}`
	var p Employee
	_ = json.Unmarshal([]byte(str), &p)
	fmt.Printf("%+v\n", p)

	str2 := `{"name":"Bob","position":"Developer","age":25}`
	var p2 Employee
	_ = json.Unmarshal([]byte(str2), &p2)
	fmt.Printf("%+v\n", p2)
	if p2.Salary == nil {
		fmt.Println("p2.Salary is nil")
	}
}

func TagMarshal() {
	p := Employee{
		Name:   "Bob",
		Salary: new(uint64),
	}
	*p.Salary = 100
	output, _ := json.MarshalIndent(p, "", "  ")
	println(string(output))
}

func main() {
	// 1. Golang object -> JSON
	// ObjectToJSON()

	// 2. JSON -> Golang object
	// JSONToObject()

	// 3. maintance cost of map
	// MapToJSON()

	// 4. careful with repeat unmarshaling
	// First JSON string with the salary field
	jsonStr3 := `{"name":"Bob","position":"Developer","age":25,"salary":4000000}`
	var employee Employee
	_ = json.Unmarshal([]byte(jsonStr3), &employee)
	fmt.Println("employee3:", employee)
	// Second JSON string without the salary field
	jsonStr4 := `{"name":"Kate","position":"Senior Developer","age":35}`
	// Salary remains 4000000
	_ = json.Unmarshal([]byte(jsonStr4), &employee)
	fmt.Println("employee4:", employee)

	ZeroValueConfusion()

	TagMarshal()
}
