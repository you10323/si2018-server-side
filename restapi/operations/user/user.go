package user

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	user_r := repositories.NewUserRepository()
	user_l_r := repositories.NewUserLikeRepository()
	user_t_r := repositories.NewUserTokenRepository()
	userByToken, err := user_t_r.GetByToken(p.Token)
	// メンターさんにバリデーション聞く
	userID := userByToken.UserID
	ids, err := user_l_r.FindLikeAll(userID)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	user, err := user_r.GetByUserID(userID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	gender := user.GetOppositeGender()
	ent, err := user_r.FindWithCondition(int(p.Limit), int(p.Offset), gender, ids)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	castedEnt := entities.Users(ent)
	sEnt := castedEnt.Build()

	return si.NewGetUsersOK().WithPayload(sEnt)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	user_r := repositories.NewUserRepository()

	profileEnt, err := user_r.GetByUserID(p.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	sEnt := profileEnt.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	user_r := repositories.NewUserRepository()
	user := entities.User{
		ID:             p.UserID,
		Nickname:       p.Params.Nickname,
		ImageURI:       p.Params.ImageURI,
		Tweet:          p.Params.Tweet,
		Introduction:   p.Params.Introduction,
		ResidenceState: p.Params.ResidenceState,
		HomeState:      p.Params.HomeState,
		Education:      p.Params.Education,
		Job:            p.Params.Job,
		AnnualIncome:   p.Params.AnnualIncome,
		Height:         p.Params.Height,
		BodyBuild:      p.Params.BodyBuild,
		MaritalStatus:  p.Params.MaritalStatus,
		Child:          p.Params.Child,
		WhenMarry:      p.Params.WhenMarry,
		WantChild:      p.Params.WantChild,
		Smoking:        p.Params.Smoking,
		Drinking:       p.Params.Drinking,
		Holiday:        p.Params.Holiday,
		HowToMeet:      p.Params.HowToMeet,
		CostOfDate:     p.Params.CostOfDate,
		NthChild:       p.Params.NthChild,
		Housework:      p.Params.Housework,
	}

	user_r.Update(&user)
	return si.NewPutProfileOK()
}
