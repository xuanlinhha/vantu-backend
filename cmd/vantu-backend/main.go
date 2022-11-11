package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"vantu.org/go-backend/common"
	"vantu.org/go-backend/phrase/api"

	"vantu.org/go-backend/phrase/repo"
	"vantu.org/go-backend/phrase/service"
)

func main() {
	common.InitLogger()
	defer common.CleanLogger()
	common.InitConfig()

	// read chars path
	charsJsonFile, err := os.Open(common.Conf.CharsPath)
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
	phrasesJsonFile, err := os.Open(common.Conf.PhrasesPath)
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
	db, err := sqlx.Connect("sqlite3", common.Conf.SqliteDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	phraseRepo := repo.NewPhraseRepo(db)
	timeoutContext := time.Duration(5) * time.Second
	phraseService := service.NewPhraseService(phraseRepo, timeoutContext)
	searchService := service.NewSearchService(allChars, allPhrases)

	// server
	e := echo.New()
	e.Use(middleware.CORS())
	api.NewPhraseHandler(e, phraseService, searchService)

	// metrics
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	if err := e.Start(common.Conf.Address); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
