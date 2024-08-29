package controller

import (
	"UMKM/app/store"
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

// Definisikan struct MenuItem
type MenuItem struct {
	ID       int
	ItemName string
	Price    float64
}

// Handler untuk membuat pesanan baru
func NewAddOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Store.Get(r, "session-name")

		userID, userIDok := session.Values["user_id"].(int)
		username, usernameok := session.Values["username"].(string)
		if !userIDok || !usernameok {
			http.Redirect(w, r, "/verfiedUserLogin", http.StatusSeeOther)
			return
		}

		if r.Method == http.MethodGet {
			// Ambil item dari tabel menu
			rows, err := db.Query("SELECT idMenu, item_name, harga FROM menu")
			if err != nil {
				http.Error(w, "Error fetching menu items: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var menuItems []MenuItem
			for rows.Next() {
				var item MenuItem
				if err := rows.Scan(&item.ID, &item.ItemName, &item.Price); err != nil {
					http.Error(w, "Error scanning menu items: "+err.Error(), http.StatusInternalServerError)
					return
				}
				menuItems = append(menuItems, item)
			}

			data := struct {
				UserData  UserData
				MenuItems []MenuItem
			}{
				UserData: UserData{
					UserID:   userID,
					Username: username,
				},
				MenuItems: menuItems,
			}

			fp := filepath.Join("app", "views", "order.html")
			tmpl, err := template.ParseFiles(fp)
			if err != nil {
				http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		if r.Method == http.MethodPost {
			orderDate := r.FormValue("order_date")
			shippingAddress := r.FormValue("shipping_address")
			paymentMethod := r.FormValue("payment_method")
			menuIDs := r.Form["menu_item"]

			// Hitung total amount
			totalAmount := 0.0

			for _, menuID := range menuIDs {
				quantityStr := r.FormValue("quantity_" + menuID)
				quantity, err := strconv.Atoi(quantityStr)
				if err != nil {
					http.Error(w, "Invalid quantity provided", http.StatusBadRequest)
					return
				}

				var price float64
				err = db.QueryRow("SELECT harga FROM menu WHERE idMenu = ?", menuID).Scan(&price)
				if err != nil {
					http.Error(w, "Error fetching menu price: "+err.Error(), http.StatusInternalServerError)
					return
				}

				totalAmount += price * float64(quantity)
			}

			// Simpan pesanan
			query := `INSERT INTO pesanan (id_User, total_amount, order_date, shipping_address, payment_method, payment_status, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?, 'Pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

			result, err := db.Exec(query, userID, totalAmount, orderDate, shippingAddress, paymentMethod)
			if err != nil {
				http.Error(w, "Error inserting order: "+err.Error(), http.StatusInternalServerError)
				return
			}

			orderID, _ := result.LastInsertId()

			// Simpan detail pesanan untuk setiap item yang dipilih
			for _, menuID := range menuIDs {
				quantity := r.FormValue("quantity_" + menuID)
				_, err := db.Exec("INSERT INTO orderDetails (order_id, menu_id, quantity) VALUES (?, ?, ?)", orderID, menuID, quantity)
				if err != nil {
					http.Error(w, "Error inserting order items: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
