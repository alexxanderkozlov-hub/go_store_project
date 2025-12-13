package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginPage отображает страницу входа в систему
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// LoginHandler обрабатывает POST запрос на авторизацию
func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "admin" && password == "admin123" {
		c.SetCookie("auth", "true", 3600, "/", "", false, true)
		c.SetCookie("username", username, 3600, "/", "", false, false)
		// Перенаправление на главную страницу
		c.Redirect(http.StatusFound, "/main")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"error": "Неверный логин или пароль",
	})
}

func LogoutHandler(c *gin.Context) {
	c.SetCookie("auth", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func CheckAuth(c *gin.Context) bool {
	auth, err := c.Cookie("auth")
	return err == nil && auth == "true"
}

func GetUsername(c *gin.Context) string {
	username, err := c.Cookie("username")
	if err != nil {
		return "Гость"
	}
	return username
}

// MainMenuPage отображает главное меню приложения
func MainMenuPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "stores.html", gin.H{
		"username": GetUsername(c),
	})
}

// AboutPage отображает страницу "О проекте"
func AboutPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "about.html", gin.H{})
}

// ContactsPage отображает страницу контактов
func ContactsPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "contacts.html", gin.H{})
}
