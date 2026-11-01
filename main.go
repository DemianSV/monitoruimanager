package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"monitoruimanager/tlsutil"

	"github.com/DemianSV/monitorclient"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gocql/gocql"
)

//go:embed all:assets all:web
var Assets embed.FS

//go:embed index.html
var Index []byte

const (
	// Version
	Version string = "2.0.0"
	// LogFile
	LogFile string = "monitoruimanager.log"
	// ConfigFile
	ConfigFile string = "monitoruimanager.json"
)

var ConsistencyRead gocql.Consistency

// FlagLog File logistics type - Stdout
var FlagLog string
var Session *gocql.Session
var TokenAuth *jwtauth.JWTAuth
var TokenSecret string

type ResponseCountT struct {
	Count int `json:"count"`
}

type DBConnectObserverT struct {
	Host    string `json:"host"`
	DC      string `json:"dc"`
	Rack    string `json:"rack"`
	Version string `json:"version"`
	Stat    string `json:"stat"`
	Up      bool   `json:"up"`
}
type DBConnectObserverAT []DBConnectObserverT

var (
	dbStatusMu sync.RWMutex
	DBStatusA  DBConnectObserverAT
)

var (
	httpDuration = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
		},
		[]string{"path", "status", "method", "func"},
	)
)

func init() {
	flag.StringVar(&FlagLog, "log", "stdout", "[-log file, stdout (default)]")
}

func setupApplication() {
	if FlagLog == "file" {
		fileLog, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			fmt.Println("Ошибка открытия LOG файла!")
			os.Exit(-1)
		}
		log.SetOutput(fileLog)
	}

	err := AppConfig.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации (", err.Error(), ")!")
		return
	}
	TokenSecret = AppConfig.JWT.SECRET
	if TokenSecret == "" {
		var err error
		TokenSecret, err = generateRandomToken(32)
		if err != nil {
			log.Fatal("Не удалось создать временный секрет JWT (", err, ")!")
		}
		log.Print("MONITOR_JWT_SECRET is not configured; JWTs will not remain valid across restarts or instances.")
	} else if len(TokenSecret) < 32 {
		log.Fatal("MONITOR_JWT_SECRET must be at least 32 bytes long.")
	}
	TokenAuth = jwtauth.New("HS256", []byte(TokenSecret), nil)

	monitorclient.Conf.SpaceID = AppConfig.UIMANAGER.SPACEID
	monitorclient.Conf.ModuleID = AppConfig.UIMANAGER.MODULEID
	monitorclient.Conf.DataManagerURL = AppConfig.UIMANAGER.DATAMANAGERURL
	monitorclient.Conf.ClientInSecureSkipVerify = AppConfig.HTTP.SKIPVERIFY
	monitorclient.Conf.DataManagerTimeOut = 1
}

func done() {
}

func (o *DBConnectObserverT) ObserveConnect(observed gocql.ObservedConnect) {
	if observed.Err != nil {
		log.Printf(
			"Connect observed: host=%s start=%v end=%v err=%v",
			observed.Host.ConnectAddress(),
			observed.Start,
			observed.End,
			observed.Err,
		)
	} else {
		log.Printf(
			"Connect observed: host=%s start=%v end=%v",
			observed.Host.ConnectAddress(),
			observed.Start,
			observed.End,
		)
	}

	status := DBConnectObserverT{
		Host:    observed.Host.ConnectAddress().String(),
		DC:      observed.Host.DataCenter(),
		Rack:    observed.Host.Rack(),
		Version: observed.Host.Version().String(),
		Stat:    observed.Host.State().String(),
		Up:      observed.Host.IsUp(),
	}
	dbStatusMu.Lock()
	defer dbStatusMu.Unlock()
	for i, item := range DBStatusA {
		if item.Host == status.Host {
			DBStatusA[i] = status
			return
		}
	}
	DBStatusA = append(DBStatusA, status)
}

