package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB

type registerUser struct {
	Id int
	FullName string
}

type loginUser struct {
	Id int
	Nama string
	Password string
}

type orderProduk struct {
	Id int
	IdLogin int
	Nama string
	Harga int
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
	r.GET("/login/", LoginUser)
	r.GET("/login/:nama", LoginUserPilihan)
	//r.Use(logger())
	r.GET("/order/", OrderProduk)
	r.GET("/order/:nama", OrderProdukPilihan)

	r.Run()
}

func RegisterUser(c *gin.Context) {
	rows, err := db.Query("SELECT id, fullname FROM users")
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
				Id:id, FullName:nama })
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
	rows, err := db.Query("SELECT id, fullname FROM users where name = ?", name)
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
			Id:id, FullName:nama })
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
	rows, err := db.Query("SELECT id, name, password FROM users")
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
		var nama string
		var password string
		err = rows.Scan(&id, &nama, &password)
		listUser = append(listUser, loginUser{
			Id:id, Password:password, Nama:nama })
	}

	c.JSON(200, gin.H{
		"code": 200,
		"success" :true,
		"status": "Success",
		"loginUser": listUser,
	})
}

func LoginUserPilihan(c *gin.Context) {
	name := c.Param("nama")
	rows, err := db.Query("SELECT id, name, password FROM users where name = ?", name)
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
		var nama string
		var password string
		err = rows.Scan(&id, &nama, &password)
		listUser = append(listUser, loginUser{
			Id:id, Password:password, Nama:nama })
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

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		aout := c.GetHeader("Authorization")
		log.Println(aout)
		if aout != "tamasia" {
			c.JSON(400, gin.H{
				"code": 400,
				"success" :false,
				"status": "Failed",
				"message": "tidak ada akses",
			})
			c.Abort()
		}

	}
}

func OrderProduk(c *gin.Context) {
	rows, err := db.Query("SELECT id, idLogin, nama, harga FROM orderProduk")
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal server error",
		})
	}
	var listOrder []orderProduk
	for rows.Next(){
		var id int
		var idlogin int
		var nama string
		var harga int
		err = rows.Scan(&id, &idlogin, &nama, &harga)
		listOrder = append(listOrder, orderProduk{
			Id:id, IdLogin:idlogin, Nama:nama, Harga:harga })
	}

	c.JSON(200, gin.H{
		"code": 200,
		"success" :true,
		"status": "Success",
		"OrderProduk": listOrder,
	})
}

func OrderProdukPilihan(c *gin.Context) {
	name := c.Param("nama")
	rows, err := db.Query("SELECT id, idLogin, nama, harga FROM orderProduk where nama = ?", name)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal Service Error",
		})
	}

	var listOrder []orderProduk
	for rows.Next(){
		var id int
		var idlogin int
		var nama string
		var harga int
		err = rows.Scan(&id, &idlogin, &nama, &harga)
		listOrder = append(listOrder, orderProduk{
			Id:id, IdLogin:idlogin, Nama:nama, Harga:harga })
	}

	if listOrder == nil {
		c.JSON(404, gin.H{
			"code": 404,
			"success" :false,
			"status": "Failed",
			"message": "data tidak ditemukan",
		})
	}else{
		c.JSON(200, gin.H{
			"Code": 200,
			"Success" :true,
			"Status": "Success",
			"OrderProduk": listOrder,
		})
	}
}

//func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
//	var user users
//	_ = json.NewDecoder(req.Body).Decode(&user)
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"username": user.Username,
//		"password": user.Password,
//	})
//	tokenString, error := token.SignedString([]byte("secret"))
//	if error != nil {
//		fmt.Println(error)
//	}
//	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
//}