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

	

	http.Handle("/", cspHandler(gzipHandler(http.HandlerFunc(weekHandler))))
	http.Handle("/api/previous/", cspHandler(http.HandlerFunc(previousWeekHandler)))
	http.Handle("/api/next/", cspHandler(http.HandlerFunc(nextWeekHandler)))
	http.Handle("/api/week/current/", cspHandler(http.HandlerFunc(currentWeekUpdateHandler)))
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

	firstDateOfWeek, lastDateOfWeek := getFirstAndLastDateOfWeek(year, week)

	weekInfo := WeekInfoTemplate{
		Week:       week,
		FirstDate:  firstDateOfWeek.Format("2006-01-02"),
		LastDate:   lastDateOfWeek.Format("2006-01-02"),
		Version:    version,
		GitHubRepo: gitHubRepo,
	}

	if err := tmpl.Execute(w, weekInfo); err != nil {
		http.Error(w, fmt.Sprint("Failed to handle the request please try again later."), http.StatusInternalServerError)
	}
}

func previousWeekHandler(w http.ResponseWriter, r *http.Request) {
	handleWeekRequest(w, r, func(date time.Time) time.Time {
		return date.AddDate(0, 0, -7)
	})
}

func nextWeekHandler(w http.ResponseWriter, r *http.Request) {
	handleWeekRequest(w, r, func(date time.Time) time.Time {
		return date.AddDate(0, 0, 7)
	})
}

func handleWeekRequest(w http.ResponseWriter, r *http.Request, operation func(time.Time) time.Time) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	year, week, err := parseWeekYearFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	current := timeFromYearAndWeek(year, week)
	current = operation(current)
	year, week = current.ISOWeek()
	sendWeekInfoResponse(w, year, week)
}

func currentWeekUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	year, week := now.ISOWeek()
	sendWeekInfoResponse(w, year, week)
}

func sendWeekInfoResponse(w http.ResponseWriter, year int, week int) {
	weekInfo := getWeekInfo(year, week)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weekInfo); err != nil {
		http.Error(w, fmt.Sprintf("Failed to handle the request please try again later: %v", err), http.StatusInternalServerError)
	}
}

func parseWeekYearFromRequest(r *http.Request) (int, int, error) {
	now := time.Now()
	year, week := now.ISOWeek()

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) > 2 {
		yearIndex := indexOf(pathParts, "year")
		if yearIndex > -1 && yearIndex+1 < len(pathParts) {
			yearParam := pathParts[yearIndex+1]
			if yearParam != "" {
				if yearArgs, err := strconv.Atoi(yearParam); err == nil {
					year = yearArgs
				} else {
					return 0, 0, fmt.Errorf("Invalid year number: %s", yearParam)
				}
			}
		}
		weekIndex := indexOf(pathParts, "week")
		if weekIndex > -1 && weekIndex+1 < len(pathParts) {
			weekParam := pathParts[weekIndex+1]
			if weekParam != "" {
				if weekArgs, err := strconv.Atoi(weekParam); err == nil {
					week = weekArgs
				} else {
					return 0, 0, fmt.Errorf("Invalid week number: %s", weekParam)
				}
			}
		}
	}

	return year, week, nil
}

func getWeekInfo(year int, week int) WeekInfo {
	firstDateOfWeek, lastDateOfWeek := getFirstAndLastDateOfWeek(year, week)
	return WeekInfo{
		Week:      week,
		FirstDate: firstDateOfWeek.Format("2006-01-02"),
		LastDate:  lastDateOfWeek.Format("2006-01-02"),
	}
}
