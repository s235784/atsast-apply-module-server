package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/njupt-sast/atsast-apply-module-server/common/jwt"
	"github.com/njupt-sast/atsast-apply-module-server/controller/response"
	"github.com/njupt-sast/atsast-apply-module-server/model/entity"
	"github.com/njupt-sast/atsast-apply-module-server/service"
)

type ReadUserProfileRequest struct {
	RequesterId *uuid.UUID `json:"requesterId"`
	UserId      *uuid.UUID `json:"userId"`
}

func ReadUserProfileRequestParser(c *gin.Context) (*ReadUserProfileRequest, error) {
	request := ReadUserProfileRequest{}

	identity := jwt.MustExtractIdentity(c)
	userId := uuid.MustParse(c.Param("userId"))

	request.RequesterId = identity.Uid
	request.UserId = &userId
	return &request, nil
}

type ReadUserProfileResponse struct {
	UserId  *uuid.UUID          `json:"userId" binding:"required"`
	Profile *entity.UserProfile `json:"profile"`
}

func ReadUserProfileRequestHandler(request *ReadUserProfileRequest) *response.Response {
	if *request.RequesterId != *request.UserId {
		isAdmin, err := service.IsAdmin(request.RequesterId)
		if err != nil {
			return response.Failed().SetMsg(err.Error())
		}
		if !isAdmin {
			return response.Failed().SetMsg(service.PermissionDeniedErr.Error())
		}
	}

	user, err := service.ReadUser(request.UserId)
	if err != nil {
		return response.Failed().SetMsg(err.Error())
	}
	return response.Success().SetData(ReadUserProfileResponse{
		UserId:  user.UserId,
		Profile: user.Profile,
	})
}
