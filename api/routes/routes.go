package routes

import (
	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/controllers"
)

func RegisterRoutes(router *mux.Router) {
	// Маршруты для продуктов
	router.HandleFunc("/api/products", controllers.GetProducts).Methods("GET")
	router.HandleFunc("/api/products/{id}", controllers.GetProduct).Methods("GET")
	router.HandleFunc("/api/products", controllers.CreateProduct).Methods("POST")
	router.HandleFunc("/api/products/{id}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/api/products/{id}", controllers.DeleteProduct).Methods("DELETE")

	// Маршруты для блюд
	router.HandleFunc("/api/dishes", controllers.GetDishes).Methods("GET")
	router.HandleFunc("/api/dishes/{id}", controllers.GetDish).Methods("GET")
	router.HandleFunc("/api/dishes", controllers.CreateDish).Methods("POST")
	router.HandleFunc("/api/dishes/{id}", controllers.UpdateDish).Methods("PUT")
	router.HandleFunc("/api/dishes/{id}", controllers.DeleteDish).Methods("DELETE")

	// Маршруты для ингредиентов блюд
	router.HandleFunc("/api/dish_ingredients", controllers.GetDishIngredients).Methods("GET")
	router.HandleFunc("/api/dish_ingredients/{id}", controllers.GetDishIngredient).Methods("GET")
	router.HandleFunc("/api/dish_ingredients", controllers.CreateDishIngredient).Methods("POST")
	router.HandleFunc("/api/dish_ingredients/{id}", controllers.UpdateDishIngredient).Methods("PUT")
	router.HandleFunc("/api/dish_ingredients/{id}", controllers.DeleteDishIngredient).Methods("DELETE")

	// Маршруты для меню
	router.HandleFunc("/api/menu", controllers.GetMenu).Methods("GET")
	router.HandleFunc("/api/menu/{id}", controllers.GetMenuItem).Methods("GET")
	router.HandleFunc("/api/menu", controllers.CreateMenuItem).Methods("POST")
	router.HandleFunc("/api/menu/{id}", controllers.UpdateMenuItem).Methods("PUT")
	router.HandleFunc("/api/menu/{id}", controllers.DeleteMenuItem).Methods("DELETE")

	// Маршруты для столиков
	router.HandleFunc("/api/tables", controllers.GetTables).Methods("GET")
	router.HandleFunc("/api/tables/{id}", controllers.GetTable).Methods("GET")
	router.HandleFunc("/api/tables", controllers.CreateTable).Methods("POST")
	router.HandleFunc("/api/tables/{id}", controllers.UpdateTable).Methods("PUT")
	router.HandleFunc("/api/tables/{id}", controllers.DeleteTable).Methods("DELETE")

	// Маршруты для заказов
	router.HandleFunc("/api/orders", controllers.GetOrders).Methods("GET")
	router.HandleFunc("/api/orders/{id}", controllers.GetOrder).Methods("GET")
	router.HandleFunc("/api/orders", controllers.CreateOrder).Methods("POST")
	router.HandleFunc("/api/orders/{id}", controllers.UpdateOrder).Methods("PUT")
	router.HandleFunc("/api/orders/{id}", controllers.DeleteOrder).Methods("DELETE")

	// Маршруты для пунктов заказа
	router.HandleFunc("/api/order_items", controllers.GetOrderItems).Methods("GET")
	router.HandleFunc("/api/order_items/{id}", controllers.GetOrderItem).Methods("GET")
	router.HandleFunc("/api/orders/{order_id}/items", controllers.AddDishToOrder).Methods("POST")
	router.HandleFunc("/api/order_items/{id}", controllers.UpdateOrderItem).Methods("PUT")
	router.HandleFunc("/api/order_items/{id}", controllers.DeleteOrderItem).Methods("DELETE")
}
