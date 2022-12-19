package service

import (
	"context"
	"fmt"
	"github.com/deemanrip/promotions/repository"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

const insertIntoTmpTableQuery = `INSERT INTO promotions.promotions_tmp
SELECT *
FROM s3('http://object-storage:9000/promotions/%v',
	'app',
	'test12345',
	'CSV',
	'id UUID, price Decimal64(6), expiration_date String'
)`

const exchangeTablesQuery = "EXCHANGE TABLES promotions.promotions_tmp AND promotions.promotions"
const truncateQuery = "TRUNCATE TABLE IF EXISTS promotions.promotions_tmp"

func CreateClient() *minio.Client {
	minioClient, err := minio.New("object-storage:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("app", "test12345", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return minioClient
}

func ListenNotifications(minioClient *minio.Client) {
	for notificationInfo := range minioClient.ListenBucketNotification(context.Background(), "promotions", "", "", []string{
		"s3:ObjectCreated:*",
	}) {
		if notificationInfo.Err != nil {
			log.Error(notificationInfo.Err)
		}
		for _, record := range notificationInfo.Records {
			processObjectCreated(&record.S3.Object.Key)
		}
	}
}

func processObjectCreated(createdFileName *string) {
	insertQuery := fmt.Sprintf(insertIntoTmpTableQuery, *createdFileName)
	clickhouseConn := repository.ClickhouseConnection
	if truncateError := clickhouseConn.Exec(context.Background(), truncateQuery); truncateError != nil {
		log.Error(truncateError)
		return
	}
	if insertErr := clickhouseConn.Exec(context.Background(), insertQuery); insertErr != nil {
		log.Error(insertErr)
		return
	}
	if exchangeErr := clickhouseConn.Exec(context.Background(), exchangeTablesQuery); exchangeErr != nil {
		log.Error(exchangeErr)
		return
	}
}
