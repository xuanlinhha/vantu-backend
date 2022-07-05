package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"vantu.org/go-backend/model"
)

type PhraseRepo interface {
	GetPhraseByHan(ctx context.Context, han string) (*model.Phrase, error)
	GetAllPhrasesInHans(ctx context.Context, hans []string) ([]*model.Phrase, error)
}

type phraseRepo struct {
	db *sqlx.DB
}

func NewPhraseRepo(db *sqlx.DB) PhraseRepo {
	return &phraseRepo{db: db}
}

func (ps *phraseRepo) GetPhraseByHan(ctx context.Context, han string) (*model.Phrase, error) {
	query := "SELECT * FROM phrases WHERE han = ?"
	list, err := ps.fetch(ctx, query, han)
	if err != nil {
		return &model.Phrase{}, err
	}
	if len(list) > 0 {
		return list[0], nil
	} else {
		return &model.Phrase{}, model.ErrNotFound
	}
}

func (ps *phraseRepo) GetAllPhrasesInHans(ctx context.Context, hans []string) ([]*model.Phrase, error) {
	query, args, err := sqlx.In("SELECT * FROM phrases WHERE han IN (?);", hans)
	query = ps.db.Rebind(query)

	list, err := ps.fetch(ctx, query, args...)
	if err != nil {
		return make([]*model.Phrase, 0), err
	}
	return list, nil
}

func (pr *phraseRepo) fetch(ctx context.Context, query string, args ...interface{}) (result []*model.Phrase, err error) {
	rows, err := pr.db.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]*model.Phrase, 0)
	for rows.Next() {
		t := model.Phrase{}
		err = rows.Scan(
			&t.Id,
			&t.Han,
			&t.Content,
			&t.Info,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, &t)
	}

	return result, nil
}
