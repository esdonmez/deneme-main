package main

import (
	"context"
	"fmt"
	"log"
	"pgo-collector"
)

func main() {
	/*
		TODO:
			2- cmd'den çağrılabilecek bir kütüphane oluştur. Minio'dan okuyup dosyaları default.pgo adı altında birleştirsin ve çağrılan yerde local'e indirsin
			3- bu dosyayla build başlatmayı deneyelim
	*/
	ctx := context.Background()

	endpoint := "localhost:9000"
	accessKeyID := "admin"
	secretAccessKey := "admin123"

	minioUploader := pgo.NewUploader(endpoint, accessKeyID, secretAccessKey, "")

	// Make a new bucket called pgos.
	bucketName := "pgos"
	location := "us-east-1"

	minioUploader.Create(ctx, bucketName, location)

	collectEndpoint := "https://discovery-indexing-buybox-score-service.moon.trendyol.com/debug/pprof/profile?seconds=30"
	collector := pgo.NewCollector(collectEndpoint, "* * * * *", minioUploader, bucketName)
	err := collector.Start()
	if err != nil {
		log.Printf("Error starting cron")
	}

	_, _ = fmt.Scanln()
}
