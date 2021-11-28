package internalhttp

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func (s *Server) GetFib(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.Atoi(r.URL.Query().Get("x"))
	if err != nil {
		zap.L().Error("Ger fib err: ", zap.Error(err))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	y, err := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil {
		zap.L().Error("Ger fib err: ", zap.Error(err))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	fib, err := s.App.GetFibSlice(r.Context(), x, y)
	if err != nil {
		zap.L().Error("Ger fib err: ", zap.Error(err))
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	_, _ = io.WriteString(w, strings.Join(fib, ", "))
}
