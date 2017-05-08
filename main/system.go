package main

import (
  "net/http"
  "runtime"
)

type BuildInfo struct {
  Name      string  `json:"name"`
  Version   string  `json:"version"`
  GoVersion string  `json:"goVersion"`
}

func SysStatus(w http.ResponseWriter, r *http.Request) {
  v := BuildInfo{"quongo", VERSION, runtime.Version() }

  Ok(w, v)
}

