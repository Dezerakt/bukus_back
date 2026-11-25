package v1

import (
	"bukus/internal/controller/restapi/v1/request"
	"bukus/internal/controller/restapi/v1/response"
	"bukus/internal/usecase/auth"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func (obj *V1) register(ctx fiber.Ctx) error {
	var body request.Register

	if err := ctx.Bind().Body(&body); err != nil {
		obj.l.Error(err, "restapi - v1 - register - bind")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := obj.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "validation failed: "+err.Error())
	}

	err := obj.authUC.Register(ctx.RequestCtx(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return errorResponse(ctx, http.StatusConflict, "user already exists")
		}

		obj.l.Error(err, "restapi - v1 - register")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to register")
	}

	return ctx.Status(http.StatusCreated).JSON(nil)
}

func (obj *V1) login(ctx fiber.Ctx) error {
	var body request.Login

	if err := ctx.Bind().Body(&body); err != nil {
		obj.l.Error(err, "restapi - v1 - login - bind")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := obj.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "validation failed: "+err.Error())
	}

	accessToken, refreshToken, err := obj.authUC.Login(ctx.RequestCtx(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return errorResponse(ctx, http.StatusUnauthorized, "invalid email or password")
		}

		obj.l.Error(err, "restapi - v1 - login")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to login")
	}

	return ctx.Status(http.StatusOK).JSON(response.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (obj *V1) refresh(ctx fiber.Ctx) error {
	var body request.Refresh

	if err := ctx.Bind().Body(&body); err != nil {
		obj.l.Error(err, "restapi - v1 - refresh - bind")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := obj.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "validation failed: "+err.Error())
	}

	accessToken, refreshToken, err := obj.authUC.Refresh(ctx.RequestCtx(), body.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			return errorResponse(ctx, http.StatusUnauthorized, "invalid refresh token")
		}

		obj.l.Error(err, "restapi - v1 - refresh")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to refresh token")
	}

	return ctx.Status(http.StatusOK).JSON(response.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
