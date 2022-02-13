package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	r := setupRouter()
	r.Run()
}

type List struct {
	Id    uint   `gorm:"primary_key" json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Member struct {
	Id       uint   `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	h := CoffeeHandler{}
	h.Initialize()

	r.GET("/list", h.GetList)
	r.POST("/list/add", h.AddList)
	r.POST("/member/login", h.MemberLogin)
	return r
}

type CoffeeHandler struct {
	DB *gorm.DB
}

func (h *CoffeeHandler) GetList(c *gin.Context) {
	list := []List{}

	h.DB.Find(&list)
	fmt.Print(list)
	c.JSON(http.StatusOK, list)
}

func (h *CoffeeHandler) AddList(c *gin.Context) {
	list := List{}

	if err := c.ShouldBindJSON(&list); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DB.Save(&list).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *CoffeeHandler) MemberLogin(c *gin.Context) {
	mem := Member{}
	var member Member
	c.Bind(&member)

	if err := h.DB.Where("(username = ? OR email = ?) AND password = ? AND role = ?", member.Username, member.Username, member.Password, member.Role).Find(&mem).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, mem)
}

func (h *CoffeeHandler) Initialize() {
	sqlDB, _ := sql.Open("mysql", "root:root@tcp(localhost:8889)/coffee")
	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	gormDB.AutoMigrate(&List{})
	gormDB.AutoMigrate(&Member{})

	h.DB = gormDB
}
