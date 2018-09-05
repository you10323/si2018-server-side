package userlike

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	user_m_r := repositories.NewUserMatchRepository()
	user_l_r := repositories.NewUserLikeRepository()
	user_r := repositories.NewUserRepository()
	user_t_r := repositories.NewUserTokenRepository()
	userByToken, err := user_t_r.GetByToken(p.Token)
	UserID := userByToken.UserID
	ids, err := user_m_r.FindAllByUserID(UserID)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	likes, err := user_l_r.FindGotLikeWithLimitOffset(UserID, int(p.Limit), int(p.Offset), ids)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	var LikeUserIDs []int64
	for _, like := range likes {
		LikeUserIDs = append(LikeUserIDs, like.UserID)
	}
	LikeUsers, _ := user_r.FindByIDs(LikeUserIDs)
	// UserLikeからLikeUserResponsesへのコンバート
	var LikeUserResponses entities.LikeUserResponses
	for _, like := range likes {
		LikeUserResponse := entities.LikeUserResponse{}
		LikeUserResponse.LikedAt = like.CreatedAt
		for _, likeUser := range LikeUsers {
			if likeUser.ID == like.UserID {
				LikeUserResponse.ApplyUser(likeUser)
			}
		}
		LikeUserResponses = append(LikeUserResponses, LikeUserResponse)
	}
	sEnt := LikeUserResponses.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	return si.NewPostLikeOK()
}
