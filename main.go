package main

import (
	"context"
	"fmt"
	"github.com/Neulhan/go-prisma2-tuto/db"
	"log"
)

func main() {
	fmt.Println("my new project")
	client := db.NewClient()
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := client.Disconnect()
		if err != nil {
			panic(fmt.Errorf("could not disconnect %w", err))
		}
	}()

	ctx := context.Background()
	// create a user
	createdUser, err := client.User.CreateOne(
		db.User.Email.Set("john.doe@example.com"),
		db.User.Name.Set("John Doe"),
		// ID is optional, which is why it's specified last. if you don't set it
		// an ID is auto generated for you
		//db.User.ID.Set("123"),
	).Exec(ctx)

	log.Printf("created user: %+v", createdUser)

	// find a single user
	user, err := client.User.FindOne(
		db.User.Email.Equals("john.doe@example.com"),
	).Exec(ctx)
	if err != nil {
		panic(err)
	}

	// for optional/nullable values, you need to check the function and create two return values
	// `name` is a string, and `ok` is a bool whether the record is null or not. If it's null,
	// `ok` is false, and `name` will default to Go's default values; in this case an empty string (""). Otherwise,
	// `ok` is true and `name` will be "John Doe".
	name, ok := user.Name()

	if !ok {
		log.Printf("user's name is null")
		return
	}

	log.Printf("The users's name is: %s", name)
}
