package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"vantu.org/go-backend/model"
	"vantu.org/go-backend/phrase/service"
)

type phraseHandler struct {
	PhraseService service.PhraseService
	SearchService service.SearchService
}

func NewPhraseHandler(e *echo.Echo, phraseSvc service.PhraseService, searchSvc service.SearchService) {
	handler := &phraseHandler{
		PhraseService: phraseSvc,
		SearchService: searchSvc,
	}
	phrases := e.Group("/phrases")
	phrases.GET("", handler.GetPhrase)
	phrases.POST("/containing-first-char", handler.GetAllPhrasesContainingFirstChar)
	phrases.POST("/in-text", handler.GetAllPhrasesInText)
}

func (ph *phraseHandler) GetPhrase(c echo.Context) error {
	param := c.QueryParam("han")
	// validate
	if strings.TrimSpace(param) == "" {
		return c.JSON(http.StatusBadRequest, model.ErrBadParamInput.Error())
	}

	ctx := c.Request().Context()
	phrase, err := ph.PhraseService.GetPhrase(ctx, param)
	if err != nil {
		return c.JSON(model.GetHttpStatus(err), err.Error())
	}
	return c.JSON(http.StatusOK, model.ResponseData{Code: model.OK_STATUS, Message: model.OK_STATUS, Data: phrase})
}

func (ph *phraseHandler) GetAllPhrasesContainingFirstChar(c echo.Context) (err error) {
	var reqData model.RequestData
	err = c.Bind(&reqData)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	// validate
	if strings.TrimSpace(reqData.Text) == "" {
		return c.JSON(http.StatusBadRequest, model.ErrBadParamInput.Error())
	}

	ctx := c.Request().Context()
	// search
	ll, err := ph.SearchService.SearchAllPhrasesContainingFirstChar(ctx, reqData.Text)
	if err != nil {
		return c.JSON(model.GetHttpStatus(err), err.Error())
	}

	// get data
	phrases := make([]string, 0)
	for _, l := range ll {
		phrases = append(phrases, l...)
	}
	phraseJsons, err := ph.PhraseService.GetAllPhrases(ctx, phrases)
	if err != nil {
		return c.JSON(model.GetHttpStatus(err), err.Error())
	}

	// han -> JsonData
	phraseData := map[string]*model.PhraseJson{}
	for _, pj := range phraseJsons {
		phraseData[pj.Han] = pj
	}

	// result
	result := make([][]*model.PhraseJson, len(ll))
	for i, l := range ll {
		tmp := make([]*model.PhraseJson, len(l))
		for j, p := range l {
			tmp[j] = phraseData[p]
		}
		result[i] = tmp
	}
	return c.JSON(http.StatusOK, model.ResponseData{Code: model.OK_STATUS, Message: model.OK_STATUS, Data: result})
}

func (ph *phraseHandler) GetAllPhrasesInText(c echo.Context) (err error) {
	var reqData model.RequestData
	err = c.Bind(&reqData)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	// validate
	if strings.TrimSpace(reqData.Text) == "" {
		err := model.ErrBadParamInput
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	// search
	phrases, err := ph.SearchService.SearchAllPhrasesInInput(ctx, reqData.Text)
	if err != nil {
		return c.JSON(model.GetHttpStatus(err), err.Error())
	}

	// get data
	phraseJsons, err := ph.PhraseService.GetAllPhrases(ctx, phrases)
	if err != nil {
		return c.JSON(model.GetHttpStatus(err), err.Error())
	}

	// han -> JsonData
	phraseData := map[string]*model.PhraseJson{}
	for _, pj := range phraseJsons {
		phraseData[pj.Han] = pj
	}

	// result
	result := make([]*model.PhraseJson, len(phrases))
	for i, p := range phrases {
		pj, _ := phraseData[p]
		result[i] = pj
	}
	return c.JSON(http.StatusOK, model.ResponseData{Code: model.OK_STATUS, Message: model.OK_STATUS, Data: result})
}
