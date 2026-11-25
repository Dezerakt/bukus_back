package v1

import (
	"bukus/internal/controller/restapi/v1/request"
	"bukus/internal/controller/restapi/v1/response"
	"bukus/internal/entity"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func (obj *V1) store(ctx fiber.Ctx) error {
	var body request.StoreBook

	if err := ctx.Bind().Body(&body); err != nil {
		obj.l.Error(err, "restapi - v1 - store - bind")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := obj.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "validation failed: "+err.Error())
	}

	book := entity.Book{
		Name: body.Name,
	}

	if err := obj.bookUC.Store(ctx.RequestCtx(), book); err != nil {
		obj.l.Error(err, "restapi - v1 - store")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to store book")
	}

	return ctx.Status(http.StatusCreated).JSON(nil)
}

func (obj *V1) list(ctx fiber.Ctx) error {
	books, err := obj.bookUC.List(ctx.RequestCtx())
	if err != nil {
		obj.l.Error(err, "restapi - v1 - list")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to list books")
	}

	resp := make([]response.Book, 0, len(books))
	for _, b := range books {
		resp = append(resp, toBookResponse(b))
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func (obj *V1) getByID(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, http.StatusBadRequest, "id is required")
	}

	book, err := obj.bookUC.GetByID(ctx.RequestCtx(), id)
	if err != nil {
		obj.l.Error(err, "restapi - v1 - getByID")
		return errorResponse(ctx, http.StatusNotFound, "book not found")
	}

	return ctx.Status(http.StatusOK).JSON(toBookResponse(book))
}

func (obj *V1) update(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, http.StatusBadRequest, "id is required")
	}

	var body request.UpdateBook
	if err := ctx.Bind().Body(&body); err != nil {
		obj.l.Error(err, "restapi - v1 - update - bind")
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := obj.v.Struct(body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "validation failed: "+err.Error())
	}

	book := entity.Book{
		ID:   id,
		Name: body.Name,
	}

	if err := obj.bookUC.Update(ctx.RequestCtx(), book); err != nil {
		obj.l.Error(err, "restapi - v1 - update")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to update book")
	}

	return ctx.Status(http.StatusOK).JSON(nil)
}

func (obj *V1) delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return errorResponse(ctx, http.StatusBadRequest, "id is required")
	}

	if err := obj.bookUC.Delete(ctx.RequestCtx(), id); err != nil {
		obj.l.Error(err, "restapi - v1 - delete")
		return errorResponse(ctx, http.StatusInternalServerError, "failed to delete book")
	}

	return ctx.Status(http.StatusNoContent).JSON(nil)
}

func toBookResponse(b entity.Book) response.Book {
	return response.Book{
		ID:        b.ID,
		Name:      b.Name,
		UserID:    b.UserID,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
