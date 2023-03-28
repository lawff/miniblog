// Copyright 2023 lawff. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package miniblog

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lawff/miniblog/internal/miniblog/controller/v1/user"
	"github.com/lawff/miniblog/internal/miniblog/store"
	"github.com/lawff/miniblog/internal/pkg/core"
	"github.com/lawff/miniblog/internal/pkg/errno"
	"github.com/lawff/miniblog/internal/pkg/log"
	mw "github.com/lawff/miniblog/internal/pkg/middleware"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func convertStringToList(str string) []string {
	list := strings.Split(str, ",")
	for i, val := range list {
		list[i] = "./upload/" + val + ".pdf" // 这里将每个值转换成文件路径
	}
	return list
}

// installRouters 安装 miniblog 接口路由.
func installRouters(g *gin.Engine) error {
	// 注册 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// 注册 /healthz handler.
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	g.POST("/upload", func(c *gin.Context) {
		log.C(c).Infow("file upload function called")

		file, _ := c.FormFile("file")
		name := uuid.New().String() + ".pdf"
		p := filepath.Join("./upload")
		// 上传文件至指定目录
		if err := c.SaveUploadedFile(file, filepath.Join(p, name)); err != nil {
			log.C(c).Errorw("Save uploaded file failed", "error", err)
			core.WriteResponse(c, err, nil)

			return
		}

		core.WriteResponse(c, nil, map[string]string{"status": "ok", "data": name})
	})

	g.GET("/upload", func(c *gin.Context) {
		log.C(c).Infow("file download function called")
		id := c.Query("id")

		// 将两个pdf文件合并为一个pdf文件
		outputFilePath := "./upload/" + uuid.New().String() + ".pdf"
		if err := api.MergeCreateFile(convertStringToList(id), outputFilePath, nil); err != nil {
			log.C(c).Errorw("Merge pdf file failed", "error", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// 打开要返回的PDF文件
		pdfPath := outputFilePath
		file, err := os.Open(pdfPath)
		if err != nil {
			log.C(c).Errorw("Open pdf file failed", "error", err)
			return
		}
		defer file.Close()

		// 将PDF文件作为二进制流返回给请求接口
		fileInfo, err := file.Stat()
		if err != nil {
			log.C(c).Errorw("Open pdf file failed", "error", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		fileSize := fileInfo.Size()

		buffer := make([]byte, fileSize)
		_, err = file.Read(buffer)
		if err != nil {
			log.C(c).Errorw("Open pdf file failed", "error", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		filename := filepath.Base(pdfPath)
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/pdf", buffer)

	})

	uc := user.New(store.S)

	g.POST("/login", uc.Login)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 users 路由分组
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
			userv1.PUT(":name/change-password", uc.ChangePassword)
			userv1.Use(mw.Authn())
		}
	}

	return nil
}
