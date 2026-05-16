package handlers

import (
	"go_store_project/internal/models"
	"go_store_project/internal/storage"
	"net/http"
	"sort" // Добавь для сортировки
	"strconv"
	"strings" // Добавь для поиска без учёта регистра

	"github.com/gin-gonic/gin"
)

// ================= SHOP =================
func ShopPage(c *gin.Context) {
	// (Твоя проверка авторизации, оставь как была в оригинале)
	if !CheckAuth(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 1. Получаем параметры фильтрации из URL
	searchQuery := c.Query("search_query")
	sortBy := c.Query("sort_by")
	direction := c.Query("direction")

	// Получаем все исходные товары из БД (они приходят с типом models.Product)
	allProducts := storage.GetAllProducts()

	// Инициализируем слайс правильным типом из пакета models
	var filteredProducts []models.Product

	// 2. ФИЛЬТРАЦИЯ (Поиск по названию, описанию или артикулу)
	if searchQuery != "" {
		searchQueryLower := strings.ToLower(searchQuery)
		for _, p := range allProducts {
			if strings.Contains(strings.ToLower(p.Name), searchQueryLower) ||
				strings.Contains(strings.ToLower(p.Description), searchQueryLower) ||
				strings.Contains(strings.ToLower(p.SKU), searchQueryLower) {
				filteredProducts = append(filteredProducts, p)
			}
		}
	} else {
		// Если поиска нет, работаем со всеми товарами
		filteredProducts = allProducts
	}

	// 3. СОРТИРОВКА
	if sortBy != "" {
		sort.Slice(filteredProducts, func(i, j int) bool {
			var isLess bool
			switch sortBy {
			case "id":
				isLess = filteredProducts[i].ID < filteredProducts[j].ID
			case "name":
				isLess = filteredProducts[i].Name < filteredProducts[j].Name
			case "price":
				// В Go для DECIMAL (из базы) обычно используется float64 в структурах
				isLess = filteredProducts[i].Price < filteredProducts[j].Price
			case "sku":
				isLess = filteredProducts[i].SKU < filteredProducts[j].SKU
			default:
				isLess = filteredProducts[i].ID < filteredProducts[j].ID
			}

			if direction == "desc" {
				return !isLess
			}
			return isLess
		})
	}

	// 4. Передаем данные в HTML шаблонизатор
	c.HTML(http.StatusOK, "shop.html", gin.H{
		"products":     filteredProducts,
		"role":         GetRole(c),
		"search_query": searchQuery,
		"sort_by":      sortBy,
		"direction":    direction,
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
