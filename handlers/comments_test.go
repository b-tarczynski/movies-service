package handlers

import (
	"log"
	"reflect"
	"testing"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/mock"
	"github.com/BarTar213/movies-service/storage"
)

func TestNewCommentHandlers(t *testing.T) {
	type args struct {
		storage storage.Storage
		logger  *log.Logger
		headers *config.Headers
	}
	tests := []struct {
		name string
		args args
		want *CommentHandlers
	}{
		{
			name: "positiveNewCommentHandlers",
			args: args{
				storage: &mock.Storage{},
				logger:  &log.Logger{},
			},
			want: &CommentHandlers{
				storage: &mock.Storage{},
				logger:  &log.Logger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommentHandlers(tt.args.storage, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommentHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
