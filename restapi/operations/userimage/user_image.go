package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	user_t_r := repositories.NewUserTokenRepository()
	user_i_r := repositories.NewUserImageRepository()
	assetsPath := os.Getenv("ASSETS_PATH")
	userByToken, err := user_t_r.GetByToken(p.Params.Token)
	if err != nil {
		si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userByToken == nil {
		si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	imgName := strconv.Itoa(int(userByToken.UserID)) + "_image.png"
	imgPath := assetsPath + imgName
	file, _ := os.Create(imgPath)
	defer file.Close()
	file.Write(p.Params.Image)
	userImage, err := user_i_r.GetByUserID(userByToken.UserID)
	if err != nil {
		si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userImage == nil {
		si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	userImage.Path = imgPath
	err = user_i_r.Update(*userImage)
	if err != nil {
		si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(imgPath),
		})
}
