package handlers

import (
	"go_store_project/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// LoginHandler обрабатывает POST запрос на авторизацию
func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// получаем пользователя из БД
	user, found := storage.GetUserByUsername(username)
	if !found {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "Пользователь не найден",
		})
		return
	}

	// проверяем пароль
	if user.Password != password {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "Неверный пароль",
		})
		return
	}

	// авторизация
	c.SetCookie("auth", "true", 3600, "/", "", false, true)
	c.SetCookie("username", user.Username, 3600, "/", "", false, false)
	c.SetCookie("role", user.Role, 3600, "/", "", false, false)

	// редирект по роли
	if user.Role == "admin" {
		c.Redirect(http.StatusFound, "/main")
	} else {
		c.Redirect(http.StatusFound, "/shop")
	}
}

// LogoutHandler обрабатывает выход пользователя
func LogoutHandler(c *gin.Context) {
	c.SetCookie("auth", "", -1, "/", "", false, true)
	c.SetCookie("username", "", -1, "/", "", false, false)
	c.SetCookie("role", "", -1, "/", "", false, false)
	c.Redirect(http.StatusFound, "/login")
}

// CheckAuth проверяет авторизацию
func CheckAuth(c *gin.Context) bool {
	auth, err := c.Cookie("auth")
	return err == nil && auth == "true"
}

// GetUsername возвращает имя пользователя
func GetUsername(c *gin.Context) string {
	username, err := c.Cookie("username")
	if err != nil {
		return "Гость"
	}
	return username
}

// GetRole возвращает роль пользователя
func GetRole(c *gin.Context) string {
	role, err := c.Cookie("role")
	if err != nil {
		return "guest"
	}
	return role
}

// MainMenuPage отображает главное меню (только для админа)
func MainMenuPage(c *gin.Context) {
	if !CheckAuth(c) || GetRole(c) != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "stores.html", gin.H{
		"username": GetUsername(c),
	})
}

// ShopPage отображает страницу магазина (для покупателей)
func ShopPage(c *gin.Context) {
	if !CheckAuth(c) || GetRole(c) != "customer" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "shop.html", gin.H{
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
