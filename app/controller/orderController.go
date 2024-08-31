package controller

import (
	"UMKM/app/store"
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

type MenuItem struct {
	ID       int
	ItemName string
	Price    float64
}

type OrderItem struct {
	ItemName string
	Quantity int
	Price    float64
}

func NewAddOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ambil session pengguna
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

			// Render template order.html
			fp := filepath.Join("app", "views", "order.html")
			tmpl, err := template.ParseFiles(fp)
			if err != nil {
				http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if err = tmpl.Execute(w, data); err != nil {
				http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		if r.Method == http.MethodPost {
			// Ambil data dari form
			orderDate := r.FormValue("order_date")
			shippingAddress := r.FormValue("shipping_address")
			paymentMethod := r.FormValue("payment_method")
			NoTelp := r.FormValue("NoTelp")
			pesanTambahan := r.FormValue("pesanTambahan")
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

			// Simpan pesanan menggunakan ON DUPLICATE KEY UPDATE
			query := `INSERT INTO pesanan (id_User, total_amount, order_date, shipping_address, payment_method, NoTelp, pesanTambahan, payment_status, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?, 'Pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
				ON DUPLICATE KEY UPDATE total_amount = VALUES(total_amount), order_date = VALUES(order_date), 
				shipping_address = VALUES(shipping_address), payment_method = VALUES(payment_method), 
				NoTelp = VALUES(NoTelp), pesanTambahan = VALUES(pesanTambahan), updated_at = CURRENT_TIMESTAMP`

			result, err := db.Exec(query, userID, totalAmount, orderDate, shippingAddress, paymentMethod, NoTelp, pesanTambahan)
			if err != nil {
				http.Error(w, "Error inserting order: "+err.Error(), http.StatusInternalServerError)
				return
			}

			orderID, _ := result.LastInsertId()

			// Simpan detail pesanan dengan ON DUPLICATE KEY UPDATE
			for _, menuID := range menuIDs {
				quantity := r.FormValue("quantity_" + menuID)
				_, err := db.Exec(`
					INSERT INTO orderDetails (order_id, menu_id, quantity) 
					VALUES (?, ?, ?)
					ON DUPLICATE KEY UPDATE quantity = VALUES(quantity)`, orderID, menuID, quantity)
				if err != nil {
					http.Error(w, "Error inserting order details: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			// Redirect ke halaman checkout
			http.Redirect(w, r, "/checkOut", http.StatusSeeOther)
		}
	}
}

func CheckoutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ambil session pengguna
		session, _ := store.Store.Get(r, "session-name")

		userID, userIDok := session.Values["user_id"].(int)
		username, usernameok := session.Values["username"].(string)
		if !userIDok || !usernameok {
			http.Redirect(w, r, "/verfiedUserLogin", http.StatusSeeOther)
			return
		}

		// Ambil pesanan terbaru
		var orderID int
		err := db.QueryRow("SELECT idPesanan FROM pesanan WHERE id_User = ? ORDER BY created_at DESC LIMIT 1", userID).Scan(&orderID)
		if err != nil {
			http.Error(w, "Error fetching order: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Ambil detail pesanan
		rows, err := db.Query(`
			SELECT m.item_name, od.quantity, m.harga
			FROM orderDetails od
			JOIN menu m ON od.menu_id = m.idMenu
			WHERE od.order_id = ?`, orderID)
		if err != nil {
			http.Error(w, "Error fetching order details: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var orderItems []OrderItem
		var total float64

		for rows.Next() {
			var item OrderItem
			if err := rows.Scan(&item.ItemName, &item.Quantity, &item.Price); err != nil {
				http.Error(w, "Error scanning order details: "+err.Error(), http.StatusInternalServerError)
				return
			}
			orderItems = append(orderItems, item)
			total += float64(item.Quantity) * item.Price
		}

		data := struct {
			UserData   UserData
			OrderItems []OrderItem
			Total      float64
		}{
			UserData: UserData{
				UserID:   userID,
				Username: username,
			},
			OrderItems: orderItems,
			Total:      total,
		}

		// Render template checkOut.html
		fp := filepath.Join("app", "views", "checkOut.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err = tmpl.Execute(w, data); err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/payed", http.StatusSeeOther)
	}
}
