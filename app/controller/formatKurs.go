package controller

import (
	"fmt"
	"strings"
)

// Fungsi untuk format harga menjadi mata uang IDR
func FormatIDR(price float64) string {
	idr := fmt.Sprintf("Rp %s", formatNumber(price))
	return idr
}

// Fungsi pembantu untuk format angka dengan koma
func formatNumber(price float64) string {
	str := fmt.Sprintf("%.0f", price) // Mengubah float ke string tanpa desimal
	var result []string
	count := 0

	for i := len(str) - 1; i >= 0; i-- {
		if count == 3 {
			result = append(result, ".")
			count = 0
		}
		result = append(result, string(str[i]))
		count++
	}

	// Membalikkan urutan hasil dan menggabungkan menjadi string
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return strings.Join(result, "")
}
