package authentication_service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/brutalzinn/api-task-list/db"
	pg "github.com/brutalzinn/go-oauth2-pg"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"

	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

var oauthserver *server.Server
var oauthClientManager *pg.ClientStore
var oauthTokenStore *pg.TokenStore

func InitOauthServer() {
	serverConfig, err := createOAuthServer()
	if err != nil {
		log.Println("Internal Error:", err.Error())
	}
	log.Printf("OAuth client Auth endpoint to %s:%s", "http://localhost", "/oauth/authorize")
	log.Printf("OAuth client Token endpoint to %s:%s", "http://localhost", "/oauth/token")
	oauthserver = serverConfig
}
func GetOauthServer() *server.Server {
	return oauthserver
}
func GetClientStore() *pg.ClientStore {
	return oauthClientManager
}
func GetTokenStore() *pg.TokenStore {
	return oauthTokenStore
}
func createOAuthServer() (*server.Server, error) {
	dbUri := db.GetConnectionUri()
	pgxConn, err := pgx.Connect(context.TODO(), dbUri)
	if err != nil {
		panic(err)
	}
	manager := manage.NewDefaultManager()
	adapter := pgx4adapter.NewConn(pgxConn)
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(5*time.Minute))
	defer tokenStore.Close()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	clientStore, _ := pg.NewClientStore(adapter)
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)
	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(UserAuthorizeHandler)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
		return
	})
	oauthClientManager = clientStore
	oauthTokenStore = tokenStore
	return srv, nil
}

func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedUserId")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/oauth/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedUserId")
	store.Save()
	return
}
