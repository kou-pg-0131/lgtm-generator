package controllers

import "github.com/koki-develop/lgtm-generator/backend/src/entities"

type Renderer interface {
	// 200
	OK(ctx Context, obj interface{})
	// 201
	Created(ctx Context, obj interface{})
	// 204
	NoContent(ctx Context)
	// 400
	BadRequest(ctx Context, code entities.ErrCode, err error)
	// 403
	Forbidden(ctx Context, code entities.ErrCode, err error)
	// 404
	NotFound(ctx Context, code entities.ErrCode, err error)
	// 500
	InternalServerError(ctx Context, err error)
}