func main() {
	defer done()
	defer log.Print("The service is stopped!")

	flag.Parse()
	setupApplication()

	envFlag := runtime.GOMAXPROCS(runtime.NumCPU())
	if envFlag > -1 {
		log.Print("GOMAXPROCS = ", runtime.NumCPU())
	} else {
		log.Print("GOMAXPROCS is default!")
	}

	// Connect to DB cluster
	cluster := gocql.NewCluster(AppConfig.DB.HOSTS...)
	cluster.Timeout = time.Duration(AppConfig.DB.TIMEOUT) * time.Millisecond
	cluster.Keyspace = AppConfig.DB.KEYSPACE

	AppConfig.DB.CONSISTENCY = strings.ToUpper(AppConfig.DB.CONSISTENCY)
	switch AppConfig.DB.CONSISTENCY {
	case "QUORUM":
		cluster.Consistency = gocql.Quorum
	case "ANY":
		cluster.Consistency = gocql.Any
	case "ONE":
		cluster.Consistency = gocql.One
	case "TWO":
		cluster.Consistency = gocql.Two
	case "THREE":
		cluster.Consistency = gocql.Three
	case "ALL":
		cluster.Consistency = gocql.All
	case "LOCALQUORUM":
		cluster.Consistency = gocql.LocalQuorum
	case "EACHQUORUM":
		cluster.Consistency = gocql.EachQuorum
	case "LOCALONE":
		cluster.Consistency = gocql.LocalOne
	default:
		cluster.Consistency = gocql.Quorum
	}
	AppConfig.DB.CONSISTENCYREAD = strings.ToUpper(AppConfig.DB.CONSISTENCYREAD)
	switch AppConfig.DB.CONSISTENCYREAD {
	case "QUORUM":
		ConsistencyRead = gocql.Quorum
	case "ANY":
		ConsistencyRead = gocql.Any
	case "ONE":
		ConsistencyRead = gocql.One
	case "TWO":
		ConsistencyRead = gocql.Two
	case "THREE":
		ConsistencyRead = gocql.Three
	case "ALL":
		ConsistencyRead = gocql.All
	case "LOCALQUORUM":
		ConsistencyRead = gocql.LocalQuorum
	case "EACHQUORUM":
		ConsistencyRead = gocql.EachQuorum
	case "LOCALONE":
		ConsistencyRead = gocql.LocalOne
	default:
		ConsistencyRead = gocql.Quorum
	}

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: AppConfig.DB.USERNAME,
		Password: AppConfig.DB.PASSWORD,
	}

	if AppConfig.DB.TLS {
		crt, err := tlsutil.LoadX509KeyPairWithPassword(AppConfig.TLS.CERTPATH, AppConfig.TLS.KEYPATH, []byte(AppConfig.TLS.KEYPASSWORD))
		if err != nil {
			log.Fatal("Failed to download the certificate (", err, ")!")
			return
		}
		var crtPool *x509.CertPool
		if AppConfig.TLS.CAPATH == "" {
			crtPool = nil
		} else {
			crtPool = x509.NewCertPool()
			if crtCA, err := os.ReadFile(AppConfig.TLS.CAPATH); err != nil {
				log.Fatal("Failed to download CA certificate (", err, ")!")
			} else if ok := crtPool.AppendCertsFromPEM(crtCA); !ok {
				log.Fatal("It was not possible to apply CA certificate!")
			}
		}

		cluster.SslOpts = &gocql.SslOptions{
			Config: &tls.Config{
				RootCAs:            crtPool,
				Certificates:       []tls.Certificate{crt},
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: AppConfig.DB.SKIPVERIFY,
			},
		}
	}

	cluster.ConnectObserver = &DBConnectObserverT{}

	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("It was not possible to establish a connection with the database (", err, ")!")
	} else {
		log.Print("The connection to the database is established... OK")
	}
	defer Session.Close()

	r := chi.NewRouter() // Initialization Gochi Router

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger) // Expanded logging
	r.Use(middleware.URLFormat)

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  allowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(httprate.LimitBy(
		AppConfig.HTTP.RL, // Requests
		1*time.Second,     // Per duration
		func(r *http.Request) (string, error) {
			return r.Header.Get("api_key"), nil
		},
	))

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verify(TokenAuth, jwtauth.TokenFromHeader))
		r.Use(jwtauth.Authenticator(TokenAuth))
		r.Use(requireAccessToken)

		r.Route("/api/v{version}/version", func(r chi.Router) {
			r.Get("/", getVersion) // Version
		})
		r.Route("/api/v{version}/userinfo", func(r chi.Router) {
			r.Get("/", getUserInfo) // Information on an authorized user
		})
		r.Route("/api/v{version}/password", func(r chi.Router) {
			r.With(requireSecureCredentialTransport).Put("/", putPassword) // Password change
		})
	})

	// Administrator branch
	r.Route("/api/v{version}/admin", func(r chi.Router) {
		r.Use(jwtauth.Verify(TokenAuth, jwtauth.TokenFromHeader))
		r.Use(jwtauth.Authenticator(TokenAuth))
		r.Use(requireAccessToken)

		r.Route("/user", func(r chi.Router) {
			r.Get("/count", getUserCount)
			r.Put("/select", putUserSelect)
			r.With(requireSecureCredentialTransport).Put("/update", putUserUpdate)
			r.Put("/remove", putUserRemove)
			r.With(requireSecureCredentialTransport).Put("/create", putUserCreate)
		})

		r.Route("/module", func(r chi.Router) {
			r.Get("/count", getModuleCount)
			r.Put("/select", putModuleSelect)
			r.Put("/update", putModuleUpdate)
			r.Put("/remove", putModuleRemove)
			r.Put("/create", putModuleCreate)
		})

		r.Route("/task", func(r chi.Router) {
			r.Get("/count", getTaskCount)
			r.Put("/select", putTaskSelect)
			r.Put("/update", putTaskUpdate)
			r.Put("/remove", putTaskRemove)
			r.Put("/create", putTaskCreate)
		})

		r.Route("/space", func(r chi.Router) {
			r.Get("/count", getSpaceCount)
			r.Put("/select", putSpaceSelect)
			r.Put("/update", putSpaceUpdate)
			r.Put("/remove", putSpaceRemove)
			r.Put("/create", putSpaceCreate)
		})

		r.Route("/problem", func(r chi.Router) {
			r.Get("/count", getProblemCount)
			r.Get("/select", getProblemSelect)
			// r.Put("/update", putProblemUpdate)
		})

		r.Route("/rawdata", func(r chi.Router) {
			r.Put("/select", putRawDataSelect)
			r.Put("/metric/select", putMetricDataSelect)
		})

		r.Route("/rawtext", func(r chi.Router) {
			r.Put("/select", putRawTextSelect)
			r.Put("/metric/select", putMetricTextSelect)
		})

		r.Route("/chartdata", func(r chi.Router) {
			r.Put("/", putChartData)
		})

		r.Route("/database", func(r chi.Router) {
			r.Get("/status", getDataBaseStatus)
		})

		r.Route("/datamanager", func(r chi.Router) {
			r.Get("/status", getDataManagerStatus)
		})

		r.Route("/dashboard", func(r chi.Router) {
			r.Get("/count", getDashboardCount)
			r.Get("/select", getDashboardSelect)
			r.Put("/data", putChartDashboardData)
			r.Put("/create", putDashboardCreate)
			r.Put("/remove", putDashboardRemove)
			r.Put("/edit", putDashboardEdit)
		})

		r.Route("/notification", func(r chi.Router) {
			r.Put("/select", putNotificationSelect)
			r.Put("/create", putNotificationCreate)
			r.Put("/remove", putNotificationRemove)
			r.Put("/edit", putNotificationEdit)
		})
	})

	// Requests to UI
	r.Get("/", webTmpl)
	r.Get("/user", webTmpl)
	r.Get("/module", webTmpl)
	r.Get("/task", webTmpl)
	r.Get("/space", webTmpl)
	r.Get("/about", webTmpl)
	r.Get("/rawtext", webTmpl)
	r.Get("/rawdata", webTmpl)
	r.Get("/dashboard", webTmpl)
	r.Get("/database", webTmpl)
	r.Get("/metric", webTmpl)
	r.Get("/problem", webTmpl)
	r.Get("/notification", webTmpl)

	r.Post("/auth", loginWithRefreshToken)
	r.Post("/auth/refresh", refreshAccessToken)
	r.Post("/auth/logout", logoutWithRefreshToken)

	// Static content for UI
	assetsFS, errFS := fs.Sub(Assets, "assets")
	if errFS != nil {
		log.Fatal("It was not possible to connect static content!")
	} else {
		fileServer(r, "/assets", http.FS(assetsFS))
	}

	webFS, errWebFS := fs.Sub(Assets, "web")
	if errWebFS != nil {
		log.Fatal("It was not possible to connect PWA content!")
	} else {
		webServer := http.FileServer(http.FS(webFS))
		r.Get("/{file}", func(w http.ResponseWriter, r *http.Request) {
			if chi.URLParam(r, "file") == "sw.js" {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			}
			webServer.ServeHTTP(w, r)
		})
	}

	r.Mount("/metrics", promhttp.Handler()) // Prometheus metrics

	if AppConfig.HTTP.TLS { // HTTPS connection
		// crt, err := tls.LoadX509KeyPair(AppConfig.TLS.CERTPATH, AppConfig.TLS.KEYPATH)
		crt, err := tlsutil.LoadX509KeyPairWithPassword(AppConfig.TLS.CERTPATH, AppConfig.TLS.KEYPATH, []byte(AppConfig.TLS.KEYPASSWORD))
		if err != nil {
			log.Fatal("Failed to download the certificate (", err, ")!")
			return
		}
		var crtPool *x509.CertPool
		if AppConfig.TLS.CAPATH == "" {
			crtPool = nil
		} else {
			crtPool = x509.NewCertPool()
			if crtCA, err := os.ReadFile(AppConfig.TLS.CAPATH); err != nil {
				log.Fatal("Failed to download CA certificate (", err, ")!")
			} else if ok := crtPool.AppendCertsFromPEM(crtCA); !ok {
				log.Fatal("It was not possible to apply CA certificate!")
			}
		}
		tlsConfig := &tls.Config{
			RootCAs:            crtPool,
			Certificates:       []tls.Certificate{crt},
			InsecureSkipVerify: AppConfig.HTTP.SKIPVERIFY,
			MinVersion:         tls.VersionTLS12,
		}

		httpServer := &http.Server{Addr: AppConfig.HTTP.HOST + ":" + AppConfig.HTTP.PORT,
			Handler:      r,
			ReadTimeout:  time.Duration(AppConfig.HTTP.READTIMEOUT) * time.Second,
			WriteTimeout: time.Duration(AppConfig.HTTP.WRITETIMEOUT) * time.Second,
			TLSConfig:    tlsConfig}
		log.Print("Сервис UIManager запущен на порту ", AppConfig.HTTP.PORT, "(HTTPS)... ОК")
		httpServer.ListenAndServeTLS("", "")
	} else { // HTTP connection
		httpServer := &http.Server{Addr: AppConfig.HTTP.HOST + ":" + AppConfig.HTTP.PORT,
			Handler:      r,
			ReadTimeout:  time.Duration(AppConfig.HTTP.READTIMEOUT) * time.Second,
			WriteTimeout: time.Duration(AppConfig.HTTP.WRITETIMEOUT) * time.Second}
		log.Print("Сервис UIManager запущен на порту ", AppConfig.HTTP.PORT, "(HTTP)... ОК")
		httpServer.ListenAndServe()
	}
}

