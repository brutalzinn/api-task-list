package oauth_api_server

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/brutalzinn/api-task-list/configs"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"github.com/go-session/session"
	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
)

var oauthserver *server.Server
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
func createOAuthServer() (*server.Server, error) {
	config := configs.GetConfig()
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	// token store
	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr: config.Redis.Host,
		DB:   config.Redis.Db,
	}))
	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	// manager.MapAccessGenerate(generates.NewAccessGenerate())
	clientStore := store.NewClientStore()
	clientStore.Set(idvar, &models.Client{
		ID:     idvar,
		Secret: secretvar,
		Domain: domainvar,
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
	return srv, nil
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
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
	store.Delete("LoggedInUserID")
	store.Save()
	w.Header().Set("Location", r.Form.Get("redirect_uri"))
	return
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
