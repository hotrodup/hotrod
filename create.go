package main

import (
  "fmt"
  "os"
  "os/exec"
  "log"
  "strings"
  "github.com/fatih/color"
  "errors"
  "time"
  "regexp"
  "net/http"
  "io/ioutil"
)

const (
  INDENT = "   "
  ARROW = "-> "
)

func checkUnique(name string) error {
  fi, err := os.Stat(name)
  switch {
    case err != nil:
      return nil
    case fi.IsDir():
      return errors.New("Project already exists")
    default:
      return errors.New("Project directory can not be created")
  }
}

func checkDeps() error {
  _, err := exec.LookPath("git")
  _, err2 := exec.LookPath("gcloud")
  if err != nil || err2 != nil {
    return err
  }
  return nil
}

func checkAuth() error {
  out, err := exec.Command("gcloud", "auth", "list").CombinedOutput()
  if err != nil {
    log.Fatal(err)
  }
  if strings.Contains(string(out[:]), "No credentialed accounts") {
    return errors.New("No credentialed accounts")
  }
  return nil
}

func checkProject() (string, error) {
  out, err := exec.Command("gcloud", "config", "list", "project").CombinedOutput()
  if err != nil {
    log.Fatal(err)
  }
  if strings.Contains(string(out[:]), "(unset)") {
    return "", errors.New("Project unset")
  }
  project := strings.Trim(strings.Split(string(out[:]), "=")[1], " \n")
  return project, nil
}

func execCustom(name string, arg ...string) (string, error) {
  done := make(chan string)
  go func(){
    out, err := exec.Command(name, arg...).CombinedOutput()
    if err != nil {
      done <- ""
    }
    done <- string(out[:])
  }()
  c := time.Tick(1 * time.Second)
  fmt.Print(INDENT)
  for {
    select {
      case _ = <-c:
        fmt.Print(".")
      case output := <-done:
        fmt.Print("\n")
        if output == "" {
          return "", errors.New("Could not execute command.")
        }
        return output, nil
    }
  }
}

func findIP(input string) string {
  numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
  regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock 
  regEx := regexp.MustCompile(regexPattern)
  ips := regEx.FindAllString(input, -1)
  return ips[len(ips)-1]
}

func createInstance(name string) (string, error) {
  err := ioutil.WriteFile("hotrod-containers.yaml", []byte(containers), 0777)
  out, err := execCustom(
    "gcloud", "compute", "instances", "create", "hotrod-" + name,
    "--image", "container-vm-v20140929",
    "--image-project", "google-containers",
    "--metadata-from-file", "google-container-manifest=hotrod-containers.yaml",
    "--zone", "us-central1-a",
    "--machine-type", "f1-micro")
  os.Remove("hotrod-containers.yaml")
  if err != nil {
    return "", err
  }
  return findIP(out), nil
}

func configureFirewall() error {
  _, _ = execCustom(
    "gcloud", "compute", "firewall-rules", "create", "allow-http",
    "--description", "Incoming http allowed",
    "--allow", "tcp:80")
  _, _ = execCustom(
    "gcloud", "compute", "firewall-rules", "create", "allow-other",
    "--description", "Incoming src files allowed",
    "--allow", "tcp:8888")
  return nil
}

func waitForContainers(ip string) {
  engineUp := make(chan bool)
  go func(){
    for {
      resp, err := http.Get("http://"+ip)
      if err != nil {
        continue
      }
      if resp.StatusCode == http.StatusOK {
        engineUp <- true
        return
      }
      time.Sleep(2 * time.Second)
    }
  }()
  fuelerUp := make(chan bool)
  go func() {
    for {
      resp, err := http.Get("http://"+ip+":8888")
      if err != nil {
        continue
      }
      if resp.StatusCode == http.StatusOK {
        fuelerUp <- true
        return
      }
      time.Sleep(2 * time.Second)
    }
  }()
  c := time.Tick(2 * time.Second)
  fmt.Print(INDENT)
  containersUp := 0
  for {
    select {
      case _ = <-c:
        fmt.Print(".")
      case _ = <-fuelerUp:
        containersUp += 1
      case _ = <-engineUp:
        containersUp += 1
    }
    if containersUp == 2 {
      fmt.Print("\n")
      return
    }
  }
}

