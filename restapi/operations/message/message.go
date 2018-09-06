package message

import (
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
	userByToken, err := user_t_r.GetByToken(Token)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userByToken == nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	UserID := userByToken.UserID
	PartnerID := p.UserID
	Match, err := user_m_r.Get(UserID, PartnerID)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	var InsertMessage entities.UserMessage
	if Match != nil {
		InsertMessage = entities.UserMessage{
			UserID:    UserID,
			PartnerID: PartnerID,
			Message:   p.Params.Message,
			CreatedAt: strfmt.DateTime(time.Now()),
			UpdatedAt: strfmt.DateTime(time.Now()),
		}
	}
	err = user_mes_r.Create(InsertMessage)
	if err != nil {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	user_mes_r := repositories.NewUserMessageRepository()
	user_t_r := repositories.NewUserTokenRepository()
	userByToken, err := user_t_r.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userByToken == nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	UserID := userByToken.UserID
	PartnerID := p.UserID
	Messages, err := user_mes_r.GetMessages(UserID, PartnerID, int(*(p.Limit)), p.Latest, p.Oldest)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if Messages == nil {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	castedMessages := entities.UserMessages(Messages)
	sEnt := castedMessages.Build()
	return si.NewGetMessagesOK().WithPayload(sEnt)
}
