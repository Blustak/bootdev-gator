package main

import (
	"fmt"
	"os"

	"github.com/Blustak/bootdev-gator/internal/config"
)

type state struct {
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
        fmt.Println("warning: command %s is being overwritten", name)
    }
    c.cmd[name] = f
}

func main() {
    cfgFile, err := config.ReadUserConfig()
    if err != nil {
        fmt.Printf("error: %v\n", err)
    }
    currentState := state{ userConfig:&cfgFile}
    cmds := commands {
        cmd: make(map[string]func(*state,command) error),
    }
    cmds.register("login", handlerLogin)

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
    if err := s.userConfig.SetUser(cmd.args[0]); err != nil {
        return err
    }
    return nil
}
