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

// templates
type Configuration struct {
  HomeDir string
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
  setUserHomeDir(c.String("dir"))
  // read user original configuration file
  var originalFile = getOsHomeDir() + "/.gitconfig"

  // rename user global .gitconfig to .gitconfig.local if it exists
  if file, err := ioutil.ReadFile(originalFile); err == nil && isFileExist(originalFile) {
    // default config file name
    var defaultFile = getUserHomeDir() + "/.gitconfig.local"
    // copy .gitconfig to .gitconfig.local
    if error := ioutil.WriteFile(defaultFile, file, 0755); error != nil {
      // errors
      println("Cannot create file: " + defaultFile)
      log.Fatal(error)
    }
  } else {
    // create new .gitconfig.local file
    createLocalConfigFile("local")
  }

  // create .gitconfig (to include other files) and .gitconfig_global for user global
  if err := gitconfigGlobal.generate(getOsHomeDir(), Configuration {}); err != nil {
    println("Cannot create ~/.gitconfig_global")
    log.Fatal(err)
  }

  switchConfigFile("local")
}

func doCreate (c *cli.Context) {
  // create configuration file from arguments
  if len(c.Args()) > 0 {
    for _, envName := range c.Args() {
      createLocalConfigFile(envName)
      doUse(c)
    }
  } else {
    println("Set 1 or more arguments to create config files")
  }
}

func doDelete (c *cli.Context) {}

func doUse (c *cli.Context) {
  // set configuration file to use by first argument
  if len(c.Args()) > 0 {
    var envName = c.Args()[0]
    switchConfigFile(envName)

    println("Now using [" + envName + "]")
  } else {
    println("Set argument to set environment")
  }
}

func doList (c *cli.Context) {
  var homeDir = getUserHomeDir()

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

func createLocalConfigFile (envName string) {
  gitconfigLocal.Name = ".gitconfig." + envName

  if error := gitconfigLocal.generate(getUserHomeDir(), Configuration {}); error != nil {
    // errors
    println("Cannot create file: " + getUserHomeDir() + "/.gitconfig." + envName)
    log.Fatal(error)
  }
}

func switchConfigFile (envName string) {
  var config = Configuration {
    EnvName : envName,
  }

  if err := gitconfig.generate(getOsHomeDir(), config); err != nil {
    println("Fialed to switch to " + envName)
    log.Fatal(err)
  }
}

func setUserHomeDir (homeDir ...string) {
  if len(homeDir) < 1 {
    homeDir = append(homeDir, getOsHomeDir())
  }

  var config = Configuration {
    HomeDir : homeDir[0],
  }

  if err := gitConf.generate(getOsHomeDir(), config); err != nil {
    println("Fialed to set user home directory")
    log.Fatal(err)
  }
}

func getUserHomeDir () string {
  // read file to find home directory
  var file, err = ioutil.ReadFile(getOsHomeDir() + "/.gitconf")

  if err != nil {
    println("Cannot open file ~/.gitconf", err.Error())
    log.Fatal(err)
  }

  // parse json to configuration to return homeDir value
  var config Configuration
  e := json.Unmarshal(file, &config)
  if e != nil {
    println("Cannot parse .gitconf", err.Error())
    log.Fatal(err)
  }

  return config.HomeDir
}

func getOsHomeDir () string {
  if runtime.GOOS == "windows" {
    var home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
    if home == "" {
      home = os.Getenv("USERPROFILE")
    }

    return home
  }

  return os.Getenv("HOME")
}

func ParseAsset(name string, path string) (*template.Template, error) {
  var src, err = Asset(path)
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
