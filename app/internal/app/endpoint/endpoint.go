package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	uu "github.com/satori/go.uuid"

	"go.uber.org/zap"
	"log"
	"net/http"
	"pehks1980/be2_hw81/internal/pkg/model"
)

// App - application contents & methods
type App struct {
	Logger *zap.Logger
	Repository RepoIf
	CTX context.Context
}
// RepoIf - repository interface (PG)
type RepoIf interface {
	New(ctx context.Context, filename, filename1 string) RepoIf
	CloseConn()
	///
	AuthUser(ctx context.Context, user model.User) (string, error)
	GetUser(ctx context.Context, name string) (model.User, error)
	AddUpdUser(ctx context.Context, user model.User) (string, error)
	DelUser(ctx context.Context, id uuid.UUID) error
	GetUserEnvs(ctx context.Context, name string) (model.Envs, error)
	// env
	AddUpdEnv(ctx context.Context, env model.Environment) (string, error)
	GetEnv(ctx context.Context, title string) (model.Environment, error)
	DelEnv(ctx context.Context, id uuid.UUID) error
	GetEnvUsers(ctx context.Context, title string) (model.Users, error)
	// goods
	AddUpdGood(ctx context.Context, good model.Good) (string, error)
	GetGood(ctx context.Context, title string) (model.Good, error)
	DelGood(ctx context.Context, id uuid.UUID) error
	FindGood(ctx context.Context, key string) ([]model.Good, error)
}
// RegisterPublicHTTP - регистрация роутинга путей типа urls.py для обработки сервером
func (app *App) RegisterPublicHTTP() *mux.Router {
	r := mux.NewRouter()
	// authorization
	r.HandleFunc("/user/auth", app.postAuth()).Methods(http.MethodPost)
	// user crud
	r.HandleFunc("/user/", app.putUser()).Methods(http.MethodPost)
	r.HandleFunc("/user/{uid}", app.getUser()).Methods(http.MethodGet)
	r.HandleFunc("/user/{uid}", app.putUser()).Methods(http.MethodPut)
	r.HandleFunc("/user/{uid}", app.delUser()).Methods(http.MethodDelete)
	// env crud
	r.HandleFunc("/env/", app.putEnv()).Methods(http.MethodPost)
	r.HandleFunc("/env/{uid}", app.getEnv()).Methods(http.MethodGet)
	r.HandleFunc("/env/{uid}", app.putEnv()).Methods(http.MethodPut)
	r.HandleFunc("/env/{uid}", app.delEnv()).Methods(http.MethodDelete)

	// GetUserEnvs
	r.HandleFunc("/user/envs", app.postUserEnvs()).Methods(http.MethodPost)
	// GetEnvUsers
	r.HandleFunc("/env/users", app.postEnvUsers()).Methods(http.MethodPost)

	// goods crud
	r.HandleFunc("/good/find/", app.findGood()).Methods(http.MethodGet)
	r.HandleFunc("/good/", app.putGood()).Methods(http.MethodPost)
	r.HandleFunc("/good/{uid}", app.getGood()).Methods(http.MethodGet)
	r.HandleFunc("/good/{uid}", app.putGood()).Methods(http.MethodPut)
	r.HandleFunc("/good/{uid}", app.delGood()).Methods(http.MethodDelete)

	return r
}
// write response in json format
func writeResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
	_, _ = w.Write([]byte("\n"))
}
// write response in json format
func writeJsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't marshal data: %s", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	writeResponse(w, status, string(response))
}
// postAuth - user auth method
func (app *App) postAuth() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		ctx := request.Context()

		defer func() {
			// update Prom objects AuthCounter tries
		}()

		//json header check
		contentType := request.Header.Get("Content-Type")
		if contentType != "application/json" {
			return
		}

		user := model.User{}

		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			return
		}

		UID, err1 := app.Repository.AuthUser(ctx, user)

		if err1 != nil || UID == "" {
			log.Printf("USER %s Log in error.\n", user.Name)
			return
		}
		log.Printf("USER %s Logged in.\n", user.Name)

		writeJsonResponse(w, http.StatusOK, UID)

	}
}
// putUser - add or update user
func (app *App) putUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		user := model.User{}
		id, _ := app.Repository.AddUpdUser(ctx, user)
		writeJsonResponse(w, http.StatusOK, id)
	}
}
// delUser - delete user
func (app *App) delUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		var id uuid.UUID
		_ = app.Repository.DelUser(ctx, id)
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// getUser - get user
func (app *App) getUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		var (
			name string
			user model.User
		)
		user, _ = app.Repository.GetUser(ctx, name)
		log.Printf("getUser = %v", user)
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// putEnv - add or update environment
func (app *App) putEnv() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		//ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// getEnv - get Env
func (app *App) getEnv() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		//ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// delEnv - delete Env
func (app *App) delEnv() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		//ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// postUserEnvs - get envs of which user is member of
func (app *App) postUserEnvs() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		var (
			envs model.Envs
			name string
		)
		envs, _ = app.Repository.GetUserEnvs(ctx,name)
		log.Printf("GetUserEnvs(%s) = %v", name, envs)
		writeJsonResponse(w, http.StatusOK, "")
	}
}
// postEnvUsers - get Users which have this env membership
func (app *App) postEnvUsers() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		var (
			users model.Users
			title string
		)
		users, _ = app.Repository.GetEnvUsers(ctx, title)
		log.Printf("GetEnvUsers(%s) = %v", title, users)
		writeJsonResponse(w, http.StatusOK, "")
	}
}

// putGood - add or update good
func (app *App) putGood() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		defer func() {
			// update Prom objects
		}()

		//json header check
		contentType := request.Header.Get("Content-Type")
		if contentType != "application/json" {
			return
		}

		good := model.Good{}

		err := json.NewDecoder(request.Body).Decode(&good)
		if err != nil {
			return
		}

		res, err1 := app.Repository.AddUpdGood(ctx, good)

		if err1 != nil || res == "" {
			log.Printf("repo good add error %v \n", err)
			return
		}

		writeJsonResponse(w, http.StatusOK, good.ID.String())
	}
}
// getGood - get good
func (app *App) getGood() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()

		params := mux.Vars(request)
		uid := params["uid"]
		good, err := app.Repository.GetGood(ctx, uid )
		if err != nil {
			log.Printf("repo good get error %v \n", err)
			return
		}

		writeJsonResponse(w, http.StatusOK, good)
	}
}
// delGood - delete Good
func (app *App) delGood() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()

		params := mux.Vars(request)
		uid := params["uid"]
		uuID, _ := uu.FromString(uid)
		err := app.Repository.DelGood(ctx, uuid.UUID(uuID))
		if err != nil {
			log.Printf("repo good remove error %v \n", err)
			return
		}

		writeJsonResponse(w, http.StatusOK, "OK")
	}
}

//findGood - elastic search find
func (app *App) findGood() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer func() {
			// update Prom objects
		}()
		params := request.URL.Query()
		key := params["key"]
		//value := params["value"]
		goods, err := app.Repository.FindGood(ctx, key[0] )
		if err != nil {
			log.Printf("repo findgood get error %v \n", err)
			return
		}

		writeJsonResponse(w, http.StatusOK, goods)
	}
}