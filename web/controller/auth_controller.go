package controller

import "codeberg.org/rchan/hmn/business"

type AuthController struct {
	b business.BusinessLayer
}

func NewAuthController(b business.BusinessLayer) *AuthController {
	return &AuthController{
		b: b,
	}
}
