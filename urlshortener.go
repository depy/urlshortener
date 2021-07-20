package main

import (
	"log"
	"net/http"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
)

const (
	requestTimeout      = 10 * time.Second
	dialTimeout         = 5 * time.Second
	redirectFallbackUrl = "https://motherfuckingwebsite.com"
)

var etcdClient *clientv3.Client
var etcdEndpoints = []string{"http://etcd:2379"}

func newEtcdClient() *clientv3.Client {
	cfg := clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: dialTimeout,
	}

	etcdclient, err := clientv3.New(cfg)

	if err != nil {
		log.Println(err)
	}

	return etcdclient
}

func etcdRedirectHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	redirectUrl := redirectFallbackUrl
	key := r.URL.Path

	response, err := etcdClient.Get(ctx, key)

	if err != nil {
		log.Println(err)
	}

	// Match the first entry it finds (currentl it's not handling duplicates)
	for _, entry := range response.Kvs {
		redirectUrl = string(entry.Value)
		break
	}

	log.Println("Redirecting for", key, "to", redirectUrl)

	w.Header().Set("Location", redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}

func main() {
	etcdClient = newEtcdClient()

	mux := http.NewServeMux()
	mux.HandleFunc("/", etcdRedirectHandler)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}
