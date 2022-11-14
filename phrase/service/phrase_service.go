package service

import (
	"context"
	"time"

	"vantu.org/go-backend/model"
	"vantu.org/go-backend/phrase/repo"
)

type PhraseService interface {
	GetPhrase(c context.Context, ph string) (*model.PhraseJson, error)
	GetAllPhrases(c context.Context, phs []string) ([]*model.PhraseJson, error)
}

type phraseService struct {
	phraseRepo     repo.PhraseRepo
	contextTimeout time.Duration
}

func NewPhraseService(repo repo.PhraseRepo, timeout time.Duration) PhraseService {
	return &phraseService{phraseRepo: repo, contextTimeout: timeout}
}

func (ps *phraseService) GetPhrase(c context.Context, ph string) (*model.PhraseJson, error) {
	ctx, cancel := context.WithTimeout(c, ps.contextTimeout)
	defer cancel()

	phrase, err := ps.phraseRepo.GetPhraseByHan(ctx, ph)
	if err == nil {
		return model.ConvertToJson(phrase), nil
	}
	return &model.PhraseJson{}, err
}

func (ps *phraseService) GetAllPhrases(c context.Context, phs []string) ([]*model.PhraseJson, error) {
	if len(phs) == 0 {
		return []*model.PhraseJson{}, nil
	}
	ctx, cancel := context.WithTimeout(c, ps.contextTimeout)
	defer cancel()

	phrases, err := ps.phraseRepo.GetAllPhrasesInHans(ctx, phs)
	if err == nil {
		return convertToPhraseJson(phrases), nil
	}
	return make([]*model.PhraseJson, 0), err
}

func convertToPhraseJson(phrases []*model.Phrase) []*model.PhraseJson {
	phraseJsons := make([]*model.PhraseJson, len(phrases))
	if phrases == nil {
		return phraseJsons
	}
	// convert
	for i, pj := range phrases {
		phraseJsons[i] = model.ConvertToJson(pj)
	}
	return phraseJsons
}
