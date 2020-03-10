package handler

import (
	"fmt"
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
	Nama     string
	Foto     string
}

const (
	TABLENAME = "user"
)

func (db *DB) GetUser(c *gin.Context) {
	client := Client{}
	result := gin.H{}

	id := c.Param("id")
	if err := db.DB.Table(TABLENAME).Where("id = ?", id).First(&client).Error; err != nil {
		result = gin.H{"result": err}
	} else {
		result = gin.H{"result": client}
	}

	c.JSON(http.StatusOK, result)
}

//bit different output
func (db *DB) GetUsers(c *gin.Context) {
	clients := []Client{}
	total := 0
	db.DB.Table(TABLENAME).Count(&total)
	if total <= 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	if err := db.DB.Table(TABLENAME).Find(&clients).Error; err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, clients)
}

func (db *DB) CreateUser(c *gin.Context) {
	file, err := c.FormFile("foto")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}
	path := "foto/" + file.Filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	client := Client{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Nama:     c.PostForm("nama"),
		Foto:     path}
	result := gin.H{}
	db.DB.Table(TABLENAME).Create(&client)
	result = gin.H{
		"result": client,
	}
	c.JSON(http.StatusOK, result)
}

func (db *DB) UpdateUser(c *gin.Context) {
	client := Client{}
	result := gin.H{}

	id := c.Query("id")
	err := db.DB.Table(TABLENAME).First(&client, id).Error
	if err != nil {
		result = gin.H{"result": "data not found"}
	}

	file, err := c.FormFile("foto")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}
	path := "foto/" + file.Filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	newClient := Client{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Nama:     c.PostForm("nama"),
		Foto:     path}
	err = db.DB.Table(TABLENAME).Model(&client).Updates(newClient).Error
	if err != nil {
		result = gin.H{"result": "update failed"}
	} else {
		result = gin.H{"result": "successfully updated data"}
	}
	c.JSON(http.StatusOK, result)
}

func (db *DB) DeleteUser(c *gin.Context) {
	client := Client{}
	result := gin.H{}

	id := c.Param("id")
	if err := db.DB.Table(TABLENAME).First(&client, id).Error; err != nil {
		result = gin.H{"result": "data not found"}
	}
	if err := db.DB.Table(TABLENAME).Delete(&client).Error; err != nil {
		result = gin.H{"result": "delete failed"}
	} else {
		result = gin.H{"result": "Data deleted successfully"}
	}

	c.JSON(http.StatusOK, result)
}
