package message

import (
	"fmt"
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	user_t_r := repositories.NewUserTokenRepository()
	user_mes_r := repositories.NewUserMessageRepository()
	user_m_r := repositories.NewUserMatchRepository()
	Token := p.Params.Token
	userByToken, _ := user_t_r.GetByToken(Token)
	UserID := userByToken.UserID
	PartnerID := p.UserID
	Match, _ := user_m_r.Get(UserID, PartnerID)
	if Match != nil {
		InsertMessage := entities.UserMessage{
			UserID:    UserID,
			PartnerID: PartnerID,
			Message:   p.Params.Message,
			CreatedAt: strfmt.DateTime(time.Now()),
			UpdatedAt: strfmt.DateTime(time.Now()),
		}
		user_mes_r.Create(InsertMessage)
	}
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
