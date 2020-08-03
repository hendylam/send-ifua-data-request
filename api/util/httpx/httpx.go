package httpx

import (
	"reflect"
	"strings"
	"time"

	"encoding/json"
	"fmt"
	"strconv"

	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog"

	"tapera/conf"
	"tapera/util/appcontext"
	"tapera/util/apperror"
	"tapera/util/identity"
)

const (
	ctjson = "application/json; charset=UTF-8"
	cttext = "text/plain; charset=utf-8"
)

type (
	// Error is general form that contains the actual information of either application or validation error that sent to user
	Error struct {
		ErrID string      `json:"errId" example:"453993f6-e433-445a-b1e4-8dbd70ef26e9"`
		Code  int         `json:"code" example:"400"`
		Msg   interface{} `json:"msg"`
	}

	// HTTPx func
	HTTPx struct {
		wr http.ResponseWriter
		rq *http.Request
	}

	// PaginationModel type
	PaginationModel struct {
		Page uint `json:"page" binding:"required"`
		Size uint `json:"size" binding:"required"`
	}
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		var name string
		tg := fld.Tag.Get("json")
		if len(tg) == 0 {
			tg = fld.Tag.Get("schema")
		}
		name = strings.SplitN(tg, ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	// validate date format : yyyy-mm-dd | 2020-02-20
	validate.RegisterValidation("isodateformat", isoDateValidation)
}

func isoDateValidation(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}

	return true
}

// New func
func New(w http.ResponseWriter, r *http.Request) *HTTPx {
	return &HTTPx{wr: w, rq: r}
}

// JSONerr func
func (h *HTTPx) JSONerr(code int, v interface{}) {
	var (
		id  string = h.ReqID()
		msg interface{}
	)

	if appErr, ok := v.(apperror.AppError); ok {
		msg = appErr.Error()
	} else if err, ok := v.(error); ok {
		switch conf.AppConf.Mode() {
		case "release":
			msg = "something wrong on our server"
		default:
			msg = err.Error()
		}
	} else {
		msg = v
	}

	h.jsonRs(code, &Error{
		ErrID: id,
		Code:  code,
		Msg:   msg,
	})
}

// JSON func
func (h *HTTPx) JSON(code int, val interface{}) {
	h.jsonRs(code, val)
}

// JSON func
func (h *HTTPx) String(code int, v string) {
	h.writeHeader(code, cttext)
	h.wr.Write([]byte(v))
}

// BindJSON func
func (h *HTTPx) BindJSON(v interface{}) error {
	decoder := json.NewDecoder(h.rq.Body)
	return decoder.Decode(v)
}

// BindJSONq func
func (h *HTTPx) BindJSONq(v interface{}) error {
	//https://stackoverflow.com/questions/39486740/how-to-get-all-get-request-query-parameters-into-a-structure-in-go
	if err := h.rq.ParseForm(); err != nil {
		return err
	}
	return schema.NewDecoder().Decode(v, h.rq.Form)
}

// Validate struct func
func (h *HTTPx) Validate(v interface{}) map[string]string {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	errs := make(map[string]string)
	verrs := err.(validator.ValidationErrors)
	for _, verr := range verrs {
		add := ""

		if verr.Kind() == reflect.String {
			add = "length "
		}

		name := verr.Field()

		switch verr.Tag() {
		case "required":
			errs[name] = "required"
		case "email":
			errs[name] = "should be a valid email"
		case "numeric":
			errs[name] = "should be a numeric value"
		case "eqfield":
			errs[name] = fmt.Sprintf("should be equal to the %s", verr.Param())
		case "lte":
			errs[name] = fmt.Sprintf("%sshould be less Than or equal to %s", add, verr.Param())
		case "gte":
			errs[name] = fmt.Sprintf("%sshould be greater than or equal to %s", add, verr.Param())
		case "gt":
			errs[name] = fmt.Sprintf("%sshould be greater than %s", add, verr.Param())
		case "lt":
			errs[name] = fmt.Sprintf("%sshould be lower than %s", add, verr.Param())
		case "min":
			errs[name] = fmt.Sprintf("%smin value should be %s", add, verr.Param())
		case "max":
			errs[name] = fmt.Sprintf("%smax value should be %s", add, verr.Param())
		case "isodateformat":
			errs[name] = fmt.Sprintf("should be a valid date yyyy-mm-dd")
		case "oneof":
			errs[name] = fmt.Sprintf("value should be one of (%s)", verr.Param())
		default:
			errs[name] = fmt.Sprintf("%s invalid", name)
		}
	}
	return errs
}

// Param func
func (h *HTTPx) Param(k string) string {
	return mux.Vars(h.rq)[k]
}

// Query func
func (h *HTTPx) Query(k string) string {
	return h.rq.URL.Query().Get(k)
}

// Header func
func (h *HTTPx) Header(k string) string {
	return h.rq.Header.Get(k)
}

// Pagination func to get pagination variable from param
func (h *HTTPx) Pagination() (uint, uint, error) {
	pg, err := uintVal(h.Param("page"), "page")
	if err != nil {
		return 0, 0, err
	}

	sz, err := uintVal(h.Param("size"), "size")
	if err != nil {
		return 0, 0, err
	}

	return pg, sz, nil
}

// PaginationQ func to get pagination variable from http query string
func (h *HTTPx) PaginationQ() (uint, uint, error) {
	pg, err := uintVal(h.Param("page"), "page")
	if err != nil {
		return 0, 0, err
	}

	sz, err := uintVal(h.Param("size"), "size")
	if err != nil {
		return 0, 0, err
	}

	return pg, sz, nil
}

// ReqID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func (h *HTTPx) ReqID() string {
	return appcontext.ReqID(h.rq.Context())
}

// Logger func, get logger by key from context
func (h *HTTPx) Logger() *zerolog.Logger {
	return appcontext.Logger(h.rq.Context())
}

// User func
func (h *HTTPx) User() identity.UserInfo {
	return appcontext.User(h.rq.Context())
}

func (h *HTTPx) jsonRs(code int, v interface{}) {
	h.writeHeader(code, ctjson)
	if err := json.NewEncoder(h.wr).Encode(v); err != nil {
		panic(err)
	}
}

func (h *HTTPx) writeHeader(code int, ctype string) {
	h.wr.Header().Set("Content-Type", ctype)
	h.wr.WriteHeader(code)
}

func uintVal(val string, name string) (uint, error) {
	var (
		i   int
		err error
	)

	if val == "" {
		err = fmt.Errorf("%s is required", name)
	} else {
		i, err = strconv.Atoi(val)
		if err == nil {
			if i < 0 {
				i = 0
				err = fmt.Errorf("%s must be greater than zero", name)
			}
		}
	}

	return uint(i), err
}
