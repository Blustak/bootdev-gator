package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/Blustak/bootdev-gator/internal/config"
	"github.com/Blustak/bootdev-gator/internal/database"
)

type state struct {
    db *database.Queries
    userConfig *config.Config
}

type command struct {
    name string
    args []string
}

type commands struct {
    cmd map[string]func(*state,command) error
}

func (c *commands) run(s *state, cmd command) error {
    mappedCmd, ok := c.cmd[cmd.name]
    if !ok {
        return fmt.Errorf("command %s is not registered",cmd.name)
    }
    if err := mappedCmd(s,cmd); err != nil {
        return err
    }
    return nil
}

func (c *commands) register(name string, f func(*state,command) error) {
    _,ok := c.cmd[name]
    if ok {
        fmt.Printf("warning: command %s is being overwritten\n", name)
    }
    c.cmd[name] = f
}

func main() {
    cfgFile, err := config.ReadUserConfig()
    if err != nil {
        fmt.Printf("error: %v\n", err)
    }
    currentState := state{ userConfig:&cfgFile}

    db,err := sql.Open("postgres",cfgFile.DbUrl)
    if err != nil {
        fmt.Printf("error: %v", err)
    }
    currentState.db = database.New(db)

    cmds := commands {
        cmd: make(map[string]func(*state,command) error),
    }
    cmds.register("login", handlerLogin)
    cmds.register("register",handlerRegister)
    cmds.register("reset",handlerReset)
    cmds.register("users",handlerUsers)

    if len(os.Args) < 2 {
        fmt.Println("error: not enough arguments")
        os.Exit(1)
    }
    cmd := command{
        name: os.Args[1],
        args: os.Args[2:],
    }
    if err := cmds.run(&currentState,cmd); err != nil {
        fmt.Printf("error: %v\n",err)
        os.Exit(1)
    }

}

func handlerLogin(s *state,cmd command) error {
    if len(cmd.args) == 0 {
        return fmt.Errorf("login requires a username")
    }
    user,err := s.db.GetUser(
        context.Background(),
        cmd.args[0],
    )
    if err != nil {
        return err
    }
    if err := s.userConfig.SetUser(user.Name); err != nil {
        return err
    }

    return nil
}

func handlerRegister(s *state,cmd command) error {
    if len(cmd.args) == 0 {
        return fmt.Errorf("register requires a username")
    }
    nowTime := time.Now()
    user,err := s.db.CreateUser(
        context.Background(),
        database.CreateUserParams{
            ID: uuid.New(),
            CreatedAt: nowTime,
            UpdatedAt: nowTime,
            Name: cmd.args[0],
    })
    if err != nil {
        return err
    }
    if err := s.userConfig.SetUser(user.Name); err != nil {
        return err
    }
    fmt.Printf("Created user %s\nData:%v\n",user.Name,user)
    return nil


}

func handlerUsers(s *state, cmd command) error {
    users, err := s.db.GetUsers(
        context.Background(),
    )
    if err != nil {
        return err
    }
    for _,u := range users {
        if u.Name == s.userConfig.CurrentUserName {
            fmt.Printf("* %s (current)\n",u.Name)
        } else {
            fmt.Printf("* %s\n", u.Name)
        }
    }
    return nil
}

func handlerReset(s *state, cmd command) error {
    if err := s.db.Reset(context.Background()); err != nil {
        return err
    }
    fmt.Println("successfully reset users table")
    return nil
}
