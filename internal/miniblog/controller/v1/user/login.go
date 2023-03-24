// Copyright 2023 lawff. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/gin-gonic/gin"

	"github.com/lawff/miniblog/internal/pkg/core"
	"github.com/lawff/miniblog/internal/pkg/errno"
	"github.com/lawff/miniblog/internal/pkg/log"
	v1 "github.com/lawff/miniblog/pkg/api/miniblog/v1"
)

// 登录 miniblog 并返回一个 JWT Token.
func (ctrl *UserController) Login(c *gin.Context) {
	log.C(c).Infow("Login function called")

	var r v1.LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	resp, err := ctrl.b.Users().Login(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, resp)
}