func webTmpl(w http.ResponseWriter, r *http.Request) {
	uaClient := r.Header.Get("User-Agent")
	raClient := r.Header.Get("X-Real-Ip")
	if raClient == "" {
		raClient = r.Header.Get("X-Forwarded-For")
	}
	if raClient == "" {
		raClient = r.RemoteAddr
	}
	if raClient == "" {
		raClient = "Address not determined"
	}
	log.Print("IP: ", raClient, " UserAgent: ", uaClient)

	w.Header().Set("Cache-control", "no-cache, must-revalidate, private, no-store, s-maxage=0, max-age=0, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "-10d")

	w.Write(Index)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	_, claims, _ := jwtauth.FromContext(r.Context())
	log.Printf("User %v requested the API version!", claims["username"])

	type VersionJSON struct {
		Name    string
		Version string
	}

	versionAPI := chi.URLParam(r, "version")
	if versionAPI == "1" {
		version := VersionJSON{
			Name:    "Πολυξένη UIManager",
			Version: Version,
		}
		jsonAPI, err := json.Marshal(version)
		if err != nil {
			log.Print("JSON ERROR (" + err.Error() + ")!")

			duration := time.Since(start).Seconds()
			httpDuration.WithLabelValues("/version", "500", "GET", "getVersion").Set(duration)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues("/version", "200", "GET", "getVersion").Set(duration)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonAPI)
	}
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func authenticateUser(userName string, password string) (jwt.MapClaims, error) {
	if len(password) > maximumPasswordBytes {
		return nil, errors.New("password exceeds maximum supported length")
	}
	var userID string
	var status int
	var passwordDB string
	var role string
	var ownerID string
	err := Session.Query(
		`SELECT user_id, status, password, role, owner_id FROM user WHERE user_name = ?`,
		userName,
	).WithContext(context.Background()).Consistency(ConsistencyRead).Scan(&userID, &status, &passwordDB, &role, &ownerID)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			log.Printf("User %s ACCESS DENIED!", userName)
			go monitorclient.Log("audit", "User "+userName+" ACCESS DENIED")
			return nil, errors.New("ACCESS DENIED")
		}
		return nil, err
	}
	if status != 1 {
		log.Printf("User %s ACCESS DENIED!", userName)
		go monitorclient.Log("audit", "User "+userName+" ACCESS DENIED")
		return nil, errors.New("ACCESS DENIED")
	}

	verified, needsUpgrade, err := verifyPassword(passwordDB, password)
	if err != nil {
		log.Printf("Invalid password hash for user %s: %v", userName, err)
		return nil, err
	}
	if !verified {
		log.Printf("User %s ACCESS DENIED!", userName)
		go monitorclient.Log("audit", "User "+userName+" ACCESS DENIED")
		return nil, errors.New("ACCESS DENIED")
	}
	if needsUpgrade {
		upgradedHash, err := hashPassword(password)
		if err != nil {
			return nil, err
		}
		var previousPassword string
		applied, err := Session.Query(
			`UPDATE user SET password = ? WHERE user_name = ? IF password = ?`,
			upgradedHash, userName, passwordDB,
		).WithContext(context.Background()).ScanCAS(&previousPassword)
		if err != nil {
			return nil, fmt.Errorf("upgrade legacy password hash: %w", err)
		}
		if !applied {
			return nil, errors.New("password changed during authentication")
		}
	}

	return jwt.MapClaims{"username": userName, "userid": userID, "role": role, "ownerid": ownerID}, nil
}

