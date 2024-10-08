package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	Groupie_tracker "groupie_tracker/handlers"
)

func main() {
	port := ":8090"
	http.HandleFunc("/", Groupie_tracker.GetDataFromJson)
	http.HandleFunc("/Artist/{id}", Groupie_tracker.HandlerShowRelation)
	http.HandleFunc("/geoMap", Groupie_tracker.GeoMap)
	http.HandleFunc("/styles/", Groupie_tracker.HandleStyle)
	fmt.Printf("http://localhost%s", port)
	fmt.Println()
	// Openbrowser("http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func Openbrowser(zz string) {
	time.Sleep(time.Second)
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", zz).Start()
	case "linux":
		err = exec.Command("xdg-open", zz).Start()
	}
	if err != nil {
		log.Fatal(err)
	}
}
