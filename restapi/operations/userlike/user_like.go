package userlike

import (
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	user_m_r := repositories.NewUserMatchRepository()
	user_l_r := repositories.NewUserLikeRepository()
	user_r := repositories.NewUserRepository()
	user_t_r := repositories.NewUserTokenRepository()
	userByToken, err := user_t_r.GetByToken(p.Token)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userByToken == nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
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
	if likes == nil {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	var LikeUserIDs []int64
	for _, like := range likes {
		LikeUserIDs = append(LikeUserIDs, like.UserID)
	}
	LikeUsers, err := user_r.FindByIDs(LikeUserIDs)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

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
	user_t_r := repositories.NewUserTokenRepository()
	user_l_r := repositories.NewUserLikeRepository()
	user_m_r := repositories.NewUserMatchRepository()
	Token := p.Params.Token
	userByToken, err := user_t_r.GetByToken(Token)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userByToken == nil {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	UserID := userByToken.UserID
	PartnerID := p.UserID
	if UserID == PartnerID {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	//いいねする相手をいいねしていないか
	Like, err := user_l_r.GetLikeBySenderIDReceiverID(UserID, PartnerID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if Like != nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	InsertLike := entities.UserLike{
		UserID:    UserID,
		PartnerID: PartnerID,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}
	err = user_l_r.Create(InsertLike)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	//いいねした相手にいいねされているかどうか
	Like, err = user_l_r.GetLikeBySenderIDReceiverID(PartnerID, UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if Like != nil {
		InsertMatch := entities.UserMatch{
			UserID:    PartnerID,
			PartnerID: UserID,
			CreatedAt: strfmt.DateTime(time.Now()),
			UpdatedAt: strfmt.DateTime(time.Now()),
		}
		err = user_m_r.Create(InsertMatch)
	}
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:    "200",
			Message: "OK",
		})
}
