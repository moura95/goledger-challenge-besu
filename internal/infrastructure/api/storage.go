package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moura95/goledger-challenge-besu/pkg/ginx"
)

type SetRequest struct {
	Value int32 `json:"value" binding:"required"`
}

type SetResponse struct {
	Value         int32  `json:"value"`
	TransactionID string `json:"transaction_id"`
}

// @Summary Set Storage Value
// @Description Set a new value for the storage
// @Tags Storage
// @Accept json
// @Produce json
// @Param Storage body SetRequest true "Storage"
// @Success 200 {object} object{message=string} "Ok"
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /storage [put]
func (s *StorageAPI) set(ctx *gin.Context) {
	var req SetRequest
	s.logger.Info("Set Storage Value")

	err := ginx.ParseJSON(ctx, &req)
	if err != nil {
		s.logger.Error("Failed to parse request: ", err)
		ctx.JSON(http.StatusBadRequest, ginx.ErrorResponse("Failed to parse request"))
		return
	}

	hash, err := s.service.Set(req.Value)
	if err != nil {
		s.logger.Error("Failed to set storage: ", err)
		ctx.JSON(http.StatusInternalServerError, ginx.ErrorResponse("Failed to update storage value"))
		return
	}
	response := SetResponse{
		Value:         req.Value,
		TransactionID: hash,
	}

	ctx.JSON(http.StatusOK, ginx.SuccessResponse(response))
}

type GetResponse struct {
	Value    int32  `json:"value"`
	LastSync string `json:"last_sync"`
}

// @Summary Get Storage Value
// @Description Get the current value of the storage
// @Tags Storage
// @Accept json
// @Produce json
// @Success 200 {object} GetResponse
// @Failure 500 {object} object{error=string}
// @Router /storage [get]
func (s *StorageAPI) get(ctx *gin.Context) {
	s.logger.Info("Get Storage Value")

	storage, err := s.service.Get()
	if err != nil {
		s.logger.Error("Failed to get storage: ", err)
		ctx.JSON(http.StatusInternalServerError, ginx.ErrorResponse("Failed to retrieve storage value"))
		return
	}

	if storage == nil {
		s.logger.Warn("Storage is empty")
		ctx.JSON(http.StatusOK, ginx.SuccessResponse("No data available"))
		return
	}

	response := GetResponse{
		Value:    storage.Value,
		LastSync: storage.LastSync.Format("2006-01-02T15:04:05Z"),
	}

	ctx.JSON(http.StatusOK, ginx.SuccessResponse(response))
}
