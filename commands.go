package main

import (
  "github.com/codegangsta/cli"
  "os"
  "runtime"
  "io/ioutil"
  "regexp"
  "strings"
  "log"
  "fmt"
)

var Commands = []cli.Command{
  commandInit,
  commandCreate,
  commandDelete,
  commandUse,
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

var commandUse = cli.Command{
  Name: "use",
  Usage: "use envname",
  Description: `
  Switch your .gitconfig environment.
  This command switches your .gitconfig environment with [env name].
  
  If you want to use git without any environment, you run
  
  gitconf switch default
  
  It will return default .gitconfig.
  `,
  Action: doUse,
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

var okay = []string{"y", "Y", "yes", "Yes", "YES"}
var no = []string{"n", "Y", "no", "No", "NO"}

func doInit (c *cli.Context) {
  homeDir := getUserHomeDir() + "/dotfiles"
  var response string

  if !isFileExist(homeDir) {
    println("Can I create home directory in " + homeDir + " ? [Y/n]")
    _, err := fmt.Scanln(&response)

    if err != nil {
      log.Fatal(err)
    }

    if containString(okay, response) {
      os.Mkdir(homeDir, 0755)
      println(homeDir + "is created")
    } else if containString(no, response) {
      println("This command must have " + homeDir + ", please try again")
      return
    } else {
      println("please type keys yes or no")
      return
    }
  }

  println("initializing global .gitconfig in " + homeDir + "...")
}

func doCreate (c *cli.Context) {}

func doDelete (c *cli.Context) {}

func doUse (c *cli.Context) {}

func doList (c *cli.Context) {
  homeDir := getUserHomeDir()

  if isDirExist(homeDir) {
    files, _ := ioutil.ReadDir(homeDir)
    for _, f := range files {
      fileName := f.Name()

      if matched, _ := regexp.MatchString("^\\.gitconfig\\..*$", fileName); matched {
        println(strings.Replace(fileName, ".gitconfig.", "", 1))
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

func isFileExist (fileName string) bool {
  if fileInfo, err := os.Stat(fileName); err == nil && !fileInfo.IsDir() {
    return true
  } else {
    return false
  }
}

func isDirExist (dirName string) bool {
  if fileInfo, err := os.Stat(dirName); err == nil && fileInfo.IsDir() {
    return true
  } else {
    return false
  }
}

func findString (slice []string, element string) int {
  for index, elem := range slice {
    if elem == element {
      return index
    }
  }

  return -1
}

func containString (slice []string, element string) bool {
  return !(findString(slice, element) == -1)
}