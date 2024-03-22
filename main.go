package main

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type StatusData struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

var statusMutex sync.Mutex
var currentStatus StatusData

func updateStatus() {
	for {
		statusMutex.Lock()
		// Mengubah nilai water dan wind secara acak dalam rentang yang diinginkan
		currentStatus.Water = generateRandomNumber(1, 10)
		currentStatus.Wind = generateRandomNumber(1, 20)
		statusMutex.Unlock()
		time.Sleep(3 * time.Second) // Update setiap 5 detik
	}
}

func generateRandomNumber(min, max int) int {
	// Menghasilkan angka acak antara min dan max
	return rand.Intn(max-min+1) + min
}

func handleDataRequest(c *gin.Context) {
	statusMutex.Lock()
	defer statusMutex.Unlock()

	// Mengirim data status saat ini sebagai JSON
	c.JSON(http.StatusOK, gin.H{"status": currentStatus})
}

func main() {
	// Membuat instance Gin
	r := gin.Default()

	// Menginisialisasi status awal
	currentStatus = StatusData{
		Water: generateRandomNumber(1, 10),
		Wind:  generateRandomNumber(1, 20),
	}

	// Memulai goroutine untuk memperbarui status secara periodik
	go updateStatus()

	// Mengatur endpoint untuk mendapatkan data
	r.GET("/data", handleDataRequest)

	// Mengatur route untuk halaman utama
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Menjalankan server pada port 8080
	r.Run(":8080")
}
