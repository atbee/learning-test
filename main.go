package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/codeforpublic/morchana-static-qr-code-api/aid"
	"github.com/codeforpublic/morchana-static-qr-code-api/internal/auth"
	"github.com/codeforpublic/morchana-static-qr-code-api/internal/middleware"
	"github.com/codeforpublic/morchana-static-qr-code-api/login"
	"github.com/codeforpublic/morchana-static-qr-code-api/qrcode"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var (
	buildcommit = "development"
	buildtime   = time.Now().Format(time.RFC3339)
)

func init() {
	initConfig()
}

func main() {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	// all routes required headers
	r.Use(middleware.Headers(viper.GetString("cors.allow_origin")))

	r.HandleFunc("/healths", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})
	r.HandleFunc("/versions", versionHandler)

	// additional headers for preflight request
	r.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Max-Age", viper.GetString("cors.max_age"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	})

	router := r.NewRoute().Subrouter()
	if viper.GetBool("app.secure") {
		router.Use(auth.Protect([]byte(viper.GetString("auth.token.access.secret"))))
	}

	db, err := newDBClient(viper.GetString("db.conn.string"))
	if err != nil {
		log.Fatal(err)
	}

	baseUrl, err := url.Parse(viper.GetString("otp.url"))
	if err != nil {
		log.Fatal(err)
		return
	}

	params := url.Values{}
	params.Add("user", viper.GetString("otp.user"))
	params.Add("pass", viper.GetString("otp.password"))
	params.Add("from", viper.GetString("otp.from"))
	params.Add("msg", viper.GetString("otp.message"))

	router.HandleFunc("/logins/{subr}", login.LoginOTP(httpClient, baseUrl, params)).Methods(http.MethodOptions, http.MethodGet)

	router.HandleFunc("/registerDevice", aid.AnonymousID(aid.StoreAnonymousID(db, viper.GetString("db.table.anonymousID")))).Methods(http.MethodOptions, http.MethodPost)
	router.HandleFunc("/qr", qrcode.Generate(viper.GetString("qr.signature"))).Methods(http.MethodOptions, http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:" + viper.GetString("port"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf(
		"Starting the service...\ncommit: %s, build time: %s, release: %s\n",
		buildcommit, buildtime, viper.GetString("app.version"),
	)
	log.Printf("serve on %s\n", ":"+viper.GetString("port"))

	go func() {
		log.Printf("%s", srv.ListenAndServe())
	}()

	gracefulshutdown(srv)
}

func gracefulshutdown(srv *http.Server) {
	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("%s", err.Error())
	}
}

func newDBClient(connStr string) (*sql.DB, error) {
	connector, err := mssql.NewConnector(connStr)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)
	return db, nil
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: viper.GetInt("http.maxconns"),
		MaxConnsPerHost:     viper.GetInt("http.maxconns"),
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: viper.GetBool("http.insecureskipverify")},
	},
	Timeout: viper.GetDuration("http.timeout"),
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"buildTime": buildtime,
		"commit":    buildcommit,
		"version":   viper.GetString("app.version"),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func initConfig() {
	viper.SetDefault("port", "1323")
	viper.SetDefault("app.secure", false)
	viper.SetDefault("app.version", "beta")
	viper.SetDefault("auth.token.access.secret", "jwtsecretkey")
	viper.SetDefault("cors.allow_origin", "*")
	viper.SetDefault("cors.max_age", "3600")
	viper.SetDefault("http.maxconns", 100)
	viper.SetDefault("http.insecureskipverify", true)
	viper.SetDefault("http.timeout", "5s")

	viper.SetDefault("otp.url", "https://ohgikdu5ed.execute-api.ap-southeast-1.amazonaws.com/smsgw-api")
	viper.SetDefault("otp.user", "morchana2")
	viper.SetDefault("otp.password", "Y8adfQzJfwUKGwUY")
	viper.SetDefault("otp.from", "Morchana")
	viper.SetDefault("otp.message", "ทดสอบ ข้อความ ภาษาไทย")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
