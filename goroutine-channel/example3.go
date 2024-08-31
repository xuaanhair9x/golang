package main

import (
	"fmt"
	//"time"
)

type Tag struct {
	Name, Type string
}

type Settings struct {
	NotificationsEnabled bool
}

type User struct {
	Id, Name, LastName, Status string
	Tags                       []*Tag
	*Settings
}

type NotificationsService struct {
}

//userToNotify <-chan *User // read-only channel
//userToNotify chan<- *User // send-only channel
func main() {
	usersToUpdate := make(chan []*User)
	userToNotify := make(chan *User)
	newUsers := []*User{
		{Name: "John", Status: "active", Settings: &Settings{NotificationsEnabled: true}},
		{Name: "Carl", Status: "active", Settings: &Settings{NotificationsEnabled: false}},
		{Name: "Paul", Status: "deactive", Settings: &Settings{NotificationsEnabled: true}},
		{Name: "Sam", Status: "active", Settings: &Settings{NotificationsEnabled: true}},
	}
	existingUsers := []*User{
		{Name: "Jessica", Status: "active", Settings: &Settings{NotificationsEnabled: true}},
		{Name: "Eric", Status: "active", Settings: &Settings{NotificationsEnabled: true}},
		{Name: "Laura", Status: "active", Settings: &Settings{NotificationsEnabled: true}},
	}

	go filterNewUsersByStatus(usersToUpdate, newUsers)
	go updateUsers(usersToUpdate, userToNotify, existingUsers)
	notifyUsers(userToNotify, existingUsers)
}

func filterNewUsersByStatus(usersToUpdate chan<- []*User, users []*User) {
	defer close(usersToUpdate)
	filteredUsers := []*User{}
	for _, user := range users {
		if user.Status == "active" && user.Settings.NotificationsEnabled {
			filteredUsers = append(filteredUsers, user)
		}
	}
	fmt.Println("1.Before usersToUpdate <- filteredUsers")
	usersToUpdate <- filteredUsers
}

func updateUsers(usersToUpdate <-chan []*User, userToNotify chan<- *User, users []*User) {
	defer close(userToNotify)
	for _, user := range users {
		user.Tags = append(user.Tags, &Tag{Name: "UserNotified", Type: "Notifications"})
	}
	fmt.Println("2.Before newUsers := <-usersToUpdate")
	newUsers := <-usersToUpdate
	// newUsers := []*User{
	// 	{Name: "John", Status: "active", Settings: &Settings{NotificationsEnabled: true}},
	// 	{Name: "Carl", Status: "active", Settings: &Settings{NotificationsEnabled: false}},
	// }
	for _, user := range newUsers {
		//time.Sleep(1 * time.Second)
		user.Tags = append(user.Tags, &Tag{Name: "NewNotification", Type: "Notifications"})
		fmt.Println("3.Before userToNotify <- user")
		userToNotify <- user
		fmt.Println("3.after userToNotify <- user")
	}
}

func notifyUsers(userToNotify <-chan *User, users []*User) {
	service := &NotificationsService{}
	for user := range userToNotify {
		service.SendEmailNotification(user, "Tags", "You got your first tag!!")
	}
	for _, user := range users {
		service.SendEmailNotification(user, "Tags", "A new tag has been added to your profile!!")
	}
}

func (n *NotificationsService) SendEmailNotification(user *User, title, message string) {
	fmt.Printf("Email Notification Sent to %v, Hi %s, %s\n", user, user.Name, message)
}