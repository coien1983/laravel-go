package http

import (
	"encoding/json"
	"net/http"
)

// Response HTTP响应接口
type Response interface {
	// 状态码
	Status() int
	// 数据
	Data() interface{}
	// 头部
	Headers() map[string]string
	SetHeader(key, value string) Response
	// 发送响应
	Send(w http.ResponseWriter)
}

// response HTTP响应实现
type response struct {
	status  int
	data    interface{}
	headers map[string]string
}

func NewResponse(status int, data interface{}) Response {
	return &response{
		status:  status,
		data:    data,
		headers: make(map[string]string),
	}
}

func (r *response) Status() int {
	return r.status
}

func (r *response) Data() interface{} {
	return r.data
}

func (r *response) Headers() map[string]string {
	return r.headers
}

func (r *response) SetHeader(key, value string) Response {
	r.headers[key] = value
	return r
}

func (r *response) Send(w http.ResponseWriter) {
	// 设置头部
	for k, v := range r.headers {
		w.Header().Set(k, v)
	}

	// 设置状态码
	w.WriteHeader(r.status)

	// 发送数据
	switch data := r.data.(type) {
	case string:
		w.Write([]byte(data))
	case []byte:
		w.Write(data)
	default:
		json.NewEncoder(w).Encode(data)
	}
}

// JsonResponse JSON响应
type JsonResponse struct {
	response
}

func NewJsonResponse(status int, data interface{}) *JsonResponse {
	return &JsonResponse{
		response: response{
			status:  status,
			data:    data,
			headers: make(map[string]string),
		},
	}
}

func (r *JsonResponse) Send(w http.ResponseWriter) {
	// 设置JSON头部
	w.Header().Set("Content-Type", "application/json")

	// 设置其他头部
	for k, v := range r.headers {
		w.Header().Set(k, v)
	}

	// 设置状态码
	w.WriteHeader(r.status)

	// 发送JSON数据
	json.NewEncoder(w).Encode(r.data)
}

// TextResponse 文本响应
type TextResponse struct {
	response
}

func NewTextResponse(status int, text string) *TextResponse {
	return &TextResponse{
		response: response{
			status:  status,
			data:    text,
			headers: make(map[string]string),
		},
	}
}

func (r *TextResponse) Send(w http.ResponseWriter) {
	// 设置文本头部
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// 设置其他头部
	for k, v := range r.headers {
		w.Header().Set(k, v)
	}

	// 设置状态码
	w.WriteHeader(r.status)

	// 发送文本数据
	w.Write([]byte(r.data.(string)))
}

// RedirectResponse 重定向响应
type RedirectResponse struct {
	response
	location string
}

func NewRedirectResponse(status int, location string) *RedirectResponse {
	return &RedirectResponse{
		response: response{
			status:  status,
			data:    nil,
			headers: make(map[string]string),
		},
		location: location,
	}
}

func (r *RedirectResponse) Send(w http.ResponseWriter) {
	// 设置重定向头部
	w.Header().Set("Location", r.location)

	// 设置其他头部
	for k, v := range r.headers {
		w.Header().Set(k, v)
	}

	// 设置状态码并重定向
	http.Redirect(w, nil, r.location, r.status)
}

// FileResponse 文件响应
type FileResponse struct {
	response
	filename string
}

func NewFileResponse(filename string) *FileResponse {
	return &FileResponse{
		response: response{
			status:  http.StatusOK,
			data:    nil,
			headers: make(map[string]string),
		},
		filename: filename,
	}
}

func (r *FileResponse) Send(w http.ResponseWriter) {
	// 设置其他头部
	for k, v := range r.headers {
		w.Header().Set(k, v)
	}

	// 发送文件
	http.ServeFile(w, nil, r.filename)
}

// Controller 基础控制器
type Controller struct{}

func (c *Controller) Json(data interface{}, status ...int) Response {
	s := http.StatusOK
	if len(status) > 0 {
		s = status[0]
	}
	return NewJsonResponse(s, data)
}

func (c *Controller) Text(text string, status ...int) Response {
	s := http.StatusOK
	if len(status) > 0 {
		s = status[0]
	}
	return NewTextResponse(s, text)
}

func (c *Controller) Redirect(location string, status ...int) Response {
	s := http.StatusFound
	if len(status) > 0 {
		s = status[0]
	}
	return NewRedirectResponse(s, location)
}

func (c *Controller) File(filename string) Response {
	return NewFileResponse(filename)
}

func (c *Controller) Success(data interface{}, message ...string) Response {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}

	response := map[string]interface{}{
		"success": true,
		"message": msg,
		"data":    data,
	}

	return NewJsonResponse(http.StatusOK, response)
}

func (c *Controller) Error(message string, status ...int) Response {
	s := http.StatusInternalServerError
	if len(status) > 0 {
		s = status[0]
	}

	response := map[string]interface{}{
		"success": false,
		"message": message,
	}

	return NewJsonResponse(s, response)
}

func (c *Controller) ValidationError(errors interface{}) Response {
	response := map[string]interface{}{
		"success": false,
		"message": "Validation failed",
		"errors":  errors,
	}

	return NewJsonResponse(http.StatusUnprocessableEntity, response)
}

func (c *Controller) NotFound(message ...string) Response {
	msg := "Resource not found"
	if len(message) > 0 {
		msg = message[0]
	}

	response := map[string]interface{}{
		"success": false,
		"message": msg,
	}

	return NewJsonResponse(http.StatusNotFound, response)
}

func (c *Controller) Unauthorized(message ...string) Response {
	msg := "Unauthorized"
	if len(message) > 0 {
		msg = message[0]
	}

	response := map[string]interface{}{
		"success": false,
		"message": msg,
	}

	return NewJsonResponse(http.StatusUnauthorized, response)
}

func (c *Controller) Forbidden(message ...string) Response {
	msg := "Forbidden"
	if len(message) > 0 {
		msg = message[0]
	}

	response := map[string]interface{}{
		"success": false,
		"message": msg,
	}

	return NewJsonResponse(http.StatusForbidden, response)
}

func (c *Controller) InternalServerError(message ...string) Response {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}

	response := map[string]interface{}{
		"success": false,
		"message": msg,
	}

	return NewJsonResponse(http.StatusInternalServerError, response)
}
