package oauth_api_server

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/brutalzinn/api-task-list/db"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

var oauthserver *server.Server
var oauthClientManager *pg.ClientStore
var oauthTokenStore *pg.TokenStore

var (
	idvar     string
	secretvar string
	domainvar string
)

func InitOauthServer() {
	flag.Parse()
	flag.StringVar(&idvar, "i", "222222", "The client id being passed in")
	flag.StringVar(&secretvar, "s", "22222222", "The client secret being passed in")
	flag.StringVar(&domainvar, "r", "https://oauth.pstmn.io/v1/callback", "The domain of the redirect url")
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
	pgxConn, _ := pgx.ConnectConfig(context.TODO(), db.GetConnectionAdapter())
	manager := manage.NewDefaultManager()

	adapter := pgx4adapter.NewConn(pgxConn)
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	defer tokenStore.Close()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	clientStore, _ := pg.NewClientStore(adapter)

	// generate jwt access token
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
	})

	oauthClientManager = clientStore
	oauthTokenStore = tokenStore
	return srv, nil
}

func DumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
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
