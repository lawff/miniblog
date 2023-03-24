// Copyright 2023 lawff. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/lawff/miniblog/internal/miniblog/biz"
	"github.com/lawff/miniblog/internal/miniblog/store"
)

// UserController 是 user 模块在 Controller 层的实现，用来处理用户模块的请求.
type UserController struct {
	b biz.IBiz
}

// New 创建一个 user controller.
func New(ds store.IStore) *UserController {
	return &UserController{b: biz.NewBiz(ds)}
}
