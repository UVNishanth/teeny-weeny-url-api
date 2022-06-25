package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type mapping struct {
	URL   string `json:"URL"`
	S_URL string `json:"shortURL"`
}

var db []mapping

type request struct {
	URL string `json:"URL"`
}

type serverConfig struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func createURL(w http.ResponseWriter, r *http.Request) {
	var url request
	json.NewDecoder(r.Body).Decode(&url)
	fmt.Println(url.URL)
	act_url := url.URL
	hash := md5.Sum([]byte(act_url))
	hexa_string := hex.EncodeToString(hash[:])
	s_url := hexa_string[:4]
	var entry = mapping{URL: act_url, S_URL: s_url}
	db = append(db, entry)
	extended_string := "blah.back/" + s_url
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(extended_string)
}

func getDB(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(db)
}

func resolveURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s_url, ok := vars["s_url"]
	if !ok {
		fmt.Println("no short url in link")
	}
	var actUrl request
	for _, entry := range db {
		if entry.S_URL == s_url {
			actUrl = request{URL: entry.URL}
		}
	}
	http.Redirect(w, r, "http://"+string(actUrl.URL), http.StatusSeeOther)
}

func main() {
	var cfg serverConfig
	f, err := os.Open("server_config.yaml")
	if err != nil {
		fmt.Println("File does not exist")
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println(err)
	}
	server_info := cfg.Server.Host + ":" + cfg.Server.Port
	//fmt.Println(cfg.Server.Host)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/create", createURL).Methods("POST")
	router.HandleFunc("/db", getDB).Methods("GET")
	router.HandleFunc("/{s_url}", resolveURL).Methods("GET")
	log.Fatal(http.ListenAndServe(server_info, router))
}