func recordSuccessfulLogin(userName string) error {
	log.Print("The user " + userName + " has entered the system!")
	if err := monitorclient.Log("audit", "User "+userName+" entered the system"); err != nil {
		log.Print("It was not possible to record the audit data (", err, ")!")
		log.Print("Access error (", err.Error(), ")!")
		return err
	}

	if err := Session.Query(`UPDATE user SET login_time = toTimestamp(now()) WHERE user_name = ?`, userName).WithContext(context.Background()).Exec(); err != nil {
		log.Print("Failed to update the date of the user entry ", userName, "!")
		go monitorclient.Log("log", "Failed to update the date of the user entry "+userName)
	}
	return nil
}

func allowOriginFunc(r *http.Request, origin string) bool {
	return false
}

func putRawTextSelect(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	_, claims, _ := jwtauth.FromContext(r.Context())
	log.Printf("User %v requested raw_text data!", claims["username"])
	go monitorclient.Log("log", fmt.Sprintf("User %v requested raw_text data", claims["username"]))

	type ResponseRawTextDataT struct {
		SpaceID    string `json:"spaceid"`
		Metric     string `json:"metric"`
		CreateTime int64  `json:"createtime"`
		EventTime  int64  `json:"eventtime"`
		Value      string `json:"value"`
		Object     string `json:"object"`
		SpaceDesc  string `json:"spacedesc"`
	}
	type ResponseRawTextDataAT []ResponseRawTextDataT
	type ResponseRawTextT struct {
		Data        ResponseRawTextDataAT `json:"data"`
		CurrentPage string                `json:"current"`
		NextPage    string                `json:"next"`
	}

	type RequestT struct {
		SpaceID     string `json:"spaceid"`
		Metric      string `json:"metric"`
		CurrentPage string `json:"current"`
		PageSize    int    `json:"pagesize"`
	}

	var spaceDesc string

	if claims["role"] == "superadmin" || claims["role"] == "admin" || claims["role"] == "user" {
		versionAPI := chi.URLParam(r, "version")
		if versionAPI == "1" {
			b, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawtext/select", "500", "PUT", "putRawTextSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var request RequestT
			if err := json.Unmarshal(b, &request); err != nil {
				log.Print("JSON UNMARSHAL ERROR (" + err.Error() + ")!")
				go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawtext/select", "500", "PUT", "putRawTextSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			pageSize := max(request.PageSize, 0)

			var responseRawTextData ResponseRawTextDataT
			var responseRawTextDataA ResponseRawTextDataAT
			var pageState []byte

			if request.CurrentPage != "" {
				var err error
				pageState, err = base64.StdEncoding.DecodeString(request.CurrentPage)
				if err != nil {
					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/rawdata/select", "500", "PUT", "putRawDataSelect").Set(duration)

					log.Print("CurrentPage decode ERROR (" + err.Error() + ")!")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			ctx := context.Background()
			var scanner gocql.Scanner
			var iter *gocql.Iter

			scannerSpace := Session.Query(`SELECT description FROM space WHERE id = ?`, request.SpaceID).WithContext(ctx).Consistency(ConsistencyRead).Iter().Scanner()
			for scannerSpace.Next() {
				err := scannerSpace.Scan(&spaceDesc)
				if err != nil {
					log.Print(err)
				}
			}

			if request.Metric != "" {
				iter = Session.Query(`SELECT space_id, metric, create_time, event_time, value, object FROM raw_text WHERE space_id = ? AND metric = ? AND event_time > toUnixTimestamp(now()) - 2592000000`, request.SpaceID, request.Metric).WithContext(ctx).Consistency(ConsistencyRead).PageSize(pageSize).PageState(pageState).Iter()
			} else {
				log.Print("Sample metric not defined")
				go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawtext/select", "500", "PUT", "putRawTextSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var responseRawText ResponseRawTextT

			if request.PageSize > 0 {
				nextPageState := iter.PageState()
				if len(nextPageState) > 0 {
					responseRawText.NextPage = base64.StdEncoding.EncodeToString(nextPageState)
				} else {
					responseRawText.NextPage = ""
				}
				responseRawText.CurrentPage = request.CurrentPage

				scanner = iter.Scanner()
			}

			for scanner.Next() {
				var spaceID string
				var metric string
				var createTime int64
				var eventTime int64
				var value string
				var object string

				err := scanner.Scan(&spaceID, &metric, &createTime, &eventTime, &value, &object)
				if err != nil {
					go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/rawtext/select", "500", "PUT", "putRawTextSelect").Set(duration)

					w.WriteHeader(http.StatusInternalServerError)
					return
				} else {
					responseRawTextData.SpaceID = spaceID
					responseRawTextData.Metric = metric
					responseRawTextData.CreateTime = createTime
					responseRawTextData.EventTime = eventTime
					responseRawTextData.Value = value
					responseRawTextData.Object = object
					responseRawTextData.SpaceDesc = spaceDesc

					responseRawTextDataA = append(responseRawTextDataA, responseRawTextData)
				}
			}

			if err := scanner.Err(); err != nil {
				log.Print(err)
				go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawtext/select", "500", "PUT", "putRawTextSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			responseRawText.Data = responseRawTextDataA

			responseJSON, err := json.Marshal(responseRawText)
			if err != nil {
				log.Print("JSON MARSHAL ERROR (" + err.Error() + ")!")
				go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawtext/select", "500", "PUT", "putRawTextSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
			} else {
				go monitorclient.Metric("httpstatus", float32(http.StatusOK))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawtext/select", "200", "PUT", "putRawTextSelect").Set(duration)

				w.Header().Set("Content-Type", "application/json")
				w.Write(responseJSON)
				return
			}
		} else {
			log.Print("Failed to choose the data!")
			go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

			duration := time.Since(start).Seconds()
			httpDuration.WithLabelValues("/rawtext/select", "500", "PUT", "putRawTextSelect").Set(duration)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues("/rawtext/select", "403", "PUT", "putRawTextSelect").Set(duration)

		w.WriteHeader(http.StatusForbidden)
		return
	}
}

func putRawDataSelect(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	_, claims, _ := jwtauth.FromContext(r.Context())
	log.Printf("User %v requested raw_data data!", claims["username"])
	go monitorclient.Log("log", fmt.Sprintf("User %v requested raw_data data", claims["username"]))

	type ResponseRawDataDataT struct {
		SpaceID    string            `json:"spaceid"`
		Metric     string            `json:"metric"`
		CreateTime int64             `json:"createtime"`
		EventTime  int64             `json:"eventtime"`
		Value      float32           `json:"value"`
		Object     string            `json:"object"`
		SpaceDesc  string            `json:"spacedesc"`
		Labels     map[string]string `json:"labels"`
	}
	type ResponseRawDataDataAT []ResponseRawDataDataT
	type ResponseRawDataT struct {
		Data        ResponseRawDataDataAT `json:"data"`
		CurrentPage string                `json:"current"`
		NextPage    string                `json:"next"`
	}

	type RequestT struct {
		SpaceID     string `json:"spaceid"`
		Metric      string `json:"metric"`
		CurrentPage string `json:"current"`
		PageSize    int    `json:"pagesize"`
	}

	var spaceDesc string

	if claims["role"] == "superadmin" || claims["role"] == "admin" || claims["role"] == "user" {
		versionAPI := chi.URLParam(r, "version")
		if versionAPI == "1" {
			b, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/select", "500", "PUT", "putRawDataSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var request RequestT
			if err := json.Unmarshal(b, &request); err != nil {
				log.Print("JSON UNMARSHAL ERROR (" + err.Error() + ")!")
				go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/select", "500", "PUT", "putRawDataSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			pageSize := max(request.PageSize, 0)

			var responseRawDataData ResponseRawDataDataT
			var responseRawDataDataA ResponseRawDataDataAT
			var pageState []byte

			if request.CurrentPage != "" {
				var err error
				pageState, err = base64.StdEncoding.DecodeString(request.CurrentPage)
				if err != nil {
					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/rawdata/select", "500", "PUT", "putRawDataSelect").Set(duration)

					log.Print("CurrentPage decode ERROR (" + err.Error() + ")!")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			ctx := context.Background()
			var scanner gocql.Scanner
			var iter *gocql.Iter

			scannerSpace := Session.Query(`SELECT description FROM space WHERE id = ?`, request.SpaceID).WithContext(ctx).Consistency(ConsistencyRead).Iter().Scanner()
			for scannerSpace.Next() {
				err := scannerSpace.Scan(&spaceDesc)
				if err != nil {
					log.Print(err)
				}
			}

			if request.Metric != "" {
				iter = Session.Query(`SELECT space_id, metric, create_time, event_time, value, object, labels FROM raw_data WHERE space_id = ? AND metric = ? AND event_time > toUnixTimestamp(now()) - 2592000000`, request.SpaceID, request.Metric).WithContext(ctx).Consistency(ConsistencyRead).PageSize(pageSize).PageState(pageState).Iter()
			} else {
				log.Print("Sample metric not defined")
				go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/select", "500", "PUT", "putRawDataSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var responseRawData ResponseRawDataT

			if request.PageSize > 0 {
				nextPageState := iter.PageState()
				if len(nextPageState) > 0 {
					responseRawData.NextPage = base64.StdEncoding.EncodeToString(nextPageState)
				} else {
					responseRawData.NextPage = ""
				}
				responseRawData.CurrentPage = request.CurrentPage

				scanner = iter.Scanner()
			}

			for scanner.Next() {
				var spaceID string
				var metric string
				var createTime int64
				var eventTime int64
				var value float32
				var object string
				var labels map[string]string

				err := scanner.Scan(&spaceID, &metric, &createTime, &eventTime, &value, &object, &labels)
				if err != nil {
					go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/rawdata/select", "500", "PUT", "putRawDataSelect").Set(duration)

					w.WriteHeader(http.StatusInternalServerError)
					return
				} else {
					responseRawDataData.SpaceID = spaceID
					responseRawDataData.Metric = metric
					responseRawDataData.CreateTime = createTime
					responseRawDataData.EventTime = eventTime
					responseRawDataData.Value = value
					responseRawDataData.Object = object
					responseRawDataData.SpaceDesc = spaceDesc
					responseRawDataData.Labels = labels

					responseRawDataDataA = append(responseRawDataDataA, responseRawDataData)
				}
			}

			if err := scanner.Err(); err != nil {
				log.Print(err)
				go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/select", "500", "PUT", "putRawDataSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			responseRawData.Data = responseRawDataDataA

			responseJSON, err := json.Marshal(responseRawData)
			if err != nil {
				log.Print("JSON MARSHAL ERROR (" + err.Error() + ")!")
				go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/select", "500", "PUT", "putRawDataSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
			} else {
				go monitorclient.Metric("httpstatus", float32(http.StatusOK))
				w.Header().Set("Content-Type", "application/json")
				w.Write(responseJSON)

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/select", "200", "PUT", "putRawDataSelect").Set(duration)

				return
			}
		} else {
			log.Print("Failed to choose the data!")
			go monitorclient.Metric("httpstatus", float32(http.StatusInternalServerError))

			duration := time.Since(start).Seconds()
			httpDuration.WithLabelValues("/rawdata/select", "500", "PUT", "putRawDataSelect").Set(duration)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues("/rawdata/select", "403", "PUT", "putRawDataSelect").Set(duration)

		w.WriteHeader(http.StatusForbidden)
		return
	}
}

func putMetricTextSelect(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	_, claims, _ := jwtauth.FromContext(r.Context())
	log.Printf("User %v requested a meter list!", claims["username"])

	type MetricT struct {
		Metric string `json:"metric"`
	}
	type MetricTA []MetricT

	type SpaceT struct {
		SpaceID string `json:"spaceid"`
	}

	if claims["role"] == "superadmin" || claims["role"] == "admin" || claims["role"] == "user" {
		versionAPI := chi.URLParam(r, "version")
		if versionAPI == "1" {
			b, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var space SpaceT
			if err := json.Unmarshal(b, &space); err != nil {
				log.Print("JSON UNMARSHAL ERROR (" + err.Error() + ")!")

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawtext/metric/select", "500", "PUT", "putMetricTextSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var responseMetric MetricT
			var responseMetricA MetricTA
			ctx := context.Background()
			scanner := Session.Query(`SELECT metric FROM space_metric_text WHERE space_id = ?`, space.SpaceID).WithContext(ctx).Consistency(ConsistencyRead).Iter().Scanner()
			for scanner.Next() {
				var metric string

				err := scanner.Scan(&metric)
				if err != nil {
					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/rawtext/metric/select", "500", "PUT", "putMetricTextSelect").Set(duration)

					w.WriteHeader(http.StatusInternalServerError)
					return
				} else {
					responseMetric.Metric = metric

					responseMetricA = append(responseMetricA, responseMetric)
				}
			}
			if err := scanner.Err(); err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawtext/metric/select", "500", "PUT", "putMetricTextSelect").Set(duration)

				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			responseJSON, err := json.Marshal(responseMetricA)
			if err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawtext/metric/select", "500", "PUT", "putMetricTextSelect").Set(duration)

				log.Print("JSON MARSHAL ERROR (" + err.Error() + ")!")
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawtext/metric/select", "200", "PUT", "putMetricTextSelect").Set(duration)

				w.Header().Set("Content-Type", "application/json")
				w.Write(responseJSON)
				return
			}
			duration := time.Since(start).Seconds()
			httpDuration.WithLabelValues("/rawtext/metric/select", "500", "PUT", "putMetricTextSelect").Set(duration)

			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			duration := time.Since(start).Seconds()
			httpDuration.WithLabelValues("/rawtext/metric/select", "500", "PUT", "putMetricTextSelect").Set(duration)

			log.Print("It was not possible to choose metrics!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues("/rawtext/metric/select", "403", "PUT", "putMetricTextSelect").Set(duration)

		w.WriteHeader(http.StatusForbidden)
		return
	}
}

func putMetricDataSelect(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	_, claims, _ := jwtauth.FromContext(r.Context())
	log.Printf("User %v requested a meter list!", claims["username"])

	type MetricT struct {
		Metric string `json:"metric"`
	}
	type MetricTA []MetricT

	type SpaceT struct {
		SpaceID string `json:"spaceid"`
	}

	if claims["role"] == "superadmin" || claims["role"] == "admin" || claims["role"] == "user" {
		versionAPI := chi.URLParam(r, "version")
		if versionAPI == "1" {
			b, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/metric/select", "500", "PUT", "putMetricDataSelect").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var space SpaceT
			if err := json.Unmarshal(b, &space); err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/metric/select", "500", "PUT", "putMetricDataSelect").Set(duration)

				log.Print("JSON UNMARSHAL ERROR (" + err.Error() + ")!")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var responseMetric MetricT
			var responseMetricA MetricTA
			ctx := context.Background()
			scanner := Session.Query(`SELECT metric FROM space_metric_data WHERE space_id = ?`, space.SpaceID).WithContext(ctx).Consistency(ConsistencyRead).Iter().Scanner()
			for scanner.Next() {
				var metric string

				err := scanner.Scan(&metric)
				if err != nil {
					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/rawdata/metric/select", "500", "PUT", "putMetricDataSelect").Set(duration)

					w.WriteHeader(http.StatusInternalServerError)
					return
				} else {
					responseMetric.Metric = metric

					responseMetricA = append(responseMetricA, responseMetric)
				}
			}
			if err := scanner.Err(); err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/metric/select", "500", "PUT", "putMetricDataSelect").Set(duration)

				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			responseJSON, err := json.Marshal(responseMetricA)
			if err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/metric/select", "500", "PUT", "putMetricDataSelect").Set(duration)

				log.Print("JSON MARSHAL ERROR (" + err.Error() + ")!")
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/rawdata/metric/select", "200", "PUT", "putMetricDataSelect").Set(duration)

				w.Header().Set("Content-Type", "application/json")
				w.Write(responseJSON)
				return
			}
			duration := time.Since(start).Seconds()
			httpDuration.WithLabelValues("/rawdata/metric/select", "500", "PUT", "putMetricDataSelect").Set(duration)

			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			duration := time.Since(start).Seconds()
			httpDuration.WithLabelValues("/rawdata/metric/select", "500", "PUT", "putMetricDataSelect").Set(duration)

			log.Print("It was not possible to choose metrics!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues("/rawdata/metric/select", "403", "PUT", "putMetricDataSelect").Set(duration)

		w.WriteHeader(http.StatusForbidden)
		return
	}
}
