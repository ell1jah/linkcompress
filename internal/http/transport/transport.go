package transport

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ell1jah/linkcompress/internal/http/domain"
)

type ReqBody struct {
	Link domain.Link `json:"link"`
}

type LinkService interface {
	Compress(domain.Link) (domain.Link, error)
	Original(domain.Link) (domain.Link, error)
}

type LinkHandlerLogger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

type LinkHandler struct {
	linkService LinkService
	logger      LinkHandlerLogger
}

func NewLinkHandler(linkService LinkService, logger LinkHandlerLogger) *LinkHandler {
	return &LinkHandler{
		linkService: linkService,
		logger:      logger,
	}
}

func (lh *LinkHandler) Get(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		lh.logger.Errorw("LinkHandler err",
			"method", "Get",
			"calling", "io.ReadAll",
			"error", err)
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	reqBody := &ReqBody{}
	err = json.Unmarshal(body, reqBody)
	if err != nil {
		lh.logger.Errorw("LinkHandler err",
			"method", "Get",
			"calling", "json.Unmarshal",
			"body", body,
			"error", err)
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	original, err := lh.linkService.Original(reqBody.Link)
	if err != nil {
		lh.logger.Errorw("LinkHandler err",
			"method", "Get",
			"calling", "lh.linkService.Original",
			"src link", reqBody.Link,
			"error", err)
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	respStruct := &ReqBody{Link: original}
	resp, err := json.Marshal(respStruct)
	if err != nil {
		lh.logger.Errorw("LinkHandler err",
			"method", "Get",
			"calling", "json.Marshal",
			"src link", reqBody.Link,
			"original link", original,
			"error", err)
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		lh.logger.Errorw("LinkHandler err",
			"method", "Get",
			"calling", "w.Write",
			"src link", reqBody.Link,
			"original link", original,
			"error", err)
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	lh.logger.Infow("LinkHandler log",
		"method", "Get",
		"src link", reqBody.Link,
		"original link", original)

	w.WriteHeader(http.StatusOK)
}

func (lh *LinkHandler) Post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		lh.logger.Errorw("LinkHandler err",
			"method", "Post",
			"calling", "io.ReadAll",
			"error", err)
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	reqBody := &ReqBody{}
	err = json.Unmarshal(body, reqBody)
	if err != nil {
		lh.logger.Errorw("LinkHandler err",
			"method", "Post",
			"calling", "json.Unmarshal",
			"body", body,
			"error", err)
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	compress, err := lh.linkService.Compress(reqBody.Link)
	if err != nil {
		lh.logger.Errorw("LinkHandler err",
			"method", "Post",
			"calling", "lh.linkService.Compress",
			"src link", reqBody.Link,
			"error", err)
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	respStruct := &ReqBody{Link: compress}
	resp, err := json.Marshal(respStruct)
	if err != nil {
		lh.logger.Errorw("LinkHandler err",
			"method", "Post",
			"calling", "json.Marshal",
			"src link", reqBody.Link,
			"compressed link", compress,
			"error", err)
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		lh.logger.Errorw("LinkHandler err",
			"method", "Post",
			"calling", "w.Write",
			"src link", reqBody.Link,
			"compressed link", compress,
			"error", err)
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	lh.logger.Infow("LinkHandler log",
		"method", "Post",
		"src link", reqBody.Link,
		"compressed link", compress)

	w.WriteHeader(http.StatusOK)
}
