package oauth_controller

import (
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
            <html>
                <body>
                    <form action="/oauth/auth" method="post">
                        <input type="text" name="email" placeholder="Email">
                        <input type="password" name="password" placeholder="Password">
                        <button type="submit">Log in</button>
                    </form>
                </body>
            </html>
        `)

}

func Authentication(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email != "user@example.com" || password != "password" {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// authCode, err := srv.NewAuthorizeRequest(r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Redirect to the OAuth2 provider's login page
	//url := authCode.GetAuthorizeURL()
	//http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
