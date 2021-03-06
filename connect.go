package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const serviceParth = "./service.json"

//Conect : function return firestore obj
func Conect() (*firestore.Client, context.Context) {
	ctx := context.Background()
	sa := option.WithCredentialsFile(serviceParth)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client, ctx
}
