// Copyright 2023 lawff. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package errno

var (
	// ErrUserAlreadyExist 代表用户已经存在.
	ErrUserAlreadyExist = &Errno{HTTP: 400, Code: "FailedOperation.UserAlreadyExist", Message: "User already exist."}
)
