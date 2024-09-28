package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.Handle("/", cspHandler(gzipHandler(http.HandlerFunc(weekHandler))))
	http.Handle("/week/", cspHandler(gzipHandler(http.HandlerFunc(weekUpdateHandler))))
	http.Handle("/static/", http.StripPrefix("/static/", cacheHandler(gzipHandler(http.FileServer(http.Dir("./static"))))))
	http.Handle("/robots.txt", http.FileServer(http.Dir(".")))

	port := "8080"
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

func weekHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	weekInfo, err := getWeekInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	template, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, weekInfo); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func weekUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	weekInfo, err := getWeekInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weekInfo); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

func getWeekInfo(r *http.Request) (WeekInfo, error) {
	now := time.Now()
	year, week := now.ISOWeek()

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) > 2 {
		weekParam := pathParts[2]
		if weekParam != "" {
			if weekArgs, err := strconv.Atoi(weekParam); err == nil {
				week = weekArgs
			} else {
				return WeekInfo{}, fmt.Errorf("invalid week number: %v", err)
			}
		}
	}

	firstDateOfWeek, lastDateOfWeek := getFirstAndLastDateOfWeek(year, week)

	return WeekInfo{
		Week:      week,
		FirstDate: firstDateOfWeek.Format("2006-01-02"),
		LastDate:  lastDateOfWeek.Format("2006-01-02"),
	}, nil
}
