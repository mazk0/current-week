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

	"CurrentWeek/internal/middleware"
	"CurrentWeek/internal/week"
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

	

	http.Handle("/", middleware.CspHandler(middleware.GzipHandler(http.HandlerFunc(weekHandler))))
	http.Handle("/api/previous/", middleware.CspHandler(http.HandlerFunc(previousWeekHandler)))
	http.Handle("/api/next/", middleware.CspHandler(http.HandlerFunc(nextWeekHandler)))
	http.Handle("/api/week/current/", middleware.CspHandler(http.HandlerFunc(currentWeekUpdateHandler)))
	http.Handle("/static/", http.StripPrefix("/static/", middleware.CacheHandler(middleware.GzipHandler(http.FileServer(http.Dir("./static"))))))
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

	year, weekNum, err := parseWeekYearFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	firstDateOfWeek, lastDateOfWeek := week.GetFirstAndLastDateOfWeek(year, weekNum)

	weekInfo := week.WeekInfoTemplate{
		Week:       weekNum,
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

	year, weekNum, err := parseWeekYearFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	current := week.TimeFromYearAndWeek(year, weekNum)
	current = operation(current)
	year, weekNum = current.ISOWeek()
	sendWeekInfoResponse(w, year, weekNum)
}

func currentWeekUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	year, weekNum := now.ISOWeek()
	sendWeekInfoResponse(w, year, weekNum)
}

func sendWeekInfoResponse(w http.ResponseWriter, year int, weekNum int) {
	weekInfo := getWeekInfo(year, weekNum)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weekInfo); err != nil {
		http.Error(w, fmt.Sprintf("Failed to handle the request please try again later: %v", err), http.StatusInternalServerError)
	}
}

func parseWeekYearFromRequest(r *http.Request) (int, int, error) {
	now := time.Now()
	year, weekNum := now.ISOWeek()

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
					weekNum = weekArgs
				} else {
					return 0, 0, fmt.Errorf("Invalid week number: %s", weekParam)
				}
			}
		}
	}

	return year, weekNum, nil
}

func getWeekInfo(year int, weekNum int) week.WeekInfo {
	firstDateOfWeek, lastDateOfWeek := week.GetFirstAndLastDateOfWeek(year, weekNum)
	return week.WeekInfo{
		Week:      weekNum,
		FirstDate: firstDateOfWeek.Format("2006-01-02"),
		LastDate:  lastDateOfWeek.Format("2006-01-02"),
	}
}
