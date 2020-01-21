package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type DB struct {
	*gorm.DB
}

type Client struct {
	ID       int
	Username string
	Password string
	MaxConn  string
}

const (
	tableName = "ListClient"
)

func (db *DB) GetClient(c *gin.Context) {
	client := Client{}
	result := gin.H{}

	id := c.Param("id")
	if err := db.DB.Table(tableName).Where("id = ?", id).First(&client).Error; err != nil {
		result = gin.H{"result": err}
	} else {
		result = gin.H{"result": client}
	}

	c.JSON(http.StatusOK, result)
}

//bit different output
func (db *DB) GetClients(c *gin.Context) {
	clients := []Client{}
	total := 0
	db.DB.Table(tableName).Count(&total)
	if total <= 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	if err := db.DB.Table(tableName).Find(&clients).Error; err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, clients)
}

func (db *DB) CreateClient(c *gin.Context) {
	client := Client{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		MaxConn:  c.PostForm("maxConn")}
	result := gin.H{}

	db.DB.Table(tableName).Create(&client)
	result = gin.H{
		"result": client,
	}
	c.JSON(http.StatusOK, result)
}

func (db *DB) UpdateClient(c *gin.Context) {
	client := Client{}
	result := gin.H{}

	id := c.Query("id")
	err := db.DB.Table(tableName).First(&client, id).Error
	if err != nil {
		result = gin.H{"result": "data not found"}
	}

	newClient := Client{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		MaxConn:  c.PostForm("maxConn")}
	err = db.DB.Table(tableName).Model(&client).Updates(newClient).Error
	if err != nil {
		result = gin.H{"result": "update failed"}
	} else {
		result = gin.H{"result": "successfully updated data"}
	}
	c.JSON(http.StatusOK, result)
}

func (db *DB) DeleteClient(c *gin.Context) {
	client := Client{}
	result := gin.H{}

	id := c.Param("id")
	if err := db.DB.Table(tableName).First(&client, id).Error; err != nil {
		result = gin.H{"result": "data not found"}
	}
	if err := db.DB.Table(tableName).Delete(&client).Error; err != nil {
		result = gin.H{"result": "delete failed"}
	} else {
		result = gin.H{"result": "Data deleted successfully"}
	}

	c.JSON(http.StatusOK, result)
}
