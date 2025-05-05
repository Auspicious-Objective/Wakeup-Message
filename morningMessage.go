package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Configuration variables (to be set before running)
const (
	quoteAPIKEY       = "" // Set your API-Ninjas key here
	weatherDataAPIKey = "" // Set your OpenWeatherMap key here
	latitude          = ""
	longitude         = ""
	endingDate        = ""
	event             = ""
	name              = ""
)

// Messaging related info

const (
	botToken   = ""
	chatID     = ""
	messageUrl = "https://api.telegram.org/bot" + botToken + "/sendMessage"
)

// Quote API Response Structure
type quoteResponse []struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

// Weather API Response Structure
type weatherData struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func getQuote(APIKey string) (string, error) {
	if APIKey == "" {
		return "", fmt.Errorf("quote API key not configured")
	}

	url := "https://api.api-ninjas.com/v1/quotes" // Added category parameter

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("X-Api-Key", APIKey)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var quotes quoteResponse
	if err := json.Unmarshal(body, &quotes); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	if len(quotes) == 0 {
		return "", fmt.Errorf("no quotes received")
	}

	return fmt.Sprintf("%s once said: \"%s\"", quotes[0].Author, quotes[0].Quote), nil
}

func correctWeatherMessage(description string) string {
	if description == "clear sky" {
		description = "clear skies  🌞"
	} else if description == "shower rain" {
		description = "rain showers"
	} else if description == "thunderstorm" {
		description = "thunderstorms"
	}

	return description
}

// make the if statement to correct the grammar
func getWeatherInfo(APIKey string) (string, error) {
	if APIKey == "" {
		return "", fmt.Errorf("weather API key not configured")
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=imperial&appid=%s",
		latitude, longitude, APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var data weatherData
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	if len(data.Weather) == 0 {
		return "", fmt.Errorf("no weather data available")
	}

	var weatherDescription string = correctWeatherMessage(data.Weather[0].Description)
	return fmt.Sprintf("Today is %.1f°F with %s.", data.Main.Temp, weatherDescription), nil
}

func getTodaysDate() string {
	now := time.Now()
	day := now.Format("Monday")
	month := now.Format("January")
	dayOfMonth := now.Day() // This returns an int (1-31)

	var dayWithSuffix string
	switch dayOfMonth {
	case 1, 21, 31:
		dayWithSuffix = fmt.Sprintf("%dst", dayOfMonth)
	case 2, 22:
		dayWithSuffix = fmt.Sprintf("%dnd", dayOfMonth)
	case 3, 23:
		dayWithSuffix = fmt.Sprintf("%drd", dayOfMonth)
	default:
		dayWithSuffix = fmt.Sprintf("%dth", dayOfMonth)
	}

	return fmt.Sprintf("Today is %s, %s %s.", day, month, dayWithSuffix)
}

func getDateDifference(endDate string, event string) string {
	layout := "2006-01-02"
	end, err := time.Parse(layout, endDate)
	if err != nil {
		return fmt.Sprintf("Error parsing date: %v", err)
	}

	now := time.Now()
	daysLeft := int(end.Sub(now).Hours() / 24)

	if daysLeft < 0 {
		return "The date of" + event + "has already passed."
	} else if daysLeft == 0 {
		return "Today is the day!"
	} else {
		return fmt.Sprintf("There are %d days left until %s.", daysLeft, event)
	}
}

func compileMessage(quote string, weather string, date string, daysTill string) []string {
	var messages []string
	var messagePartOne string = "Good morning 👋 " + name + ". " + date + weather
	messages = append(messages, messagePartOne)
	var messagePartTwo string = quote + ". " + daysTill
	messages = append(messages, messagePartTwo)
	return messages
}

func sendMessage(message string, id string, url string) {
	client := &http.Client{}
	values := map[string]string{"text": message, "chat_id": id}
	jsonParameters, _ := json.Marshal(values)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonParameters))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.Status)
		defer res.Body.Close()
	}
}

func main() {
	// Get and display quote
	quote, err := getQuote(quoteAPIKEY)
	if err != nil {
		fmt.Printf("Error getting quote: %v\n", err)
	}

	// Get and display weather
	weather, err := getWeatherInfo(weatherDataAPIKey)
	if err != nil {
		fmt.Printf("Error getting weather: %v\n", err)
	}

	var date string = getTodaysDate()
	var dateDif string = getDateDifference(endingDate, event)

	var messages []string = compileMessage(quote, weather, date, dateDif)

	for message := range len(messages) {
		sendMessage(messages[message], chatID, messageUrl)
	}

}
