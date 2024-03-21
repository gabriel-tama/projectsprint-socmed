package image

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gabriel-tama/projectsprint-socmed/common/response"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	s3Service S3Service
}

func NewImageController(s3Service S3Service) *ImageController {
	return &ImageController{s3Service: s3Service}
}

func (ic *ImageController) UploadImage(c *gin.Context) {
	var data = &ImageResponse{}
	res := &response.ResponseBody{
		Message: "file not found",
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Check file size
	if file.Size > 2*1024*1024 || file.Size < 10*1024 {
		res.Message = "image is wrong (not *.jpg | *.jpeg, more than 2MB or less than 10KB)"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	fileHeader, err := file.Open()
	if err != nil {
		res.Message = "file not found"
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	defer fileHeader.Close()

	buffer := make([]byte, 512) // Why 512 bytes? See http://golang.org/pkg/net/http/#DetectContentType
	_, err = fileHeader.Read(buffer)
	if err != nil {
		res.Message = "file not found"
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" {
		res.Message = "image is wrong (not *.jpg | *.jpeg, more than 2MB or less than 10KB)"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	fileBuffer, err := file.Open()
	if err != nil {
		res.Message = "file not found"
		c.JSON(http.StatusBadRequest, res)
		return
	}
	fileBuffer.Close()

	objKey := fmt.Sprintf("%s/%v-%s", "ngab-gab", time.Now().Unix(), file.Filename)

	_, err = ic.s3Service.UploadFile(objKey, fileBuffer, contentType)
	if err != nil {
		res.Message = "failed to upload image"
		c.JSON(http.StatusInternalServerError, res.Message)
		return
	}

	data.Url = ic.s3Service.GetObjectWithUrl(objKey)

	res = &response.ResponseBody{
		Message: "image uploaded successfully",
		Data:    data,
	}

	c.JSON(http.StatusOK, res)
}
