// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"
)


 //exported field since they begins with a capital letter
type ServerValues struct {
	ServerAddress string
	ServerHostname string
	ServerDate string
	ServerUrl string
	Host string 
	ClusterName string
}

func main() {
	var addr = ":80"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if sleep := r.URL.Query().Get("sleep"); sleep != "" {
			d, err := time.ParseDuration(sleep)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			time.Sleep(d)
		}

		defer r.Body.Close()
		log.Println("{url:",r.URL,"}, cluster:", os.Getenv("CLUSTER_NAME"), "}")


		tmplt := template.New("index.html")       //create a new template with some name
		var err error;
		tmplt, err = tmplt.ParseFiles("index.html") //parse some content and generate a template, which is an internal representation
		if err != nil {
			log.Println(err)
		}

		currentTime := time.Now()
		hostname, _ := os.Hostname()

		p := ServerValues{
 			ServerAddress: r.RemoteAddr,
			ServerHostname: hostname,
			ServerDate: currentTime.Format("2006.01.02 15:04:05"),
			ServerUrl: r.RequestURI,
			Host: r.Host, 
			ClusterName: os.Getenv("CLUSTER_NAME"),
		} //define an instance with required field

		err = tmplt.Execute(w, p) //merge template ‘t’ with content of ‘p’
				if err != nil {
			log.Println(err)
		}
		
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("server started on " + addr)

	<-done
	log.Println("server stop requested")

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
	log.Print("server exit properly")
}