func copySource(name string, ip string) error {
  err := exec.Command(
    "git","clone", "https://github.com/hotrodup/engine", name).Run()
  if err != nil {
    return err
  }
  err = exec.Command(
    "rm", "-rf", name + "/.git").Run()
  if err != nil {
    return err
  }
  err = exec.Command(
    "rm", name + "/Dockerfile").Run()
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(name + "/.hotrod.yml", []byte(fmt.Sprintf(hotrodConfig, name, ip)), 0777)
  if err != nil {
    return err
  }

  return nil
}

func create(name string) {
  
  fmt.Println(CHECKERED_FLAG, color.YellowString("Creating new project"), color.GreenString(name))
  
  err := checkUnique(name)
  if err != nil {
    fmt.Println(INDENT, color.RedString("Can't create directory"), name)
    fmt.Println(INDENT, color.RedString("Choose a unique project name"))
    return    
  }

  err = checkDeps()
  if err != nil {
    fmt.Println(INDENT, color.RedString("Hot Rod requires `git` and `gcloud`"))
    fmt.Println(INDENT, color.RedString("Make sure both dependencies are on the $PATH"))
    return
  }
  err = checkAuth()
  if err != nil {
    fmt.Println(INDENT, color.RedString("Hot Rod requires an credentialed `gcloud` account"))
    fmt.Println(INDENT, color.RedString("Run `gcloud auth login`"))
    return
  }
  project, err := checkProject()
  if err != nil {
    fmt.Println(INDENT, color.RedString("Hot Rod requires an active `gcloud` project"))
    fmt.Println(INDENT, color.RedString("Create a project at"), "https://console.developers.google.com/project")
    fmt.Println(INDENT, color.RedString("and set it as default with `gcloud config set project <PROJECT>"))
    return
  }

  fmt.Println(ARROW, "Spinning up an instance")
  ip, err := createInstance(name)
  if err != nil {
    fmt.Println(INDENT, color.RedString("Hot Rod failed to create an instance"))
    fmt.Println(INDENT, color.RedString("Please enable billing and turn on the Compute API"))
    fmt.Println(INDENT, color.RedString("at"), fmt.Sprintf("https://console.developers.google.com/project/%s/apiui/api", project))
    return
  }

  done := make(chan bool)
  go func() {
    err := copySource(name, ip)
    if err != nil {
      done <- false
    }
    done <- true
  }()

  fmt.Println(ARROW, "Opening ports for traffic")
  err = configureFirewall()
  if err != nil {
    fmt.Println(INDENT, color.RedString("Opening ports failed"))
    fmt.Println(INDENT, color.RedString("Please try again later"))
    return
  }

  fmt.Println(ARROW, "Starting containers")
  waitForContainers(ip)

  d := <-done
  if !d {
    fmt.Println(INDENT, color.RedString("Hot Rod failed to create the source directory"))
    fmt.Println(INDENT, color.RedString("Check folder permissions"))
  }

  fmt.Println(ARROW, "Done")
  fmt.Println(RED_CAR, color.YellowString("Now `cd "+name+"` and run `hotrod up`"))

}

const containers = `
version: v1beta2
containers:
  - name: engine
    image: hotrod/engine
    ports:
      - name: http
        hostPort: 80
        containerPort: 8080
    volumeMounts:
      - name: app
        mountPath: /app
    env:
      - name: PORT
        value: 8080
  - name: fueler
    image: hotrod/fueler
    ports:
      - name: upload
        hostPort: 8888
        containerPort: 8888
    volumeMounts:
      - name: app
        mountPath: /app
volumes:
  - name: app
`

const hotrodConfig = `
name: %s
ip: %s
`
