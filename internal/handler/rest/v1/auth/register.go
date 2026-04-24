package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"asona/internal/constants"
	"asona/internal/controller/auth"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
)

// registerRequest holds the unified registration request body for all 3 steps.
type registerRequest struct {
	Step     int    `json:"step"     binding:"required,oneof=1 2 3"`
	Email    string `json:"email"    binding:"required,email"`
	OTP      string `json:"otp"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register handles user registration in 3 steps
// @Summary      User Registration
// @Description  Register a new user in the system using 3-step OTP flow
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      registerRequest  true  "Registration details"
// @Success      200      {object}  response.Response
// @Success      201      {object}  response.Response{data=loginResponse}
// @Failure      400      {object}  response.Response
// @Failure      409      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /auth/register [post]
func (h Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ERROR.Printf("[Register] failed request param: %+v", err)
		c.JSON(http.StatusBadRequest, response.NewResponse(
			constants.InvalidRequestParams.Code,
			constants.InvalidRequestParams.Message,
			nil,
		))
		return
	}

	ctx := c.Request.Context()

	switch req.Step {
	case 1:
		err := h.ctrl.RegisterStep1SendOTP(ctx, req.Email)
		if err != nil {
			logger.ERROR.Printf("[Register] step 1 failed for %s: %+v", req.Email, err)
			if errors.Is(err, auth.ErrUserAlreadyExists) {
				c.JSON(http.StatusConflict, response.NewResponse(
					constants.EmailExists.Code,
					constants.EmailExists.Message,
					nil,
				))
				return
			}
			c.JSON(http.StatusInternalServerError, response.NewResponse(
				constants.InternalServerError.Code,
				constants.InternalServerError.Message,
				nil,
			))
			return
		}
		c.JSON(http.StatusOK, response.NewResponse(
			constants.SendEmailRegisterSuccess.Code,
			constants.SendEmailRegisterSuccess.Message,
			nil,
		))
		return

	case 2:
		if req.OTP == "" {
			c.JSON(http.StatusBadRequest, response.NewResponse(
				constants.InvalidRequestParams.Code,
				"OTP is required for step 2",
				nil,
			))
			return
		}
		err := h.ctrl.RegisterStep2VerifyOTP(ctx, req.Email, req.OTP)
		if err != nil {
			logger.ERROR.Printf("[Register] step 2 failed for %s: %+v", req.Email, err)
			if errors.Is(err, auth.ErrWrongOTP) {
				c.JSON(http.StatusBadRequest, response.NewResponse(
					constants.VerifyCodeExpired.Code,
					constants.VerifyCodeExpired.Message,
					nil,
				))
				return
			}
			c.JSON(http.StatusInternalServerError, response.NewResponse(
				constants.InternalServerError.Code,
				constants.InternalServerError.Message,
				nil,
			))
			return
		}
		c.JSON(http.StatusOK, response.NewResponse(
			constants.EmailVerifiedSuccess.Code,
			constants.EmailVerifiedSuccess.Message,
			nil,
		))
		return

	case 3:
		if req.OTP == "" || req.Name == "" || req.Username == "" || len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, response.NewResponse(
				constants.InvalidRequestParams.Code,
				"Missing required fields for step 3",
				nil,
			))
			return
		}

		user, token, err := h.ctrl.RegisterStep3Complete(ctx, auth.RegisterInput{
			Name:     req.Name,
			Username: req.Username,
			Email:    req.Email,
			OTP:      req.OTP,
			Password: req.Password,
		})
		if err != nil {
			logger.ERROR.Printf("[Register] step 3 failed for %s: %+v", req.Email, err)
			if errors.Is(err, auth.ErrUserAlreadyExists) {
				c.JSON(http.StatusConflict, response.NewResponse(
					constants.EmailExists.Code,
					constants.EmailExists.Message,
					nil,
				))
				return
			}
			if errors.Is(err, auth.ErrUsernameAlreadyExists) {
				c.JSON(http.StatusConflict, response.NewResponse(
					constants.UsernameExists.Code,
					constants.UsernameExists.Message,
					nil,
				))
				return
			}
			if errors.Is(err, auth.ErrWrongOTP) {
				c.JSON(http.StatusBadRequest, response.NewResponse(
					constants.VerifyCodeExpired.Code,
					constants.VerifyCodeExpired.Message,
					nil,
				))
				return
			}
			c.JSON(http.StatusInternalServerError, response.NewResponse(
				constants.InternalServerError.Code,
				constants.InternalServerError.Message,
				nil,
			))
			return
		}

		logger.INFO.Printf("[Register] user registered successfully: %s", user.Username)
		c.JSON(http.StatusCreated, response.NewResponse(
			constants.RegisterUserSuccess.Code,
			constants.RegisterUserSuccess.Message,
			loginResponse{
				User: loginUserResponse{
					ID:       user.ID,
					Name:     user.Name,
					Username: user.Username,
					Email:    user.Email,
					Image:    user.Image,
				},
				SessionToken: token,
			},
		))
		return
	}
}
