package files

import (
	"errors"

	"gitlab.com/_spacemc_/web/users/internal/domain/ports"

	"github.com/gin-gonic/gin"
	usersErrors "gitlab.com/_spacemc_/web/users/errors"

	"gitlab.com/_spacemc_/web/gokit/ginx"
)

type Handler struct {
	service ports.FilesServicePort
}

func NewHandler(service ports.FilesServicePort) *Handler {
	return &Handler{service: service}
}

// GetSkin godoc
// @Summary      Получить URL скина пользователя
// @Description  Получить URL скина пользователя по его ID
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID пользователя"
// @Success      200  {object} ginx.SuccessResponse[UploadResponse] "Успешный ответ"
// @Failure      400  {object}  ginx.ErrorResponse "Неверный ID"
// @Failure      404  {object}  ginx.ErrorResponse "Пользователя не существует"
// @Router       /users/{id}/skin [get]
func (h *Handler) GetSkin(c *gin.Context) {
	p := ginx.NewGinxParser(c)

	userID, err := p.GetPathString("id")

	if err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
		return
	}

	skinURL, err := h.service.GetUserSkin(c.Request.Context(), *userID)
	if err != nil {
		switch {
		case errors.Is(err, usersErrors.ErrUserNotFound):
			ginx.WriteNotFound(c)
		}
		return

	}

	response := GetSkin{
		Skin: skinURL,
	}

	ginx.WriteSuccessResponse(c, &response)
}

// GetCloak godoc
// @Summary      Получить URL плаща пользователя
// @Description  Получить URL плаща пользователя по его ID
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID пользователя"
// @Success      200  {object} ginx.SuccessResponse[UploadResponse] "Успешный ответ"
// @Failure      400  {object}  ginx.ErrorResponse "Неверный ID"
// @Failure      404  {object}  ginx.ErrorResponse "Пользователя не существует"
// @Router       /users/{id}/cloak [get]
func (h *Handler) GetCloak(c *gin.Context) {
	p := ginx.NewGinxParser(c)

	userID, err := p.GetPathString("id")

	if err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
		return
	}

	cloakURL, err := h.service.GetUserCloak(c.Request.Context(), *userID)
	if err != nil {
		switch {
		case errors.Is(err, usersErrors.ErrUserNotFound):
			ginx.WriteNotFound(c)
		}

		return

	}

	response := GetCloak{
		Cloak: cloakURL,
	}

	ginx.WriteSuccessResponse(c, &response)
}

// UploadSkin godoc
// @Summary      Загрузить скин пользователю
// @Description  Загружает скин пользователя по его ID. Вырезает голову из скина и загружает в avatar
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID пользователя"
// @Param 		 skin formData file true "Файл скина (png)"
// @Success      200  {object} ginx.SuccessResponse[UploadResponse] "Успешный ответ"
// @Failure      404  {object}  ginx.ErrorResponse "Пользователя не существует"
// @Failure      500  {object}  ginx.ErrorResponse "Ошибка загрузки файла"
// @Router       /users/{id}/skin [post]
func (h *Handler) UploadSkin(c *gin.Context) {
	p := ginx.NewGinxParser(c)

	userID, err := p.GetPathString("id")

	if err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
		return
	}

	file, err := c.FormFile("skin")
	if err != nil {
		ginx.WriteBadRequest(c)
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		ginx.WriteInternalError(c, usersErrors.FailedToOpenFile.Error())
		return
	}
	defer fileReader.Close()

	url, err := h.service.UploadSkin(c.Request.Context(), *userID, fileReader, file.Filename, file.Size)

	if err != nil {
		ginx.WriteInternalError(c, err.Error())
		return
	}

	response := UploadResponse{
		Message: "skin uploaded successfully",
		UserID:  *userID,
		Url:     url,
	}

	ginx.WriteSuccessResponse(c, &response)
}

// DeleteSkin godoc
// @Summary      Удалить скин пользоваетя
// @Description  Удаляет скин пользователя по его ID
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID пользователя"
// @Success      200  {object} ginx.SuccessResponse[Response] "Успешный ответ"
// @Failure      404  {object}  ginx.ErrorResponse "Пользователя не существует"
// @Failure      500  {object}  ginx.ErrorResponse "Ошибка удаления"
// @Router       /users/{id}/skin [delete]
func (h *Handler) DeleteSkin(c *gin.Context) {
	p := ginx.NewGinxParser(c)

	userID, err := p.GetPathString("id")

	if err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
		return
	}

	err = h.service.DeleteSkin(c.Request.Context(), *userID) // Добавить 404 ошибку если юзера нет

	if err != nil {
		ginx.WriteInternalError(c, err.Error())
		return
	}

	response := Response{
		Message: "skin delete successfully",
	}

	ginx.WriteSuccessResponse(c, &response)
}

// UploadCloak godoc
// @Summary      Загрузить плащ пользователю
// @Description  Загружает плащ пользователя по его ID
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID пользователя"
// @Param 		 cloak formData file true "Файл плаща (png)"
// @Success      200  {object} ginx.SuccessResponse[UploadResponse] "Успешный ответ"
// @Failure      404  {object}  ginx.ErrorResponse "Пользователя не существует"
// @Failure      500  {object}  ginx.ErrorResponse "Ошибка загрузки файла"
// @Router       /users/{id}/cloak [post]
func (h *Handler) UploadCloak(c *gin.Context) {
	p := ginx.NewGinxParser(c)

	userID, err := p.GetPathString("id")

	if err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
		return
	}

	file, err := c.FormFile("cloak")
	if err != nil {
		ginx.WriteBadRequest(c)
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		ginx.WriteInternalError(c, usersErrors.FailedToOpenFile.Error())
		return
	}
	defer fileReader.Close()

	url, err := h.service.UploadCloak(c.Request.Context(), *userID, fileReader, file.Filename, file.Size)

	if err != nil {
		ginx.WriteInternalError(c, err.Error())
		return
	}

	response := UploadResponse{
		Message: "cloak uploaded successfully",
		UserID:  *userID,
		Url:     url,
	}

	ginx.WriteSuccessResponse(c, &response)
}

// DeleteCloak godoc
// @Summary      Удалить плащ пользоваетя
// @Description  Удаляет плащ пользователя по его ID
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID пользователя"
// @Success      200  {object} ginx.SuccessResponse[Response] "Успешный ответ"
// @Failure      404  {object}  ginx.ErrorResponse "Пользователя не существует"
// @Failure      500  {object}  ginx.ErrorResponse "Ошибка удаления"
// @Router       /users/{id}/cloak [delete]
func (h *Handler) DeleteCloak(c *gin.Context) {
	p := ginx.NewGinxParser(c)

	userID, err := p.GetPathString("id")

	if err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
		return
	}

	err = h.service.DeleteCloak(c.Request.Context(), *userID)

	if err != nil {
		ginx.WriteInternalError(c, err.Error())
		return
	}

	response := Response{
		Message: "cloak delete successfully",
	}

	ginx.WriteSuccessResponse(c, &response)
}
