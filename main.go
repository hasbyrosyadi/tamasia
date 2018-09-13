package main

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type registerUser struct {
	ID int
	Nama string
}

var regis []registerUser

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/tamasia")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/regis", RegisterUser)
	r.GET("/regis/:nama", RegisterUserPilihan)
	r.Run()
}

func RegisterUser(c *gin.Context) {
	rows, err := db.Query("SELECT id, nama FROM registrasiUser")
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal server error",
		})
	}
	var listUser []registerUser
		for rows.Next(){
			var id int
			var nama string
			err = rows.Scan(&id, &nama)
			listUser = append(listUser, registerUser{
				ID:id, Nama:nama })
			fmt.Print(id, nama)
		}

		c.JSON(200, gin.H{
			"code": 200,
			"success" :true,
			"status": "Success",
			"registerUser": listUser,
		})
}

func RegisterUserPilihan(c *gin.Context) {
	name := c.Param("nama")
	rows, err := db.Query("SELECT id, nama FROM registrasiUser where nama = ?", name)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal Service Error",
		})
	}

	var listUser []registerUser
	for rows.Next(){
		var id int
		var nama string
		err = rows.Scan(&id, &nama)
		log.Println(err)
		listUser = append(listUser, registerUser{
			ID:id, Nama:nama })
	}

	if listUser == nil {
		c.JSON(404, gin.H{
			"code": 404,
			"success" :false,
			"status": "Failed",
			"message": "data tidak ditemukan",
		})
	}else{
		c.JSON(200, gin.H{
			"code": 200,
			"success" :true,
			"status": "Success",
			"registerUser": listUser,
		})
	}
}