package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	user_m_r := repositories.NewUserMatchRepository()
	user_r := repositories.NewUserRepository()
	user_t_r := repositories.NewUserTokenRepository()
	userByToken, err := user_t_r.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userByToken == nil {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	UserID := userByToken.UserID
	Matches, err := user_m_r.FindByUserIDWithLimitOffset(UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if Matches == nil {
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	var MatchUserIDs []int64
	for _, Match := range Matches {
		MatchUserIDs = append(MatchUserIDs, Match.PartnerID)
	}
	MatchUsers, err := user_r.FindByIDs(MatchUserIDs)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	var MatchUserResponses entities.MatchUserResponses
	for _, match := range Matches {
		MatchUserResponse := entities.MatchUserResponse{}
		MatchUserResponse.MatchedAt = match.CreatedAt
		for _, MatchUser := range MatchUsers {
			if MatchUser.ID == match.PartnerID {
				MatchUserResponse.ApplyUser(MatchUser)
			}
		}
		MatchUserResponses = append(MatchUserResponses, MatchUserResponse)
	}
	sEnt := MatchUserResponses.Build()
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
