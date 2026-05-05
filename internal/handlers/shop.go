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

	productID, _ := strconv.Atoi(c.Param("id"))

	// получаем корзину
	cart := storage.GetCart(userID)

	for _, item := range cart {
		if item.ProductID == productID {
			storage.CreateOrder(userID, productID, item.Quantity, item.Total)
			break
		}
	}

	// удаляем из корзины
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

func OrdersPage(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	orders := storage.GetOrders()

	c.HTML(http.StatusOK, "orders.html", gin.H{
		"orders": orders,
	})
}

func UpdateOrderStatusHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	status := c.Param("status")

	storage.UpdateOrderStatus(id, status)

	c.Redirect(http.StatusFound, "/orders")
}

func DeleteOrder(c *gin.Context) {
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/orders")
		return
	}

	storage.DeleteOrder(id)

	c.Redirect(http.StatusFound, "/orders")
}
