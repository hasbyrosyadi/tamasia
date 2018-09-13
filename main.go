package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

type registerUser struct {
	ID int
	Nama string
}

type loginUser struct {
	Id int
	IdRegister int
	Nama string
}

type orderProduk struct {
	id int
	idLogin int
	nama string
	harga int
}

var regis []registerUser
var login []loginUser
var order []orderProduk

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/tamasia")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/regis/", RegisterUser)
	r.GET("/regis/:nama", RegisterUserPilihan)
	r.GET("/login/viewall", LoginUser)
	r.GET("/login", LoginUserPilihan)
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

func LoginUser(c *gin.Context) {
	rows, err := db.Query("SELECT id, idRegister, nama FROM loginUser")
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal server error",
		})
	}
	var listUser []loginUser
	for rows.Next(){
		var id int
		var idregis int
		var nama string
		err = rows.Scan(&id, &idregis, &nama)
		listUser = append(listUser, loginUser{
			Id:id, IdRegister:idregis, Nama:nama })
	}

	c.JSON(200, gin.H{
		"code": 200,
		"success" :true,
		"status": "Success",
		"loginUser": listUser,
	})
}

func LoginUserPilihan(c *gin.Context) {
	name := c.Query("nama")
	idRegis := c.Query("idregis")
	rows, err := db.Query("SELECT id, idRegister, nama FROM loginUser where nama = ? and idRegister = ?", name, idRegis)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal Service Error",
		})
	}

	var listUser []loginUser
	for rows.Next(){
		var id int
		var idregis int
		var nama string
		err = rows.Scan(&id, &idregis, &nama)
		log.Println(err)
		listUser = append(listUser, loginUser{
			Id:id, IdRegister:idregis, Nama:nama })
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
			"loginUser": listUser,
		})
	}
}

