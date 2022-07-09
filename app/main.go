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
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"vantu.org/go-backend/model"
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
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		logrus.Fatal("Invalid Timeout: ", err)
	}
	config := model.Config{
		Port:        port,
		CharsPath:   os.Getenv("CHARS_PATH"),
		PhrasesPath: os.Getenv("PHRASES_PATH"),
		SqliteDB:    os.Getenv("SQLITE_DB"),
		Timeout:     timeout,
	}

	// read chars path
	charsJsonFile, err := os.Open(config.CharsPath)
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

	// read phrases path
	phrasesJsonFile, err := os.Open(config.PhrasesPath)
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

	// Initialization
	db, err := sqlx.Connect("sqlite3", config.SqliteDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	phraseRepo := repo.NewPhraseRepo(db)
	timeoutContext := time.Duration(config.Timeout) * time.Second
	phraseService := service.NewPhraseService(phraseRepo, timeoutContext)
	searchService := service.NewSearchService(allChars, allPhrases)

	// server
	e := echo.New()
	e.Use(middleware.CORS())
	api.NewPhraseHandler(e, phraseService, searchService)

	// metrics
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	if err := e.Start(port); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
