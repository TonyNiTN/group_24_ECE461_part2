package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/19chonm/461_1_23/db"
	"github.com/19chonm/461_1_23/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func RunServer() {

	//create client to interact with google cloud storage
	const (
		projectid  = "trusted-package-registry"
		bucketname = "package-repos"
		tablename  = "packages"
	)

	client, err := db.NewBucketClient(context.Background(), projectid, bucketname)
	if err != nil {
		fmt.Println("failed to create GCS client!")
	}
	defer client.Close()

	firestoreClient, err := db.NewFirestoreClient(context.Background(), projectid, tablename)

	if err != nil {
		fmt.Println("failed to create Firestore client!")
	}
	defer firestoreClient.Close()

	//gin.SetMode(gin.ReleaseMode)

	//Initialize go gin router
	r := gin.Default()
	//r.LoadHTMLGlob("views/*")
	r.Use(CORSMiddleware())
	authRoutes := r.Group("/")
	authRoutes.Use(AuthMiddleware())
	{
		authRoutes.GET("/", func(c *gin.Context) {
			pageSize := 10
			pageToken := ""

			// Create a slice to hold the packages
			packages := make([]*db.Package, 0)

			// Fetch the packages from Firestore
			for {
				// Create a query to get the next page of packages
				q := firestoreClient.GetClient().Collection("packages").OrderBy("name", firestore.Asc)
				if pageToken != "" {
					q = q.StartAfter(pageToken)
				}
				q = q.Limit(pageSize)

				// Execute the query
				iter := q.Documents(firestoreClient.GetCtx())
				for {
					// Get the next document
					doc, err := iter.Next()
					if err == iterator.Done {
						break
					}
					if err != nil {
						c.JSON(http.StatusInternalServerError, "Error getting documents from the database")
						log.Fatalf("Failed to iterate Firestore documents: %v", err)
					}

					// Unmarshal the document into a Package struct
					var p *db.Package
					err = doc.DataTo(&p)
					if err != nil {
						log.Fatalf("Failed to unmarshal Firestore document: %v", err)
					}
					p.ID = doc.Ref.ID

					// Add the package to the slice
					packages = append(packages, p)
				}

				// Check if there are more pages
				if len(packages) >= pageSize {
					doc, err := iter.Next()
					if err == iterator.Done {
						break
					} else if err != nil {
						c.JSON(http.StatusInternalServerError, "Error getting documents from the database")
						log.Fatalf("Failed to iterate Firestore documents: %v", err)
					} else {
						pageToken = doc.Ref.ID
					}
				} else {
					break
				}
			}

			// c.HTML(http.StatusOK, "index.html", gin.H{
			// 	"packages":  packages,
			// 	"pageSize":  pageSize,
			// 	"pageToken": pageToken,
			// })

			c.JSON(http.StatusOK, packages)
		})
		authRoutes.POST("/packages/upload", func(c *gin.Context) {
			if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			}

			name := c.Request.FormValue("name")
			url := c.Request.FormValue("url")

			id := uuid.New().String()
			packageData := &db.Package{
				ID:   id,
				Name: name,
				URL:  url,
			}

			err = firestoreClient.UploadPackage(context.Background(), firestoreClient.GetClient(), packageData, id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "error uploading package")
			}

			file, header, err := c.Request.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, "error getting file from form")
				return
			}
			defer file.Close()

			if err := client.UploadFile(header.Filename, file, id); err != nil {
				c.JSON(http.StatusInternalServerError, "error uploading file")
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "package uploaded"})

		})
		authRoutes.GET("/packages/:id", func(c *gin.Context) {
			packageID := c.Param("id")
			packageInfo, err := firestoreClient.GetPackage(context.Background(), firestoreClient.GetClient(), packageID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, packageInfo)
		})
		authRoutes.GET("/packages/:id/score", func(c *gin.Context) {
			packageID := c.Param("id")
			packageInfo, err := firestoreClient.GetPackage(context.Background(), firestoreClient.GetClient(), packageID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			firestoreClient.ScorePackage(context.Background(), firestoreClient.GetClient(), packageInfo.URL, packageInfo)
			c.JSON(http.StatusOK, packageInfo)
		})
		authRoutes.GET("/packages/:id/download", func(c *gin.Context) {
			packageID := c.Param("id")
			packageInfo, err := firestoreClient.GetPackage(context.Background(), firestoreClient.GetClient(), packageID)
			if err != nil {
				logger.DebugMsg("error getting package info in Gin package download handler")
			}

			w := c.Writer
			r := c.Request

			// set the response headers to indicate a zip file
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Content-Disposition", "attachment; filename="+packageInfo.Name+".zip")

			// call the DownloadFile function to write the file content to the response body
			err = client.DownloadFile(packageID, w, r)
			if err != nil {
				logger.DebugMsg("error downloading file")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, "Download successful")
		})
		authRoutes.GET("/packages/search", func(c *gin.Context) {
			packageName := c.Query("name")
			searchResults, err := firestoreClient.SearchPackage(context.Background(), firestoreClient.GetClient(), packageName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, searchResults)
		})
		authRoutes.DELETE("/packages/:id/delete", func(c *gin.Context) {
			packageName := c.Param("id")
			if err := firestoreClient.DeletePackage(context.Background(), firestoreClient.GetClient(), packageName); err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "package deleted"})
		})
	}

	//LOGIN A USER
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Email             string `json:"email"`
			Password          string `json:"password"`
			ReturnSecureToken bool   `json:"returnSecureToken"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Initialize the Identity Toolkit client
		creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

		idtClient, err := identitytoolkit.NewService(context.Background(), option.WithCredentialsFile(creds))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Identity Toolkit client"})
			return
		}

		// Authenticate the user
		authReq := &identitytoolkit.IdentitytoolkitRelyingpartyVerifyPasswordRequest{
			Email:             req.Email,
			Password:          req.Password,
			ReturnSecureToken: true,
		}
		authResp, err := idtClient.Relyingparty.VerifyPassword(authReq).Do()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Return the ID token
		c.JSON(http.StatusOK, gin.H{"token": authResp.IdToken})
	})

	// REGISTER A USER
	r.POST("/register", func(c *gin.Context) {
		var req struct {
			Email             string `json:"email"`
			Password          string `json:"password"`
			ReturnSecureToken bool   `json:"returnSecureToken"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Initialize the Identity Toolkit client
		saPath := filepath.Join(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
		opt := option.WithCredentialsFile(saPath)
		idtClient, err := identitytoolkit.NewService(context.Background(), opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Identity Toolkit client"})
			return
		}

		// Register the user
		registerReq := &identitytoolkit.IdentitytoolkitRelyingpartySignupNewUserRequest{
			Email:    req.Email,
			Password: req.Password,
		}

		registerResp, err := idtClient.Relyingparty.SignupNewUser(registerReq).Do()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to register user"})
			return
		}

		// Return the ID token
		c.JSON(http.StatusOK, gin.H{"token": registerResp.IdToken})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.DebugMsg("Listening on port " + srv.Addr + ":")
			log.Println("Listening on port " + srv.Addr + ":")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.DebugMsg("Shutdown server ...")
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.DebugMsg("Server Shutdown")
		log.Println("Server Shutdown")
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds")
		logger.DebugMsg("timeout of 5 seconds.")
	}
	logger.DebugMsg("Server exiting")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the ID token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID token is required"})
			c.Abort()
			return
		}

		// Initialize the Firebase app
		credentials, err := db.ReadCredentialsFromGCS(context.Background(), "proj-env", "keys.json")
		if err != nil {
			logger.DebugMsg("error reading credentials from the proj-env bucket")
		}

		opt := option.WithCredentialsJSON(credentials)
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firebase app"})
			c.Abort()
			return
		}

		// Get the Firebase Auth client
		authClient, err := app.Auth(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firebase Auth client"})
			c.Abort()
			return
		}

		// Verify the ID token
		_, err = authClient.VerifyIDToken(context.Background(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
			c.Abort()
			return
		}

		// Proceed to the next handler if the ID token is valid
		c.Next()
	}
}
