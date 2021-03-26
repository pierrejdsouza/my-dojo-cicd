package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Payload struct {
	Cluster          string      `json:"Cluster"`
	TaskARN          string      `json:"TaskARN"`
	Revision         string      `json:"Revision"`
	AvailabilityZone string      `json:"AvailabilityZone"`
	LaunchType       string      `json:"LaunchType"`
	Containers       []Container `json:"Containers"`
	Version          string
}

type Container struct {
	Name      string    `json:"Name"`
	Image     string    `json:"Image"`
	Limits    Limits    `json:"Limits"`
	CreatedAt string    `json:"CreatedAt"`
	StartedAt string    `json:"StartedAt"`
	Networks  []Network `json:"Networks"`
}

type Limits struct {
	CPU    int `json:"CPU"`
	Memory int `json:"Memory"`
}

type Network struct {
	NetworkMode              string   `json:"NetworkMode"`
	IPv4Addresses            []string `json:"IPv4Addresses"`
	MACAddress               string   `json:"MACAddress"`
	IPv4SubnetCIDRBlock      string   `json:"IPv4SubnetCIDRBlock"`
	PrivateDNSName           string   `json:"PrivateDNSName"`
	SubnetGatewayIpv4Address string   `json:"SubnetGatewayIpv4Address"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func jsonify(w http.ResponseWriter, r *http.Request) {
	payload := new(Payload)
	version := os.Getenv("COMMIT_ID")
	getJson(os.Getenv("ECS_CONTAINER_METADATA_URI_V4")+"/task", payload)

	profile := Payload{payload.Cluster, payload.TaskARN, payload.Revision, payload.AvailabilityZone, payload.LaunchType, payload.Containers, version}
	fmt.Println(profile)
	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(js)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/json", jsonify).Methods("GET")

	router.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		})

	corsObj := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(corsObj)(router)))
}
