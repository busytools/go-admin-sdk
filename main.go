package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

	fmt.Println(*command)
	switch *command {

	case "delete":
		uid := getUserUIDByEmail(email, ctx, client)
		deleteUser(ctx, client, uid)

	case "create":
		createUser(ctx, client, email)

	default:
		uid := getUserUIDByEmail(email, ctx, client)
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

func createUser(ctx context.Context, client *auth.Client, email string) *auth.UserRecord {
	reader := bufio.NewReader(os.Stdin)

	/*
		fmt.Print("Enter email address : ")
		email, _ := reader.ReadString('\n')
		email = strings.Replace(email, "\n", "", -1)
		email = strings.Replace(email, "\r", "", -1)
	*/

	fmt.Print("Enter Password : ")
	password, _ := reader.ReadString('\n')
	password = strings.Replace(password, "\n", "", -1)
	password = strings.Replace(password, "\r", "", -1)

	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		Disabled(false)
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
	}
	log.Printf("Successfully created user: %v\n", u)
	return u
}

func check(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
