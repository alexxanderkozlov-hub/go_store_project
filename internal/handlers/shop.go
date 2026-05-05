package handlers

import (
	"go_store_project/internal/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ================= SHOP =================
func ShopPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	products := storage.GetAllProducts()

	c.HTML(http.StatusOK, "shop.html", gin.H{
		"products": products,
		"role":     GetRole(c),
	})
}

// ================= ADD TO CART =================
func AddToCartHandler(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/shop")
		return
	}

	storage.AddToCart(userID, productID)

	c.Redirect(http.StatusFound, "/shop")
}

// ================= CART PAGE =================
func CartPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	cart := storage.GetCart(userID)
	total := storage.GetCartTotal(userID)

	c.HTML(http.StatusOK, "shop_form.html", gin.H{
		"cart":  cart,
		"total": total,
	})
}

// ================= REMOVE FROM CART =================
func RemoveFromCart(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/cart")
		return
	}

	storage.RemoveFromCart(userID, productID)

	c.Redirect(http.StatusFound, "/cart")
}

// ================= BUY FROM CART =================
func BuyFromCart(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/cart")
		return
	}

	// TODO: позже заменить на orders
	storage.RemoveFromCart(userID, productID)

	c.Redirect(http.StatusFound, "/cart")
}

// ================= USER ID =================
func getUserID(c *gin.Context) int {
	username, err := c.Cookie("username")
	if err != nil {
		return 0
	}

	user, found := storage.GetUserByUsername(username)
	if !found {
		return 0
	}

	return user.ID
}
