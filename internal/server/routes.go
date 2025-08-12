package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

// var userTemplate = `<!DOCTYPE html><html><head><title>User</title></head><body>{{.Name}}</body></html>`

// Initialize the session store in an init() function or in your Server constructor.
// Remove top-level statements that are not declarations.

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte(os.Getenv("LOGIN_KEY")))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // example: 1 day
		HttpOnly: true,
		Secure:   false, // set to true in production
	}
	gothic.Store = store
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	// r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		// Clear Go session
		session, _ := store.Get(r, "auth-session")
		session.Options.MaxAge = -1 // expire immediately
		_ = session.Save(r, w)

		// Also clear gothic session if used
		gothicSess, _ := gothic.Store.Get(r, "gothic-session")
		gothicSess.Options.MaxAge = -1
		_ = gothicSess.Save(r, w)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"message":"Logged out"}`)
	})

	// r.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
	// 	// try to get the user without re-authenticating
	// 	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
	// 		t, _ := template.New("foo").Parse(userTemplate)
	// 		t.Execute(res, gothUser)
	// 	} else {
	// 		gothic.BeginAuthHandler(res, req)
	// 	}
	// })
	r.Get("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")

		if _, ok := session.Values["userID"]; ok {
			// Already logged in → go straight to profile
			http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
			return
		}

		// Not logged in → start OAuth
		gothic.BeginAuthHandler(w, r)
	})

	r.Get("/api/auth/status", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		if _, ok := session.Values["userID"]; ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"authenticated": true, "name": "%s"}`, session.Values["name"])
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"authenticated": false}`)
		}
	})

	r.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")

		r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Clear gothic session so it doesn't try to store the big payload
		sess, _ := gothic.Store.Get(r, "gothic-session")
		sess.Values = make(map[interface{}]interface{})
		_ = sess.Save(r, w)

		// Store only minimal info
		session, _ := store.Get(r, "auth-session")
		session.Values["userID"] = user.UserID
		session.Values["email"] = user.Email
		session.Values["name"] = user.Name
		_ = session.Save(r, w)

		// Redirect to a protected page
		redirectURL := "http://localhost:5173/dashboard"
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	})

	// Protected route
	r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		if _, ok := session.Values["userID"]; ok {
			fmt.Fprintf(w, "Welcome, %s! Email: %s", session.Values["name"], session.Values["email"])
		} else {
			http.Redirect(w, r, "/auth/google", http.StatusTemporaryRedirect)
		}
	})

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

// func (s *Server) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {

// 	provider := chi.URLParam(r, "provider")

// 	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

// 	user, err := gothic.CompleteUserAuth(w, r)
// 	if err != nil {
// 		fmt.Fprintln(w, err)
// 		return
// 	}
// 	fmt.Println(user)

// 	http.Redirect(w, r, "http://localhost:5173/dashboard?name="+url.QueryEscape(user.Name), http.StatusFound)
// }
