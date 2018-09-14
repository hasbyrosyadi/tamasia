package main

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
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
	IdUsers int
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

	r.GET("/regis", RegisterUser)
	r.GET("/regis/:nama", RegisterUserPilihan)
	r.GET("/login", LoginUser)
	r.GET("/login/:nama", LoginUserPilihan)
	r.GET("/order/:nama", CreateTokenEndpoint)
	r.Use(logger())
	r.GET("/order", OrderProduk)
	r.Run()
}

/* code dibawah ini merupakan proses pembuatan API yang mengembalikan semua users yang melakukan register.
	data yang ditampilakan sebatas id user dan nama lengkap dari user tersebut */

func RegisterUser(c *gin.Context) {
	rows, err := db.Query("SELECT id, fullname FROM users")
	if err != nil {
		// mengembalikan error 500 apabila tidak tersambung dengan database
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal server error",
		})
	}
	var listUser []registerUser
	// memasukan data di database kedalam struct yang dibuat
		for rows.Next(){
			var id int
			var nama string
			err = rows.Scan(&id, &nama)
			listUser = append(listUser, registerUser{
				Id:id, FullName:nama })
		}
	// mengembalikan code 200 apabila data yang dicari ditemukan
		c.JSON(200, gin.H{
			"code": 200,
			"success" :true,
			"status": "Success",
			"registerUser": listUser,
		})
}

/* code dibawah ini merupakan proses pembuatan API yang mengembalikan users yang melakukan register dengan parameter nama.
	data yang ditampilakan sebatas id user dan nama lengkap dari user tersebut */

func RegisterUserPilihan(c *gin.Context) {
	name := c.Param("nama")
	rows, err := db.Query("SELECT id, fullname FROM users where name = ?", name)
	if err != nil {
		// mengembalikan error 500 apabila tidak tersambung dengan database
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal Service Error",
		})
	}

	var listUser []registerUser
	// memasukan data di database kedalam struct yang dibuat
	for rows.Next(){
		var id int
		var nama string
		err = rows.Scan(&id, &nama)
		log.Println(err)
		listUser = append(listUser, registerUser{
			Id:id, FullName:nama })
	}

	if listUser == nil {
		// mengembalikan error 404 apabila data yang dicari tidak ditemukan
		c.JSON(404, gin.H{
			"code": 404,
			"success" :false,
			"status": "Failed",
			"message": "data tidak ditemukan",
		})
	}else{
		// mengembalikan code 200 apabila data yang dicari ditemukan
		c.JSON(200, gin.H{
			"code": 200,
			"success" :true,
			"status": "Success",
			"registerUser": listUser,
		})
	}
}

/* code dibawah ini merupakan proses pembuatan API yang mengembalikan semua users yang melakukan login.
	data yang ditampilakan sebatas id user, nama dari user tersebut */

func LoginUser(c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		// mengembalikan error 500 apabila tidak tersambung dengan database
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal server error",
		})
	}
	var listUser []loginUser
	// memasukan data di database kedalam struct yang dibuat
	for rows.Next(){
		var id int
		var nama string
		err = rows.Scan(&id, &nama)
		listUser = append(listUser, loginUser{
			Id:id, Nama:nama })
	}
	// mengembalikan code 200 apabila data yang dicari ditemukan
	c.JSON(200, gin.H{
		"code": 200,
		"success" :true,
		"status": "Success",
		"loginUser": listUser,
	})
}

/* code dibawah ini merupakan proses pembuatan API yang mengembalikan users yang melakukan login dengan parameter nama.
	data yang ditampilakan sebatas id user, nama dari user tersebut */

func LoginUserPilihan(c *gin.Context) {
	name := c.Param("nama")
	rows, err := db.Query("SELECT id, name FROM users where name = ?", name)
	if err != nil {
		// mengembalikan error 500 apabila tidak tersambung dengan database
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal Service Error",
		})
	}

	var listUser []loginUser
	// memasukan data di database kedalam struct yang dibuat
	for rows.Next(){
		var id int
		var nama string
		err = rows.Scan(&id, &nama)
		listUser = append(listUser, loginUser{
			Id:id, Nama:nama })
	}

	if listUser == nil {
		// mengembalikan error 404 apabila data yang dicari tidak ditemukan
		c.JSON(404, gin.H{
			"code": 404,
			"success" :false,
			"status": "Failed",
			"message": "data tidak ditemukan",
		})
	}else{
		// mengembalikan code 200 apabila data yang dicari ditemukan
		c.JSON(200, gin.H{
			"code": 200,
			"success" :true,
			"status": "Success",
			"loginUser": listUser,
		})
	}
}
/* code dibawah ini merupakan proses pembuatan Authorization yang token users yang melakukan login dengan parameter nama.
	pass decode yang yang digunakan adalah secret */
func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		aout := c.GetHeader("Authorization")
		if aout == "" {
			// mengembalikan error 404 apabila data yang dicari tidak ditemukan
			c.JSON(404, gin.H{
				"code": 404,
				"success" :false,
				"status": "Failed",
				"message": "tidak ada akses",
			})
			c.Abort()
		}
		// proses penggunaan token yang telah dibuat di order/:nama
		token, _ := jwt.Parse(aout, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte("secret"), nil
		})

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		} else {
			// mengembalikan error 404 apabila data yang dicari tidak ditemukan
			c.JSON(404, gin.H{
				"code": 404,
				"success" :false,
				"status": "Failed",
				"message": "tidak ada akses",
			})
			c.Abort()
		}
	}
}

/* code dibawah ini merupakan proses pembuatan API yang mengembalikan semua order produk.
	data yang ditampilakan sebatas id user, nama dari user tersebut */

func OrderProduk(c *gin.Context) {
	rows, err := db.Query("SELECT id, idUsers, nama, harga FROM orderProduk")
	if err != nil {
		// mengembalikan error 500 apabila tidak tersambung dengan database
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal server error",
		})
	}
	var listOrder []orderProduk
	// memasukan data di database kedalam struct yang dibuat
	for rows.Next(){
		var id int
		var idusers int
		var nama string
		var harga int
		err = rows.Scan(&id, &idusers, &nama, &harga)
		listOrder = append(listOrder, orderProduk{
			Id:id, IdUsers:idusers, Nama:nama, Harga:harga })
	}
	// mengembalikan code 200 apabila data yang dicari ditemukan
	c.JSON(200, gin.H{
		"code": 200,
		"success" :true,
		"status": "Success",
		"OrderProduk": listOrder,
	})
}

/* code dibawah ini merupakan proses pembuatan Token yang mengembalikan hasil decode dari user yang login dengan
	parameter nama.
	data yang ditampilkan sebatas id user, nama, dan password dari user tersebut */

func CreateTokenEndpoint(c *gin.Context) {
	name := c.Param("nama")
	rows, err := db.Query("SELECT id, name, password FROM users where name = ?", name)
	if err != nil {
		// mengembalikan error 500 apabila tidak tersambung dengan database
		c.JSON(500, gin.H{
			"code": 500,
			"message": "Internal Service Error",
		})
	}

	var id int
	var nama string
	var password string
	// memasukan data di database kedalam struct yang dibuat
	for rows.Next(){
		err = rows.Scan(&id, &nama, &password)
	}

	// pengambilan id, nama, dan password dari db dan dimasukan ke variabel user
	user := loginUser{Id:id, Nama:nama, Password:password}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Nama,
		"password": user.Password,
	})
	// password decode yang digunakan adalah secret
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	// mengembalikan pesan sukses apabila berhasil
	c.JSON(200, tokenString)
}