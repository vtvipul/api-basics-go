package main

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func New() *Person {
	return &Person{}
}

func main() {
	db, _ := gorm.Open("postgres", "host=localhost port=5435 user=postgres dbname=gorm password=root sslmode=disable")
	defer db.Close()
	p1 := Person{FirstName: "firstname1", LastName: "lastname1"}
	p2 := Person{FirstName: "firstname2", LastName: "lastname2"}

	db.AutoMigrate(New())

	db.Create(&p1)
	db.Create(&p2)

	var p3 = New()
	db.First(p3)
	fmt.Println(p3)
	j := json.NewEncoder(os.Stdout).Encode(p3)
	fmt.Println(j)
	fmt.Println("-------------------------------------------------------------")
	type persons []Person
	p := db.Find(&persons{}).Value
	j1 := json.NewEncoder(os.Stdout).Encode(p)
	fmt.Println(j1)
}
