package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
}

func NewPerson() *Person {
	return &Person{}
}

func NewPeople() *[]Person {
	return &[]Person{}
}

func seedData() *[]Person {
	p1 := Person{FirstName: "firstname1", LastName: "lastname1"}
	p2 := Person{FirstName: "firstname2", LastName: "lastname2"}
	p3 := Person{FirstName: "firstname3", LastName: "lastname3"}
	p4 := Person{FirstName: "firstname4", LastName: "lastname4"}
	people := NewPeople()
	*people = append(*people, p1, p2, p3, p4)
	return people
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5435 user=postgres dbname=gorm password=root sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.AutoMigrate(NewPerson())

	// for _, person := range *seedData() {
	// 	fmt.Println(person)
	// 	db.Create(&person)
	// }

	r := gin.Default()

	getPeople := GetPeople(db)
	r.GET("/person", getPeople)

	getPerson := GetPerson(db)
	r.GET("/person/:id", getPerson)

	createPerson := CreatePerson(db)
	r.POST("/person", createPerson)

	updatePerson := UpdatePerson(db)
	r.PUT("/person/:id", updatePerson)

	deletePerson := DeletePerson(db)
	r.DELETE("/person/:id", deletePerson)

	log.Fatal(r.Run())
}

func GetPeople(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		people := NewPeople()
		if err := db.Find(people).Error; err != nil {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(200, people)
	}
}

func GetPerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		person := NewPerson()
		id, err := strconv.Atoi(c.Params.ByName("id"))
		if err != nil {
			c.AbortWithStatusJSON(500, err)
		}
		if err := db.Where("id=?", id).First(person).Error; err != nil {
			c.AbortWithStatus(404)
		}
		c.JSON(200, person)
	}
}

func CreatePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := NewPerson()
		c.BindJSON(p)
		fmt.Println(p)
		db.Create(p)
		c.JSON(http.StatusCreated, p)
	}
}

func UpdatePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Params.ByName("id"))
		if err != nil {
			c.AbortWithStatusJSON(500, err)
			return
		}

		p := NewPerson()
		if err = db.Where("id=?", id).First(p).Error; err != nil {
			c.AbortWithStatus(404)
			return
		}

		c.BindJSON(p)
		db.Save(p)

		c.JSON(http.StatusCreated, p)
	}
}

func DeletePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Params.ByName("id"))
		if err != nil {
			c.AbortWithStatusJSON(500, err)
		}

		p := NewPerson()
		if err = db.Where("id=?", id).Delete(p).Error; err != nil {
			c.AbortWithStatus(400)
		}
		c.JSON(204, gin.H{"message": fmt.Sprintf("person with id %d deleted", id)})
	}
}
