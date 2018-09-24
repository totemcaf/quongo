package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"runtime"
)

type BuildInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	GoVersion string `json:"goVersion"`
}

func SysStatus(w rest.ResponseWriter, r *rest.Request) {
	v := BuildInfo{"quongo", VERSION, runtime.Version()}
	w.WriteJson(v)
}
