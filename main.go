package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/public", "./public")

	client := r.Group("/api")
	{
		client.GET("/category/:id", Read)
		client.POST("/category/create", Create)
		client.PATCH("/category/update/:id", Update)
		client.DELETE("/category/:id", Delete)
	}

	return r
}

func main() {
	r := setupRouter()
	r.Run(":2301")
}

//Create database connection with config
func DBConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "12345678"
	dbName := "Sales"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

type Category struct {
	CategoryID   int    `json:"CategoryID"`
	CategoryName string `json:"CategoryName"`
	Description  string `json:"Description"`
}

func Read(c *gin.Context) {

	db := DBConn()
	rows, err := db.Query("SELECT CategoryID, CategoryName, Description FROM Categories WHERE CategoryID = " + c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"messages": "Category not found",
		})
	}

	category := Category{}

	for rows.Next() {
		var id int
		var name, description string

		err = rows.Scan(&id, &name, &description)
		if err != nil {
			panic(err.Error())
		}

		category.CategoryID = id
		category.CategoryName = name
		category.Description = description
	}

	c.JSON(200, category)
	defer db.Close()
	//Delay close database until Read() complete
}

//Create new category API
func Create(c *gin.Context) {
	db := DBConn()

	type CreateCategory struct {
		name        string `form:"name" json:"title" binding:"required"`
		description string `form:"description" json:"body" binding:"required"`
	}

	var json CreateCategory

	if err := c.ShouldBindJSON(&json); err == nil {
		insCategory, err := db.Prepare("INSERT INTO Categories(CategoryName, Description) VALUES(?,?)")
		if err != nil {
			c.JSON(500, gin.H{
				"messages": err,
			})
		}

		insCategory.Exec(json.name, json.description)
		c.JSON(200, gin.H{
			"messages": "new category inserted",
		})

	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	defer db.Close()
}

func Update(c *gin.Context) {
	db := DBConn()
	type UpdateCategory struct {
		Title string `form:"name" json:"title" binding:"required"`
		Body  string `form:"description" json:"body" binding:"required"`
	}

	var json UpdateCategory
	if err := c.ShouldBindJSON(&json); err == nil {
		edit, err := db.Prepare("UPDATE Category SET CategoryName=?, Description=? WHERE CategoryID= " + c.Param("id"))
		if err != nil {
			panic(err.Error())
		}
		edit.Exec(json.Title, json.Body)

		c.JSON(200, gin.H{
			"messages": "category was edited",
		})
	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	defer db.Close()
}

func Delete(c *gin.Context) {
	db := DBConn()

	delete, err := db.Prepare("DELETE FROM Category WHERE CategoryIDç=?")
	if err != nil {
		panic(err.Error())
	}

	delete.Exec(c.Param("id"))
	c.JSON(200, gin.H{
		"messages": "category was deleted",
	})

	defer db.Close()
}
