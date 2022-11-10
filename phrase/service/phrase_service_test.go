package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"vantu.org/go-backend/model"
	"vantu.org/go-backend/phrase/repo"
)

func Test_phraseService_GetPhrase(t *testing.T) {
	type fields struct {
		phraseRepo     repo.PhraseRepo
		contextTimeout time.Duration
	}
	type args struct {
		c  context.Context
		ph string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.PhraseJson
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &phraseService{
				phraseRepo:     tt.fields.phraseRepo,
				contextTimeout: tt.fields.contextTimeout,
			}
			got, err := ps.GetPhrase(tt.args.c, tt.args.ph)
			if (err != nil) != tt.wantErr {
				t.Errorf("phraseService.GetPhrase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("phraseService.GetPhrase() = %v, want %v", got, tt.want)
			}
		})
	}
}
