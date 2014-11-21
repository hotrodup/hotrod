package main

import (
  "os"
  "fmt"
  "gopkg.in/alecthomas/kingpin.v1"
  "github.com/fatih/color"
  "github.com/kyokomi/emoji"
)

const (
  VERSION = "0.0.1"
)

var (

  CHECKERED_FLAG = emoji.Sprintf(":checkered_flag:")
  RED_CAR = emoji.Sprintf(":red_car:")

  app = kingpin.New(color.YellowString("hotrod"), CHECKERED_FLAG + color.YellowString(" Turbocharge your Node.js development cycle"))
  
  createCmd  = app.Command("create", "Create a new Hot Rod app.")
  createName = createCmd.Arg("name", "The name of your app.").Required().String()

  upCmd = app.Command("up", "Beam up the source to your Hot Rod app.")
)

func main() {
  switch kingpin.MustParse(app.Parse(os.Args[1:])) {
    case createCmd.FullCommand():
      create(*createName)
    case upCmd.FullCommand():
      up()
    default:
      fmt.Printf("%s (v %s)\n\n", color.YellowString("Hot Rod"), VERSION)
      app.Usage(os.Stderr)
  }
}