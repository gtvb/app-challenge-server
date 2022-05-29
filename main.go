package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gtvb/app-challenge-server/pkg/db"
	"github.com/gtvb/app-challenge-server/pkg/repository"
	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load env variables")
	}
}

func main() {
	pusherClient := pusher.Client{
		AppID:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: os.Getenv("PUSHER_CLUSTER"),
        Secure: true,
	}

	db, err := db.OpenDB(os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := repository.NewInstallerRequestsRepository(db)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	router.POST("/request_installer", func(c *gin.Context) {
		buf := new(strings.Builder)
		io.Copy(buf, c.Request.Body)

		var installerRequest repository.InstallerRequest
		err = json.Unmarshal([]byte(buf.String()), &installerRequest)
		if err != nil {
			c.String(500, err.Error())
		}

		bytes, _ := json.Marshal(installerRequest)

		err = repo.CreateRequest(installerRequest)
		if err != nil {
			c.String(500, err.Error())
		}

		channel := fmt.Sprintf("channel_%s", installerRequest.PlanId)
		event := fmt.Sprintf("event_%s", installerRequest.InstallerId)

		err = pusherClient.Trigger(channel, event, string(bytes))
		if err != nil {
			c.String(500, err.Error())
		}

		defer c.Request.Body.Close()
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	router.GET("/installer_requests/:planId/:installerId", func(c *gin.Context) {
		param1 := c.Param("planId")
		param2 := c.Param("installerId")

		planId, _ := strconv.Atoi(param1)
		installerId, _ := strconv.Atoi(param2)

		requests, err := repo.GetRequestsByPlanAndInstallerId(planId, installerId)
		if err != nil {
			c.String(500, err.Error())
		}

		data, _ := json.Marshal(requests)
		c.String(201, string(data))
	})

	router.Run()
}
