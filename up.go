package main

import (
  "log"
  "os"
  "bytes"
  "io"
  "mime/multipart"
  "net/http"
  "gopkg.in/fsnotify.v1"
  fp "path/filepath"
)

func post(path string, add, isFile bool, baseURL string) {

  var b bytes.Buffer
  w := multipart.NewWriter(&b)

  cwd, _ := os.Getwd()
  relPath, err := fp.Rel(cwd, path)
  if err != nil {
    log.Fatal(err)
    return
  }
  _ = w.WriteField("path", relPath)

  if (add && isFile) {
    // Add your image file
    f, err := os.Open(path)
    if err != nil {
        return
    }
    fw, err := w.CreateFormFile("file", path)
    if err != nil {
        return
    }
    if _, err = io.Copy(fw, f); err != nil {
        return
    }    
  }

  w.Close()

  var url string
  switch {
    case !add:
      url = baseURL + "/remove"
    case add && isFile:
      url = baseURL + "/addFile"
    case add && !isFile:
      url = baseURL + "/addFolder"
  }

  req, err := http.NewRequest("POST", url, &b)
  if err != nil {
      return
  }
  // Don't forget to set the content type, this will contain the boundary.
  req.Header.Set("Content-Type", w.FormDataContentType())

  // Submit the request
  client := &http.Client{}
  _, err = client.Do(req)
  return

}

func handle(event fsnotify.Event, watcher *fsnotify.Watcher, baseURL string) {
  filepath := event.Name

  // event type
  var add bool
  switch {
    case event.Op&fsnotify.Write == fsnotify.Write:
      add = true
    case event.Op&fsnotify.Create == fsnotify.Create:
      add = true
    case event.Op&fsnotify.Remove == fsnotify.Remove:
      add = false
    case event.Op&fsnotify.Rename == fsnotify.Rename:
      add = false
    default:
      return
  }

  if add {

    fi, err := os.Stat(filepath);
    switch {
      case err != nil:
        return
      case fi.IsDir():
        post(filepath, true, false, baseURL)
        err = watcher.Add(filepath)
        if err != nil {
          return
        }
      default:
        post(filepath, true, true, baseURL)
    }

  } else {
    post(filepath, false, false, baseURL)
  }
}

func start_up(baseURL string) {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    log.Fatal(err)
  }
  defer watcher.Close()

  events := make(chan fsnotify.Event)
  go func() {
    for {
      event := <-events
      sem := make(chan int)
      go func() {
        handle(event, watcher, baseURL)
        sem <- 1
      }()
      <-sem
    }
  }()

  done := make(chan bool)
  go func() {
    for {
      select {
      case event := <-watcher.Events:
        events <- event
      case err := <-watcher.Errors:
        log.Println("error:", err)
      }
    }
  }()

  cwd, _ := os.Getwd()
  err = fp.Walk(cwd, func(path string, info os.FileInfo, _ error) error {
      if info.IsDir() {
        watcher.Add(path)
      }
      return nil
    })
  err = watcher.Add(cwd)
  if err != nil {
    log.Fatal(err)
  }
  <-done
}
