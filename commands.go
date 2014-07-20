package main

import (
  "github.com/codegangsta/cli"
  "os"
  "runtime"
  "io/ioutil"
  "regexp"
  "strings"
  "text/template"
  "log"
  "encoding/json"
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
  Flags: []cli.Flag {
    cli.StringFlag { Name: "dir, d", Value: getOsHomeDir(), Usage: "Directory where you put .gitconfigs" },
  },
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

// variables
var okay = []string{"y", "Y", "yes", "Yes", "YES"}
var no = []string{"n", "Y", "no", "No", "NO"}

// templates
type Configuration struct {
  HomeDir string
}

type Gitconfig struct {
  EnvName string
}

// .gitconf
var gitConfTemplate = template.Must(ParseAsset(".gitconf", "templates/.gitconf.tmpl"))
var gitConf = Source {
  Name: ".gitconf",
  Template: *gitConfTemplate,
}
// .gitconfig
var gitconfigTemplate = template.Must(ParseAsset(".gitconfig", "templates/.gitconfig.tmpl"))
var gitconfig = Source {
  Name: ".gitconfig",
  Template: *gitconfigTemplate,
}
// .gitconfig_global
var gitconfigGlobalTemplate = template.Must(ParseAsset(".gitconfig_global", "templates/.gitconfig_global.tmpl"))
var gitconfigGlobal = Source {
  Name: ".gitconfig_global",
  Template: *gitconfigGlobalTemplate,
}
// .gitconfig.local
var gitconfigLocalTemplate = template.Must(ParseAsset(".gitconfig.local", "templates/.gitconfig.local.tmpl"))
var gitconfigLocal = Source {
  Name: ".gitconfig.local",
  Template: *gitconfigLocalTemplate,
}

// methods
func doInit (c *cli.Context) {
  // set user home directory
  home_dir := setUserHomeDir(c.String("dir"))
  config_file := getOsHomeDir() + "/.gitconfig" 

  // rename user global .gitconfig to .gitconfig.local if it exists
  if file, err := ioutil.ReadFile(config_file); err == nil && isFileExist(config_file) {
    // copy .gitconfig to .gitconfig.local
    if error := ioutil.WriteFile(home_dir + "/.gitconfig.local", file, 0755); error != nil {
      // errors
      println("Cannot create file: " + home_dir + "/.gitconfig.local")
      log.Fatal(error)
    }
  } else {
    // create new .gitconfig.local file
    if error := gitconfigLocal.generate(getUserHomeDir(), ""); error != nil {
      // errors
      println("Cannot create file: " + home_dir + "/.gitconfig.local")
      log.Fatal(error)
    }
  }

  // create .gitconfig (to include other files) and .gitconfig_global for user global
  if err := gitconfigGlobal.generate(getOsHomeDir(), ""); err == nil {
    println("~/.gitconfig_global is created")
  } else {
    println("Cannot create ~/.gitconfig_global")
    log.Fatal(err)
  }

  envConf := Gitconfig { ".gitconfig.local" }
  if err := gitconfig.generate(getOsHomeDir(), envConf); err == nil {
    println("~/.gitconfig is created")
  } else {
    println("Cannot create ~/.gitconfig")
    log.Fatal(err)
  }
}

func doCreate (c *cli.Context) {
  homeDir := getUserHomeDir()

  println(homeDir)
}

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

func setUserHomeDir (homeDir ...string) string {
  if len(homeDir) < 1 {
    homeDir = append(homeDir, getOsHomeDir())
  }

  config := Configuration { homeDir[0] }

  if err := gitConf.generate(getOsHomeDir(), config); err == nil {
    println("~/.gitconf is created")
  } else {
    log.Fatal(err)
  }

  return homeDir[0]
}

func getUserHomeDir () string {
  file, err := ioutil.ReadFile(getOsHomeDir() + "/.gitconf")

  if err != nil {
    println("Cannot open file ~/.gitconf", err.Error())
    os.Exit(1)
  }

  var config Configuration
  e := json.Unmarshal(file, &config)
  if e != nil {
    println("Cannot parse .gitconf", err.Error())
    os.Exit(1)
  }

  return config.HomeDir
}

func getOsHomeDir () string {
  if runtime.GOOS == "windows" {
    home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
    if home == "" {
      home = os.Getenv("USERPROFILE")
    }

    return home
  }

  return os.Getenv("HOME")
}

func ParseAsset(name string, path string) (*template.Template, error) {
  src, err := Asset(path)
  if err != nil {
    return nil, err
  }

  return template.New(name).Parse(string(src))
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
