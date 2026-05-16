package handlers

import (
	"go_store_project/internal/models"
	"go_store_project/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterPage отображает страницу регистрации
func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

// RegisterHandler обрабатывает POST запрос регистрации
func RegisterHandler(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	// проверка совпадения паролей
	if password != confirmPassword {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "Пароли не совпадают",
		})
		return
	}

	// проверяем существует ли пользователь
	_, found := storage.GetUserByUsername(username)

	if found {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "Пользователь уже существует",
		})
		return
	}

	user := models.User{
		Username: username,
		Password: password,
		Role:     "customer",
	}

	id := storage.CreateUser(user)

	if id == 0 {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "Ошибка регистрации",
		})
		return
	}

	// после регистрации переходим на логин
	c.Redirect(http.StatusFound, "/login")
}
