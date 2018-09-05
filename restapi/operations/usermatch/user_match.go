package usermatch

import (
	"fmt"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	user_m_r := repositories.NewUserMatchRepository()
	//user_r := repositories.NewUserRepository()
	user_t_r := repositories.NewUserTokenRepository()
	userByToken, err := user_t_r.GetByToken(p.Token)
	UserID := userByToken.UserID
	MatchUsers, err := user_m_r.FindByUserIDWithLimitOffset(UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	fmt.Println(MatchUsers)
	return si.NewGetMatchesOK().WithPayLoad(sEnt)
}
