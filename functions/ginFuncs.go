package functions

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
)

func PrepareCtx(method string) (*httptest.ResponseRecorder, *gin.Context) {
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ctx.Request.Method = method
	return resp, ctx
}

func SetReqBody(ctx *gin.Context, req any) {
	json_bytes, _ := json.Marshal(req)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(json_bytes))
	ctx.Request.Header.Set("Content-Type", "application/json")
}

func GetRespBody(resp *httptest.ResponseRecorder, target any) error {
	json_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(json_bytes, target)
	return err
}
