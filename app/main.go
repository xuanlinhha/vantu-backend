package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"vantu.org/go-backend/phrase/api"

	"vantu.org/go-backend/phrase/repo"
	"vantu.org/go-backend/phrase/service"
)

func main() {
	// env variables
	port := os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		port = ":3000"
	}
	charsPath := os.Getenv("CHARS_PATH")
	phrasesPath := os.Getenv("PHRASES_PATH")
	sqliteDB := os.Getenv("SQLITE_DB")
	timeout := os.Getenv("TIMEOUT")
	timeoutValue, err := strconv.Atoi(timeout)
	if err != nil {
		logrus.Fatal("Invalid Timeout: ", timeout)
	}
	logrus.Info("timeoutValue = ", timeoutValue)

	// read static data
	charsJsonFile, err := os.Open(charsPath)
	if err != nil {
		logrus.Fatal("Cannot Open Chars Json: ", err)
	}
	defer charsJsonFile.Close()
	byteValue, err := ioutil.ReadAll(charsJsonFile)
	if err != nil {
		logrus.Fatal("Cannot Read All Chars: ", err)
	}
	var allChars []string
	json.Unmarshal(byteValue, &allChars)

	phrasesJsonFile, err := os.Open(phrasesPath)
	if err != nil {
		logrus.Fatal("Cannot open Phrases: ", err)
	}
	defer phrasesJsonFile.Close()
	byteValue, err = ioutil.ReadAll(phrasesJsonFile)
	if err != nil {
		logrus.Fatal("Cannot Read All Phrases: ", err)
	}
	var allPhrases []string
	json.Unmarshal(byteValue, &allPhrases)

	// server
	e := echo.New()
	e.Use(middleware.CORS())

	// DB
	db, err := sqlx.Connect("sqlite3", sqliteDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialization
	phraseRepo := repo.NewPhraseRepo(db)
	timeoutContext := time.Duration(timeoutValue) * time.Second
	phraseService := service.NewPhraseService(phraseRepo, timeoutContext)
	searchService := service.NewSearchService(allChars, allPhrases)
	api.NewPhraseHandler(e, phraseService, searchService)

	if err := e.Start(port); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
