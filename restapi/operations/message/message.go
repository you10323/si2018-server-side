package message

import (
	"fmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	return si.NewPostMessageOK()
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	user_mes_r := repositories.NewUserMessageRepository()
	user_t_r := repositories.NewUserTokenRepository()
	userByToken, _ := user_t_r.GetByToken(p.Token)
	UserID := userByToken.UserID
	PartnerID := p.UserID
	Messages, _ := user_mes_r.GetMessages(UserID, PartnerID, int(*(p.Limit)), p.Latest, p.Oldest)
	fmt.Println(Messages)
	castedMessages := entities.UserMessages(Messages)
	sEnt := castedMessages.Build()
	return si.NewGetMessagesOK().WithPayload(sEnt)
}
