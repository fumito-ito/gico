package main

import (
  "github.com/codegangsta/cli"
  "os"
  "runtime"
  "io/ioutil"
  "regexp"
  "strings"
)

var Commands = []cli.Command{
  commandInit,
  commandCreate,
  commandDelete,
  commandSwitch,
  commandList,
}

var commandInit = cli.Command{
  Name: "init",
  Usage: "Initial setup",
  Description: `
  Initial setup for your .gitconfig environment.
  This command creates $HOME/dotfiles directory if it doesn't exist.
  `,
  Action: doInit,
}

var commandCreate = cli.Command{
  Name: "create",
  Usage: "Create new .gitconfig",
  Description: `
  Create new .gitconfig environment files.
  This command creates some files in $HOME/dotfiles/[env name].
  
  If you want to override existing env, you have to use -f option.
  `,
  Action: doCreate,
}

var commandDelete = cli.Command{
  Name: "delete",
  Usage: "Delete .gitconfig file",
  Description: `
  Delete your .gitconfig environment files.
  This command deletes $HOME/dotfiles/[env name] directory.
  `,
  Action: doDelete,
}

var commandSwitch = cli.Command{
  Name: "switch",
  Usage: "Switch .gitconfig",
  Description: `
  Switch your .gitconfig environment.
  This command switches your .gitconfig environment with [env name].
  
  If you want to use git without any environment, you run
  
  gitconf switch default
  
  It will return default .gitconfig.
  `,
  Action: doSwitch,
}

var commandList = cli.Command{
  Name: "list",
  Usage: "List your all .gitconfig",
  Description: `
  List your .gitconfig environment.
  This command lists your .gitconfig environments in console.
  `,
  Action: doList,
}

func doInit (c *cli.Context) {}

func doCreate (c *cli.Context) {}

func doDelete (c *cli.Context) {}

func doSwitch (c *cli.Context) {}

func doList (c *cli.Context) {
  homeDir := getUserHomeDir() + "/dotfiles"

  if isFileExist(homeDir) {
    files, _ := ioutil.ReadDir(homeDir)
    for _, f := range files {
      fileName := f.Name()

      if matched, _ := regexp.MatchString(".*\\.gitconfig$", fileName); matched {
        println(strings.Replace(fileName, ".gitconfig", "", 1))
      }
    }
  }
}

func getUserHomeDir () string {
  if runtime.GOOS == "windows" {
    home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
    if home == "" {
      home = os.Getenv("USERPROFILE")
    }

    return home
  }

  return os.Getenv("HOME")
}

func isFileExist (filename string) bool {
  if _, err := os.Stat(filename); err == nil {
    return true
  } else {
    return false
  }
}