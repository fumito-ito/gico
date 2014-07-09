package main

import (
  "github.com/codegangsta/cli"
)

var Commands = []cli.Command{
  commandList,
}

var commandList = cli.Command{
  Name: "list",
  Usage: "",
  Description: ``,
  Action: doList,
}

func doList (c *cli.Context) {
}