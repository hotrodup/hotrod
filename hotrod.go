package main

import (
  "fmt"
  "os"
  "gopkg.in/alecthomas/kingpin.v1"
)

var (
  app = kingpin.New("hotrod", "Turbocharge your Node.js development cycle")
  
  create     = app.Command("create", "Create a new Hot Rod app.")
  createName = create.Arg("name", "The name of your app.").Required().String()

  up = app.Command("up", "Beam up the source to your Hot Rod app.")

)

func main() {
  kingpin.Version("0.0.1")

  switch kingpin.MustParse(app.Parse(os.Args[1:])) {
    case create.FullCommand():
      fmt.Printf("Creating Hot Rod app %s\n", *createName)
    case up.FullCommand():
      _, err := os.Stat("static")
      if err != nil {
        fmt.Printf("This command must be run inside your app directory\n")
        return
      }
      start_up("http://107.178.216.59:8888")
    default:
      app.Usage(os.Stderr)
  }
}