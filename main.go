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

var tmpl *template.Template
var version string
var gitHubRepo = "https://github.com/mazk0/current-week"

func main() {
	var err error
	tmpl, err = template.ParseFiles("template.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	version = strconv.FormatInt(time.Now().Unix(), 10)

	http.Handle("/", cspHandler(gzipHandler(http.HandlerFunc(weekHandler))))
	http.Handle("/week/", cspHandler(gzipHandler(http.HandlerFunc(weekUpdateHandler))))
	http.Handle("/week/current/", cspHandler(gzipHandler(http.HandlerFunc(currentWeekUpdateHandler))))
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

	year, week, err := parseWeekYearFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	weekInfo := getWeekInfo(year, week)
	if err := tmpl.Execute(w, weekInfo); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func weekUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	year, week, err := parseWeekYearFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	weekInfo := getWeekInfo(year, week)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weekInfo); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

func currentWeekUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	year, week := now.ISOWeek()
	weekInfo := getWeekInfo(year, week)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weekInfo); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

func parseWeekYearFromRequest(r *http.Request) (int, int, error) {
	now := time.Now()
	year, week := now.ISOWeek()

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) > 2 {
		weekParam := pathParts[2]
		if weekParam != "" {
			if weekArgs, err := strconv.Atoi(weekParam); err == nil {
				week = weekArgs
			} else {
				return 0, 0, fmt.Errorf("invalid week number: %v", err)
			}
		}
	}

	numberOfWeeks := getNumberOfWeeks(year)
	if week < 1 {
		week = 1
	} else if week > numberOfWeeks {
		week = numberOfWeeks
	}

	return year, week, nil
}

func getWeekInfo(year int, week int) WeekInfo {
	firstDateOfWeek, lastDateOfWeek := getFirstAndLastDateOfWeek(year, week)

	return WeekInfo{
		Week:       week,
		FirstDate:  firstDateOfWeek.Format("2006-01-02"),
		LastDate:   lastDateOfWeek.Format("2006-01-02"),
		Version:    version,
		GitHubRepo: gitHubRepo,
	}
}
