package main

import (
	"fmt"
	"net/http"
	// "strconv"

	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"github.com/labstack/echo/v4"
  
)

var (
	DB *gorm.DB
)

func init() {
	InitDB()
	InitialMigration()
}

type Config struct {
	DB_Username string
	DB_Password string
	DB_Port 	string
	DB_Host 	string
	DB_Name 	string
}

func InitDB() {
	config := Config{
		DB_Username: "root",
		DB_Password: "",
		DB_Port: "3306",
		DB_Host: "localhost",
		DB_Name: "crud_go",
	}

	connectionString :=
	fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_Name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

type Book struct {
	gorm.Model
	Tittle     string `json:"tittle" form:"tittle"`
  	Author    string `json:"author" form:"author"`

}

func InitialMigration() {
	DB.AutoMigrate(&Book{})
}

func GetBooksController(c echo.Context) error {
	var books []Book

	if err := DB.Find(&books).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message" : "success get all Book",
		"Book" : books,
	})
}

func GetBookController(c echo.Context) error {
	id := c.Param("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
	  return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
	  "message": "success get book",
	  "book":    book,
	})
   
  }
  
  
func CreateBookController(c echo.Context) error {
	book := Book{}
	c.Bind(&book)

	if err := DB.Save(&book).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create new book",
		"book":    book,
	})
}

func DeleteBookController(c echo.Context) error {
	id := c.Param("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
	  return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := DB.Delete(&book).Error; err != nil {
	  return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
	  "message": "success delete book",
	  "book":    book,
	})
  }
  
  func UpdateBookController(c echo.Context) error {
	id := c.Param("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
	  return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	c.Bind(&book)
	if err := DB.Save(&book).Error; err != nil {
	  return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
	  "message": "success update book",
	  "book":    book,
	})
  }
  

func main() {
	e := echo.New()
	e.GET("/book", GetBooksController)
    e.GET("/book/:id", GetBookController)
 	e.POST("/book", CreateBookController)
  	e.DELETE("/book/:id", DeleteBookController)
  	e.PUT("/book/:id", UpdateBookController)

	//start the server, and log if it fails
  	e.Logger.Fatal(e.Start(":8000"))

}

