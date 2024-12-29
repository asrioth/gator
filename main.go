package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/asrioth/gator/internal/config"
	"github.com/asrioth/gator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const Login string = "login"
const Register string = "register"
const Reset string = "reset"
const Users string = "users"

func initializeCommands() Commands {
	commands := Commands{Commands: make(map[string]func(*State, Command) error), BaseCommands: make(map[string]Command)}
	commands.Register(Login, 1, handlerLogin)
	commands.Register(Register, 1, handlerRegister)
	commands.Register(Reset, 0, handlerReset)
	commands.Register(Users, 0, handlerUsers)
	return commands
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Command name required")
		os.Exit(1)
	}
	commands := initializeCommands()

	configData, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	state := State{ConfigData: &configData}
	db, err := sql.Open("postgres", state.ConfigData.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	state.Db = database.New(db)

	for i := 1; i < len(os.Args); i++ {
		cmdWord := sanitizeWord(os.Args[i])
		cmdBase, ok := commands.BaseCommands[cmdWord]
		if !ok {
			fmt.Println("Unrecognised command")
			os.Exit(1)
		}
		if len(os.Args)-2 < cmdBase.NumberOfArguments {
			fmt.Println("Not enough arguments for command")
			os.Exit(1)
		}
		for argCount := 0; argCount < cmdBase.NumberOfArguments; argCount++ {
			i++
			cmdBase.Arguments[argCount] = os.Args[i]
		}
		err = commands.Run(&state, cmdBase)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func sanitizeWord(word string) string {
	word = strings.ToLower(word)
	word = strings.TrimSpace(word)
	return word
}

func handlerLogin(s *State, cmd Command) error {
	_, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}
	err = s.ConfigData.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User name has been set to %v\n", cmd.Arguments[0])
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	_, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err == nil {
		return errors.New("a user with that name already exits")
	}

	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Arguments[0]})
	if err != nil {
		return err
	}
	err = s.ConfigData.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}
	println(user.ID.String(), user.CreatedAt.String(), user.UpdatedAt.String(), user.Name)
	return nil
}

func handlerReset(s *State, cmd Command) error {
	err := s.Db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Database succesfully deleted all users")
	return nil
}

func handlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		current := ""
		if s.ConfigData.CurrentUserName == user {
			current = " (current)"
		}
		fmt.Printf(" * %v%v\n", user, current)
	}
	return nil
}
