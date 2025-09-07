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
	"github.com/Blustak/bootdev-gator/internal/rss"
)

type state struct {
	db         *database.Queries
	userConfig *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmd map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	mappedCmd, ok := c.cmd[cmd.name]
	if !ok {
		return fmt.Errorf("command %s is not registered", cmd.name)
	}
	if err := mappedCmd(s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.cmd[name]
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
	currentState := state{userConfig: &cfgFile}

	db, err := sql.Open("postgres", cfgFile.DbUrl)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	currentState.db = database.New(db)

	cmds := commands{
		cmd: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
    cmds.register("following",middlewareLoggedIn(handlerFollowing))
    cmds.register("unfollow",middlewareLoggedIn(handlerUnfollow))

	if len(os.Args) < 2 {
		fmt.Println("error: not enough arguments")
		os.Exit(1)
	}
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	if err := cmds.run(&currentState, cmd); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login requires a username")
	}
	user, err := s.db.GetUser(
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

func handlerRegister(s *state, cmd command) error {
	if err := checkArgs(1, cmd); err != nil {
		return err
	}
	nowTime := time.Now()
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
			Name:      cmd.args[0],
		})
	if err != nil {
		return err
	}
	if err := s.userConfig.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Printf("Created user %s\nData:%v\n", user.Name, user)
	return nil

}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(
		context.Background(),
	)
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.Name == s.userConfig.CurrentUserName {
			fmt.Printf("* %s (current)\n", u.Name)
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

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.NewFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", feed)
	return nil

}


func handlerFeeds(s *state, cmd command) error {
	if err := checkArgs(0, cmd); err != nil {
		return err
	}
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Feed %v:\n", feed.ID)
		fmt.Printf("\tName: %s\n", feed.Name)
		fmt.Printf("\tURL: %s\n", feed.Url)
		fmt.Printf("\tUsername: %s\n", user.Name)
	}
	return nil
}

// logged in functions

func handlerUnfollow(s *state, cmd command, user database.User) error {
    if err := checkArgs(1,cmd); err != nil {
        return err
    }
    feed, err := s.db.GetFeedByUrl(context.Background(),cmd.args[0])
    if err != nil {
        return err
    }
    if err := s.db.UnfollowFeed(context.Background(),database.UnfollowFeedParams{
        UserID: user.ID,
        FeedID: feed.ID,
    }); err != nil {
        return err
    }
    return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if err := checkArgs(1, cmd); err != nil {
		return err
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	nowTime := time.Now()
	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	fmt.Printf("Feed name: %s\nFeed URL:%s\n", feed.Name, feed.Url)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	res, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
    fmt.Printf("User %s is currently following:\n",user.Name)
	for _, feed := range res {
        fmt.Printf("\t - %s\n",feed.FeedName)
	}
    return nil
}

func handlerAddFeed(s *state, cmd command, user database.User ) error {
	if err := checkArgs(2, cmd); err != nil {
		return err
	}
	nowTime := time.Now()
	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
    s.db.CreateFeedFollow(context.Background(),database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: nowTime,
        UpdatedAt: nowTime,
        UserID: user.ID,
        FeedID: feed.ID,

    })
	fmt.Printf("%v\n", feed)
	return nil

}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
    return func (s *state, cmd command) error {
        user,err := s.db.GetUser(context.Background(),s.userConfig.CurrentUserName)
        if err != nil {
            return err
        }
        if err := handler(s,cmd,user); err != nil {
            return err
        }
        return nil
    }
}

func checkArgs(expectedLength int, cmd command) error {
	if len(cmd.args) != expectedLength {
		return fmt.Errorf("%s expected %d args; got %d.", cmd.name, expectedLength, len(cmd.args))
	}
	return nil
}
