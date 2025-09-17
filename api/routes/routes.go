package routes

import (
	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/controllers"
)

func RegisterRoutes(router *mux.Router) {
	// Маршруты для продуктов
	router.HandleFunc("/api/products/all", controllers.GetProducts).Methods("GET")
	router.HandleFunc("/api/products/get/{id}", controllers.GetProduct).Methods("GET")
	router.HandleFunc("/api/products/add", controllers.CreateProduct).Methods("POST")
	router.HandleFunc("/api/products/update/{id}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/api/products/delete/{id}", controllers.DeleteProduct).Methods("DELETE")

	// Маршруты для блюд
	router.HandleFunc("/api/dishes/all", controllers.GetDishes).Methods("GET")
	router.HandleFunc("/api/dishes/get/{id}", controllers.GetDish).Methods("GET")
	router.HandleFunc("/api/dishes/add", controllers.CreateDish).Methods("POST")
	router.HandleFunc("/api/dishes/update/{id}", controllers.UpdateDish).Methods("PUT")
	router.HandleFunc("/api/dishes/delete/{id}", controllers.DeleteDish).Methods("DELETE")

	// Маршруты для ингредиентов блюд
	router.HandleFunc("/api/dish_ingredients/all", controllers.GetDishIngredients).Methods("GET")
	router.HandleFunc("/api/dish_ingredients/get/{id}", controllers.GetDishIngredient).Methods("GET")
	router.HandleFunc("/api/dish_ingredients/add", controllers.CreateDishIngredient).Methods("POST")
	router.HandleFunc("/api/dish_ingredients/update/{id}", controllers.UpdateDishIngredient).Methods("PUT")
	router.HandleFunc("/api/dish_ingredients/delete/{id}", controllers.DeleteDishIngredient).Methods("DELETE")

	// Маршруты для меню
	router.HandleFunc("/api/menu/all", controllers.GetMenu).Methods("GET")
	router.HandleFunc("/api/menu/get/{id}", controllers.GetMenuItem).Methods("GET")
	router.HandleFunc("/api/menu/add", controllers.CreateMenuItem).Methods("POST")
	router.HandleFunc("/api/menu/update/{id}", controllers.UpdateMenuItem).Methods("PUT")
	router.HandleFunc("/api/menu/delete/{id}", controllers.DeleteMenuItem).Methods("DELETE")

	// Маршруты для столиков
	router.HandleFunc("/api/tables/all", controllers.GetTables).Methods("GET")
	router.HandleFunc("/api/tables/get/{id}", controllers.GetTable).Methods("GET")
	router.HandleFunc("/api/tables/add", controllers.CreateTable).Methods("POST")
	router.HandleFunc("/api/tables/update/{id}", controllers.UpdateTable).Methods("PUT")
	router.HandleFunc("/api/tables/delete/{id}", controllers.DeleteTable).Methods("DELETE")

	// Маршруты для заказов
	router.HandleFunc("/api/orders/all", controllers.GetOrders).Methods("GET")
	router.HandleFunc("/api/orders/get/{id}", controllers.GetOrder).Methods("GET")
	router.HandleFunc("/api/orders/add", controllers.CreateOrder).Methods("POST")
	router.HandleFunc("/api/orders/update/{id}", controllers.UpdateOrder).Methods("PUT")
	router.HandleFunc("/api/orders/delete/{id}", controllers.DeleteOrder).Methods("DELETE")

	// Маршруты для пунктов заказа
	router.HandleFunc("/api/order_items/all", controllers.GetOrderItems).Methods("GET")
	router.HandleFunc("/api/order_items/get/{id}", controllers.GetOrderItem).Methods("GET")
	router.HandleFunc("/api/orders/items/{order_id}", controllers.AddDishToOrder).Methods("POST")
	router.HandleFunc("/api/order_items/update/{id}", controllers.UpdateOrderItem).Methods("PUT")
	router.HandleFunc("/api/order_items/delete/{id}", controllers.DeleteOrderItem).Methods("DELETE")
}
