package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"

	"gorm.io/driver/sqlite"

	"strconv"
)

type Order struct {
	gorm.Model
	ID       uint
	FullName string
}

func StringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Order{})

	e := echo.New()
	e.GET("/orders", func(c echo.Context) error {
		var orders []Order
		db.Find(&orders)

		return c.JSON(http.StatusOK, orders)
	})

	e.GET("/orders/:id", func(c echo.Context) error {
		var order Order

		db.Last(&order, c.Param("id"))

		return c.JSON(http.StatusOK, order)
	})

	e.POST("/orders", func(c echo.Context) error {
		fullName := c.FormValue("fullName")
		id := StringToUint(c.FormValue("id"))

		order := Order{FullName: fullName, ID: id}

		db.Create(&order)

		return c.String(http.StatusCreated, "Order created")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
