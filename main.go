package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/debugger1709/blog-aggregator/internal/config"
	"github.com/debugger1709/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal("cannot open database")
	}
	dbQuerries := database.New(db)
	st := state{
		cfg: &cfg,
		db:  dbQuerries,
	}
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerAllUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerAllFeeds)
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}
	err = cmds.run(&st, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
