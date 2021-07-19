package main

import (
	"log"
	"net/http"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
)

const (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

func main() {
	etcdEndpoints := []string{"http://etcd:2379"}

	cfg := clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	}

	etcdclient, err := clientv3.New(cfg)

	if err != nil {
		log.Println(err)
	}

	cli, _ := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"etcd:2379"},
	})
	defer cli.Close()

	kv := clientv3.NewKV(cli)
	defer etcdclient.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
		redirectUrl := "https://motherfuckingwebsite.com" // default
		key := r.URL.Path
		response, err := kv.Get(ctx, key)

		if err != nil {
			log.Println(err)
		}

		for _, entry := range response.Kvs {
			redirectUrl = string(entry.Value)
			break
		}

		log.Println("Redirecting for", key, "to", redirectUrl)

		w.Header().Set("Location", redirectUrl)
		w.WriteHeader(http.StatusMovedPermanently)
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}
