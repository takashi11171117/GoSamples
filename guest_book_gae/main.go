package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"github.com/BurntSushi/toml"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type config struct {
	Firebase firebaseConf
}

type firebaseConf struct {
	DatabaseURL   string
	ProjectID     string
	StorageBucket string
}

var configVar config

var (
	firebaseConfig = &firebase.Config{
		DatabaseURL:   configVar.Firebase.DatabaseURL,
		ProjectID:     configVar.Firebase.ProjectID,
		StorageBucket: configVar.Firebase.StorageBucket,
	}
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

type postStruct struct {
	Author  string
	UserID  string
	Message string
	Posted  time.Time
}

type templateParams struct {
	Notice  string
	Name    string
	Message string
	Posts   []postStruct
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	params := templateParams{}

	_, err := toml.DecodeFile("config.toml", &configVar)
	if err != nil {
		params.Notice = fmt.Sprintf("Setup file can not be read")
		indexTemplate.Execute(w, params)
		return
	}

	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	q := datastore.NewQuery("Post").Order("-Posted").Limit(20)
	if _, err := q.GetAll(ctx, &params.Posts); err != nil {
		log.Errorf(ctx, "Getting posts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't get latest posts. Refresh?"
		indexTemplate.Execute(w, params)
		return
	}

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
		return
	}

	message := r.FormValue("message")

	// Create a new Firebase App.
	app, err := firebase.NewApp(ctx, firebaseConfig)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Create a new authenticator for the app.
	auth, err := app.Auth(ctx)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Verify the token passed in by the user is valid.
	tok, err := auth.VerifyIDTokenAndCheckRevoked(ctx, r.FormValue("token"))
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Use the validated token to get the user's information.
	user, err := auth.GetUser(ctx, tok.UID)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}

	post := postStruct{
		UserID:  user.UID,
		Author:  user.DisplayName,
		Message: message,
		Posted:  time.Now(),
	}

	if post.Author == "" {
		post.Author = "Anonymous Gopher"
	}

	params.Name = post.Author

	if post.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		params.Notice = "No message provided"
		indexTemplate.Execute(w, params)
		return
	}

	key := datastore.NewIncompleteKey(ctx, "Post", nil)

	if _, err := datastore.Put(ctx, key, &post); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't add new post. Try again?"
		params.Message = post.Message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}

	params.Posts = append([]postStruct{post}, params.Posts...)
	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", post.Author)
	indexTemplate.Execute(w, params)
}
