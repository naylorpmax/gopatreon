# gopatreon

Lightweight package to expose Patreon API via gopkg.in/mxpv/patreon-go.v1.

## Examples

```golang
package example

import (
	"github.com/gorilla/mux"
	"github.com/naylorpmax/gopatreon"
)

type Callback struct {
	OAuth2Config *oauth2.Config
}

func (c *Callback) Handler(w http.ResponseWriter, r *http.Request) {
	go func() {
		code := r.FormValue("code")
		if code == "" {
			panic(err)
		}

		client, err := gopatreon.NewClient(r.Context(), code, c.OAuth2Config)
		if err != nil {
			panic(err)
		}

		service := gopatreon.NewService(client)

		userName, err := service.AuthenticateUser()
		if err != nil {
			panic(err)
		}

		welcomeMsg := map[string]string{
			"message": "welcome! you're logged in",
			"name":    userName,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(welcomeMsg)
		if err != nil {
			panic(err)
		}
	}()
}

type Authorize struct {
	OAuth2Config *oauth2.Config
}

func (a *Authorize) Handler(w http.ResponseWriter, r *http.Request) {
	go func() {
		url := l.OAuth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}()
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	authorize := &Authorize{OAuth2Config: cfg.OAuth2Config}
	router.Methods("GET").
		Path("/authorize").
		Handler(authorize.Handler)

	callback := &Callback{OAuth2Config: cfg.OAuth2Config}
	router.Methods("GET").
		Path("/callback").
		Handler(callback.Handler)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	panic(server.ListenAndServe())
}
```