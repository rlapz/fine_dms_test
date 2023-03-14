package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"enigmacamp.com/fine_dms/config"
	"enigmacamp.com/fine_dms/middleware"
	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repo"
	"enigmacamp.com/fine_dms/usecase"
	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileUsecase usecase.FileUsecase
}

func NewFileController(router *gin.RouterGroup, f usecase.FileUsecase, cfg *config.Secret) {
	fc := &FileController{fileUsecase: f}
	authMiddleware := middleware.ValidateToken(cfg.Key)

	router.GET("/files/:user_id", authMiddleware, fc.getFilesByUserId)
	router.PUT("/files/:user_id", authMiddleware, fc.updateFile)
	router.GET("/files/:user_id/search", fc.searchFilesByUserId)
}

func (fc *FileController) getFilesByUserId(ctx *gin.Context) {
	id, err := GetUserId(ctx)
	if err != nil {
		return
	}
	files, err := fc.fileUsecase.GetFilesByUserId(id)
	if err != nil {
		if errors.Is(err, usecase.ErrUsecaseNoData) {
			FailedJSONResponse(ctx, http.StatusNotFound, "no data")
		} else {
			FailedJSONResponse(ctx, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	SuccessJSONResponse(ctx, http.StatusOK, "files retrieved successfully", files)
}

func (fc *FileController) updateFile(ctx *gin.Context) {
	userID, err := GetUserId(ctx)
	if err != nil {
		return
	}

	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var file model.File

	if err := ctx.BindJSON(&file); err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = fc.fileUsecase.UpdateFile(userID, file.Path, file.Ext)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidFileData) {
			FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, usecase.ErrUsecaseNoData) {
			FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		} else {
			FailedJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		}

		return
	}
	SuccessJSONResponse(ctx, http.StatusOK,
		fmt.Sprintf("file with id = %d has been updated", userID),
		nil,
	)
}

func (fc *FileController) deleteFile(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fileId, err := strconv.Atoi(ctx.Query("file_id"))
	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, "invalid file id")
		return
	}

	token := ctx.GetHeader("Authorization")
	if token == "" {
		FailedJSONResponse(ctx, http.StatusUnauthorized, "missing token")
		return
	}

	err = fc.fileUsecase.DeleteFile(userId, fileId, token, []byte("secret_key"))
	if err != nil {
		if errors.Is(err, usecase.ErrUsecaseInvalidAuth) {
			FailedJSONResponse(ctx, http.StatusUnauthorized, "invalid authentication")
		} else if errors.Is(err, usecase.ErrUsecaseNoData) {
			FailedJSONResponse(ctx, http.StatusNotFound, "no data")
		} else {
			FailedJSONResponse(ctx, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	SuccessJSONResponse(ctx, http.StatusOK, "file deleted successfully", nil)
}

func (fc *FileController) searchFilesByUserId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		FailedJSONResponse(ctx, http.StatusBadRequest, "invalid user id")
		return
	}
	query := ctx.Query("q")
	if query == "" {
		FailedJSONResponse(ctx, http.StatusBadRequest, "invalid query")
		return
	}

	files, err := fc.fileUsecase.SearchFilesByUserId(id, query)
	if err != nil {
		if errors.Is(err, repo.ErrRepoNoData) {
			FailedJSONResponse(ctx, http.StatusNotFound, "no data")
		} else if errors.Is(err, usecase.ErrInvalidUserID) {
			FailedJSONResponse(ctx, http.StatusBadRequest, err.Error())
		} else {
			FailedJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		}

		return
	}
	SuccessJSONResponse(ctx, http.StatusOK, "file search successfully", files)
}
