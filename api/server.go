package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
	"github.com/19chonm/461_1_23/db"
	"github.com/19chonm/461_1_23/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

func RunServer() {

	//create client to interact with google cloud storage
	const (
		projectid  = "group-24-ece461" 
		bucketname = "package-repo-ece461-24"
	)

	client, err := db.NewBucketClient(context.Background(), projectid, bucketname)
	if err != nil {
		fmt.Println("failed to create GCS client!")
	}
	defer client.Close()

	//Initialize go gin router
	r := gin.Default()

	//ROUTES

	//Get all packages
	r.GET("/packages", func(c *gin.Context) {
		packages, err := client.ListAllPackages()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(packages)
		c.JSON(http.StatusOK, packages)
	})

	// GetPackageInfo
	r.GET("/package/:name", func(c *gin.Context) {
		packageName := c.Param("name")
		packageInfo, err := client.GetPackageInfo(packageName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, packageInfo)
	})

	r.GET("/package", func(c *gin.Context) {
		packageName := c.Query("name")
		if packageName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name query parameter is missing"})
			return
		}
		searchResults, err := client.SearchPackage(packageName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, searchResults)
	})

	r.GET("/package/:name/download", func(c *gin.Context) {
		packageName := c.Param("name")
		f, err := os.Create(packageName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer f.Close()

		// download file from GCP and save to local file
		if err := client.DownloadPackage(packageName, f); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "Download successful")
	})

	r.DELETE("/package/:name", func(c *gin.Context) {
		packageName := c.Param("name")
		if err := client.RemovePackage(packageName); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "package deleted"})
	})

	r.POST("/upload", func(c *gin.Context) {
		fmt.Println("Hit upload a package")
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		if err := client.UploadPackage(header.Filename, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "package uploaded"})
	})

	srv := &http.Server{
		Addr:    ":8000",
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

func ListPackages(c *gin.Context) {
	const (
		projectID  = "your-project-id"
		bucketName = "your-bucket-name"
	)

	// create client to interact with Google Cloud Storage
	client, err := storage.NewClient(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Close()

	// get the bucket handle
	bucket := client.Bucket(bucketName)

	// list all objects in the bucket
	var packages []string
	query := &storage.Query{}
	objects := bucket.Objects(context.Background(), query)
	for {
		attrs, err := objects.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		packages = append(packages, attrs.Name)
	}

	c.JSON(http.StatusOK, gin.H{"packages": packages})
}
