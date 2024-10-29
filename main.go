package main

import (
	controller "PrayerService/controller"
	_ "PrayerService/docs"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//	@title			Prayer Service
//	@version		1.0
//	@description	This is a websocket server offering broadcast and user specific messaging.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host		prayer-service-495160257238.us-east4.run.app
// @BasePath	/
func main() {
	controller := controller.GetInstance()
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
	  AllowedOrigins:   []string{"https://*", "http://*", "ws://*", "wss://*"},
	  AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	  AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	  ExposedHeaders:   []string{"Link"},
	  AllowCredentials: false,
	  MaxAge:           300, 
	}))
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	router.Handle("/subscribe", controller.Auth(controller.Subscribe))
	router.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
		httpSwagger.BeforeScript(`const UrlMutatorPlugin = (system) => ({
			rootInjects: {
			  setScheme: (scheme) => {
				const jsonSpec = system.getState().toJSON().spec.json;
				const schemes = Array.isArray(scheme) ? scheme : [scheme];
				const newJsonSpec = Object.assign({}, jsonSpec, { schemes });
		  
				return system.specActions.updateJsonSpec(newJsonSpec);
			  },
			  setHost: (host) => {
				const jsonSpec = system.getState().toJSON().spec.json;
				const newJsonSpec = Object.assign({}, jsonSpec, { host });
		  
				return system.specActions.updateJsonSpec(newJsonSpec);
			  },
			  setBasePath: (basePath) => {
				const jsonSpec = system.getState().toJSON().spec.json;
				const newJsonSpec = Object.assign({}, jsonSpec, { basePath });
		  
				return system.specActions.updateJsonSpec(newJsonSpec);
			  }
			}
		  });`),
		httpSwagger.Plugins([]string{"UrlMutatorPlugin"}),
		httpSwagger.UIConfig(map[string]string{
			"showExtensions":        "true",
			"onComplete": fmt.Sprintf(`() => {
			window.ui.setScheme(['wss', 'https', 'http', 'ws']);
			window.ui.setHost('%s');
			window.ui.setBasePath('%s');
		}`, "prayer-service-495160257238.us-east4.run.app", "/"),
		}),
	))
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "moderator/dist"))
	FileServer(router, "/", filesDir)

	log.Println("Server started on port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalln("Server error:", err)
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
