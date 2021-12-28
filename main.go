package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

var email string

func main() {
	var _email = flag.String("email", "anonymous", "enter user's email")
	var command = flag.String("command", "nothing", "enter command")
	flag.Parse()

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	email = *_email
	ctx := context.Background()
	uid := getUserUIDByEmail(email, ctx, client)

	if *command == "delete" {
		deleteUser(ctx, client, uid)
	} else {
		fmt.Printf("Email :%v\n", email)
		fmt.Printf("User ID :%v\n", uid)
	}
}

func getUserUIDByEmail(email string, ctx context.Context, client *auth.Client) string {
	u, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		log.Fatalf("error getting user by email %s: %v\n", email, err)
	}
	//log.Printf("Successfully fetched user data: \nemail:%v\nUID:%v\n", u.Email, u.UID)
	return u.UID
}

func deleteUser(ctx context.Context, client *auth.Client, uid string) {
	fmt.Printf("Careful !, Are you sure you want to delete %v (y/n) : ", email)
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	check(err)

	if char == 'y' {
		err = client.DeleteUser(ctx, uid)
		if err != nil {
			fmt.Println("Sorry, there is something wrong with me")
			log.Fatal(err.Error())
		} else {
			fmt.Printf("Yeay, user %v has been deleted\n", email)
		}
	} else {
		fmt.Println("Operation canceled, bye")
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
