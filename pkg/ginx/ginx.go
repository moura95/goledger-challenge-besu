package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
}

func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return err
	}
	return nil
}

func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return err
	}
	return nil
}
func ParseUri(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindUri(obj); err != nil {
		return err
	}
	return nil
}

type Response struct {
	Error interface{} `json:"error"`
	Data  interface{} `json:"data"`
}
type ResponseWithPageInfo struct {
	Error    interface{} `json:"error"`
	Data     interface{} `json:"data"`
	PageInfo PageInfo    `json:"pageInfo"`
}

func SuccessResponseWithPageInfo(data interface{}, pageInfo PageInfo) ResponseWithPageInfo {
	return ResponseWithPageInfo{
		Data:     data,
		Error:    "",
		PageInfo: pageInfo,
	}
}

func SuccessResponse(data interface{}) Response {
	return Response{
		Data:  data,
		Error: "",
	}
}

func ErrorResponse(error string) Response {
	return Response{
		Data:  "",
		Error: error,
	}
}

type PageInfo struct {
	TotalRecords int
	CurrentPage  int
	TotalPages   int
	Limit        int
}

func NewPageInfo(limit, totalRecords, currentPage, totalPages int) PageInfo {
	return PageInfo{
		TotalRecords: totalRecords,
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		Limit:        limit,
	}
}
