package mock

import (
	"context"
	"net/http"

	"github.com/BarTar213/notificator/models"
	"github.com/BarTar213/notificator/senders"
)

type Notificator struct {
	SentInternalErr bool
	SentEmailErr    bool

	SentInternalStatus int
	SentEmailStatus    int
}

func (n *Notificator) SendInternal(ctx context.Context, templateName string, internal *senders.Internal) (int, *models.Response, error) {
	if n.SentInternalErr {
		return n.SentInternalStatus, nil, exampleErr
	}
	return http.StatusOK, &models.Response{}, nil
}

func (n *Notificator) SendEmail(ctx context.Context, templateName string, email *senders.Email) (int, *models.Response, error) {
	if n.SentEmailErr {
		return n.SentEmailStatus, nil, exampleErr
	}
	return http.StatusOK, &models.Response{}, nil
}
