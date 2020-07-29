package main

import (
	"flag"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	r "github.com/kinderyj/k8s-scheduler-extender-demo/pkg/router"
	"k8s.io/klog"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {

	klog.InitFlags(nil)
	flag.Parse()
	router := httprouter.New()
	router.GET("/", r.Index)
	router.POST("/filter", r.Filter)
	router.POST("/prioritize", r.Prioritize)
	klog.V(3).Infof("Starting http server.")
	klog.Fatalf("Failed to setup http server, err: %v", http.ListenAndServe(":8000", router))
}
