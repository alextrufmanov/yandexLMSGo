package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alextrufmanov/yandexLMSGo/pkg/calc"
)

type cfg struct {
	Addr string
}

func cfgEnv() *cfg {
	cfg := new(cfg)
	cfg.Addr = os.Getenv("PORT")
	if cfg.Addr == "" {
		cfg.Addr = "8080"
	}
	return cfg
}

type Application struct {
	cfg    *cfg
	logger *log.Logger
}

func New() *Application {
	return &Application{
		cfg:    cfgEnv(),
		logger: log.New(os.Stdout, "[APP]", log.Flags()),
	}
}

type RequestBody struct {
	Expression string `json:"expression"`
}

type AnswerBody struct {
	Result string `json:"result"`
}

func NotFoundError404(w http.ResponseWriter) {
	http.Error(w, "404 not found.", http.StatusNotFound)
}

func CalcError422(w http.ResponseWriter) {
	http.Error(w, "{\"error\":\"Expression is not valid\"}", http.StatusUnprocessableEntity)
}

func InternalError500(w http.ResponseWriter) {
	http.Error(w, "{\"error\":\"Internal server error\"}", http.StatusInternalServerError)
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {

	var requestBody RequestBody
	var answerBody AnswerBody

	if strings.ToUpper(r.Method) != "POST" {
		NotFoundError404(w)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err == nil {
		if json.Unmarshal(bodyBytes, &requestBody) == nil {
			calcResult, err := calc.Calc(requestBody.Expression)
			if err == nil {
				answerBody.Result = fmt.Sprint(calcResult)
				outData, _ := json.Marshal(answerBody)
				fmt.Fprintf(w, string(outData))
			} else {
				CalcError422(w)
			}
			return
		}
	}
	InternalError500(w)
}

func (a *Application) StartRESTServer() error {
	a.logger.Printf("Listen localhost:%s", a.cfg.Addr)
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.cfg.Addr, nil)
}
