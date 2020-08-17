package main

import (
	"log"
	"os"

	crawler "github.com/siangyeh8818/seleium.exporter.juniuhome/internal/crawler"
	server "github.com/siangyeh8818/seleium.exporter.juniuhome/internal/server"
)

func main() {
	log.Println("Exporter is start ro running")
	account_email := os.Getenv("ACCOUNT_EMAIL")
	log.Printf("ACCOUNT_EMAIL : %s \n", account_email)

	account_password := os.Getenv("ACCOUNT_PASSWORD")
	log.Printf("ACCOUNT_PASSWORD : %s \n", account_password)

	interval_time := os.Getenv("SELEIUM_INTERNAL_TIME")
	log.Printf("SELEIUM_INTERNAL_TIME : %s \n", interval_time)

	go func() {
		crawler.CallSelium()
	}()

	server.Run_Exporter_Server()
}
