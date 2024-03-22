package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Data struct untuk menyimpan data water dan wind
type Data struct {
	Water      int    `json:"water"`
	Wind       int    `json:"wind"`
	WaterStat  string `json:"water_status"`
	WindStat   string `json:"wind_status"`
	UpdateTime string `json:"update_time"`
}

// UpdateJSON akan meng-update data JSON setiap 15 detik
func UpdateJSON() {
	for {
		data := Data{
			Water:      rand.Intn(12) + 1, // Nilai maksimal air adalah 12
			Wind:       rand.Intn(20) + 1, // Nilai maksimal angin adalah 20
			UpdateTime: time.Now().Format(time.RFC3339),
		}

		// Menentukan status air
		switch {
		case data.Water < 5:
			data.WaterStat = "Aman"
		case data.Water >= 5 && data.Water <= 8:
			data.WaterStat = "Siaga"
		default:
			data.WaterStat = "Bahaya"
		}

		// Menentukan status angin
		switch {
		case data.Wind < 6:
			data.WindStat = "Aman"
		case data.Wind >= 6 && data.Wind <= 15:
			data.WindStat = "Siaga"
		default:
			data.WindStat = "Bahaya"
		}

		// Menulis data ke file JSON
		file, err := os.Create("data.json")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		err = encoder.Encode(data)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		fmt.Println("File JSON berhasil di-update")

		time.Sleep(15 * time.Second)
	}
}

// Handler untuk menampilkan halaman HTML
func indexHandler(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlContent)
}

// Handler untuk mengambil data status dari file JSON
func statusHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("data.json")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var data Data
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	// Jalankan fungsi UpdateJSON sebagai goroutine
	go UpdateJSON()

	// Menyiapkan route
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/status", statusHandler)

	// Menjalankan server
	fmt.Println("Server berjalan di http://localhost:2323")
	http.ListenAndServe(":2323", nil)
}
