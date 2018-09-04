package message

import (
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	return si.NewPostMessageOK()
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	user_m_r := repositories.NewUserMatchRepository()
	user_mes_r := repositories.NewUserMessageRepository()
	user_t_r := repositories.NewUserTokenRepository()
	userByToken, _ := user_t_r.GetByToken(p.Token)
	UserID := userByToken.UserID
	MatchingUserIDs, _ := user_m_r.FindAllByUserID(UserID)
	var MessageAll []*models.UserMessage
	for key, value := range MatchingUserIDs {
		Messages, _ := user_mes_r.GetMessages(UserID, value, p.Limit, p.Latest, p.Oldest)
		MessageAll = append(MessageAll, Messages)
	}
	return si.NewGetMessagesOK().WithPayload(MessageAll)
}
