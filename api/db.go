package api

import (
	"context"
	"log"
	"net/http"

	"encoding/json"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

const creds = `{
	"type": "service_account",
	"project_id": "ngosangns-authapp",
	"private_key_id": "d05a9d531cfe6a540e3d263f7f168f105aac438b",
	"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCWTytFMPgFs4lV\n0hn3ILOni/se57pbhQaICDssy8dZdCMpX3xAW/nQftpOpVTOFv49cOZfgbHFlQSS\nhtytcVAgZlksnMhCAHh5L4zwtmz6WmjZrbROn+u9wKvkWz2ih9rVzGHzgSK/Vxtf\ne3ePY0F21bEly7krMDJUZxplo2EtOauA25SKzOL+YfWMl2EHi2OkfvskFvqea+iL\nKM2XC/i3W29WU+Rdic2shu9JAf+oW00S525/QZ0OCRG7Vu9mApOI1Ium4C1NMA1m\nZUEj3soD+uG+QwQAr12IGadvZVuX3yUWQ9Gdiv3vO0nDldMPFFBTEg2IhHDqqKc8\n/Xl05aiLAgMBAAECggEAK1K3VJKYMkB8uZU8saulzIl+wHonDyDo8IV61N8JXuGo\nDtE8RWqdNviRILDroBklf30OkYGWxS296yPe8EjkQOpvEnpACDINh1RqyaizAXfU\nl8VG3pCeiItDoCaiPbntm8j0hJR4ePD47UzveY8wu7k8/XlaLjYOk6BCDgCIzllL\nTu4CAsz3ROtjQcs6tSye1Nb34E+EInVkyMcxx8ye6aSiwSXVpspmnmLNozTlWl/B\npMGZ/whIUxkExC+raBVlQ0+YFUq1zTZtUBmrPqUd0AwDYJ0UCDhCaoeNK7OrKyvE\n9CF0eB81EGt8y4OovNngDCYqo/e1qTcnrDWwwIY+XQKBgQDKukmX7HJriv0z6ZD0\nbMd9LPA5QYfTWMe7xEPSKPTowtUeac355Fh1NTDYUuOZTm6No+zALBJtqfM0ow/6\nu4xkBgnuPLTi92VDhtyLHju5SCVLTZUiPh4rZ5p5tUVOIVwtiE+Al26KzrV5nB2O\nzfGeSxJQ6teZdCP+pkIfkfZT/QKBgQC9zqLxGEIqgE6WwRWw7IzWHBS6SrRxoNC9\nMKiSMeeecU8F9hAX2i43Vwc2vDAo2KqT1AKCySQ62C2i9yKpLoiSVvisrSe30DO3\n7eZhmWbb4+sQCA+wNUIPl1wFTaO8u5wgI4qV1e35xPmP5thl25QMyQY88Vf0LMnO\nEDRsRHxhJwKBgH4Ft+H9VlOMH0K6GyYWyRwDZ4HwPqMfOWp74z0twLFBDILPev+w\n91xFKIwRWvLeUtMx5+a+fuY1E77Q7woKtIcpSdsTWc8Nw6FqoIbK0I2pT6W0INUc\nkFyjFuA5009yZX7YkbN6b3lSbnWemrE/TMf+GvC6TDRUglu2trwxXFNxAoGAH6ZF\nJU6WOeALOrgXldjb4xfrnOS4Efpu1B04/qezp4fdVXEiXdfTvQaV/VqD7UuzpdLE\nWGqRz/4cHgB3lx5bR4uZC7IT3WKqPRE+SSf7Ls9ictnlQ8ydp1vlzDyWAPKHphRH\nF5UEiZF+y7AChXmHxln+4EqeD+kn/WRWWXdSBiECgYEAv9T+sT82DIau13cszDTU\nBi4Zg3jCCpoDO9NyEg7shnH22UaUMhGNBW5AljzD6KiWwF08wXFf2js2la0mghym\nRXbPj+U6VjPlUDVOdVsDEJcAa+Vo4M+Db5OTR1dmIBsHbvce8NNIXYtQeFx6KGN/\nJ8mrasjr3PBN/hjI4d3XvfM=\n-----END PRIVATE KEY-----\n",
	"client_email": "firebase-adminsdk-5ugft@ngosangns-authapp.iam.gserviceaccount.com",
	"client_id": "106763696388976129492",
	"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-5ugft%40ngosangns-authapp.iam.gserviceaccount.com"
}`

// Db Handler
func Db(w http.ResponseWriter, r *http.Request) {
	client, ctx, err := connectDb()
	defer client.Close()
	if err != nil {
		log.Println(err)
		b, _ := json.Marshal(err)
		w.Write(b)
	}

	// Create a record
	type Record struct {
		ID   string `firestore:"id"`
		Name string `firestore:"name"`
	}
	_, _, err = client.Collection("users").Add(ctx, Record{
		ID:   "00",
		Name: "Nguyen Van A",
	})
	if err != nil {
		log.Println(err)
		b, _ := json.Marshal(err)
		w.Write(b)
	}
	w.Write([]byte("Inserted a record"))
}

func connectDb() (*firestore.Client, context.Context, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "ngosangns-authapp", option.WithCredentialsJSON([]byte(creds)))

	// defer client.Close()
	if err != nil {
		return client, ctx, err
	}
	return client, ctx, nil
}

// // MongoDB
// const dbName = "ngosangns"
// const dbCollection = "authenticationapp"
// const dbConnectString = "mongodb+srv://ngosangns:jikmli@cluster0.oxs6m.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
// // Db Handler
// func Db(w http.ResponseWriter, r *http.Request) {
// 	// Connect database
// 	client, cancel, err := connectDatabase()
// 	defer cancel()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	collection := client.Database(dbName).Collection(dbCollection)

// 	_, err = collection.InsertOne(context.TODO(), bson.D{
// 		{"id", "01"},
// 		{"name", "Nguyen Van A"},
// 		{"email", "123@gmail.com"},
// 	})
// 	if err != nil {
// 		fmt.Fprintf(w, err.Error())
// 		log.Fatal(err)
// 	}
// 	fmt.Fprintf(w, "Inserted a record!")
// }
// func connectDatabase() (*mongo.Client, context.CancelFunc, error) {
// 	clientOpts := options.Client().ApplyURI(dbConnectString)
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	client, err := mongo.Connect(ctx, clientOpts)
// 	if err != nil {
// 		defer cancel()
// 		return nil, nil, err
// 	}
// 	return client, cancel, nil
// }
