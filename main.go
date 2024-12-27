package main

import (
	"fmt"

	"github.com/asrioth/gator/internal/config"
)

func main() {
	configData, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = configData.SetUser("asrioth")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("db url: %v\n", configData.DbUrl)
	fmt.Printf("user name: %v\n", configData.CurrentUserName)
}
