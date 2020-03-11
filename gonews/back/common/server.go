package common

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	lock sync.Mutex
)

// StartServ start server at 8080
func StartServ(config *Config) {
	gin.SetMode(config.Mode)
	r := gin.Default()
	r.Use(GinLogger(LoggerFromConfig(config)), gin.Recovery())

	staticPath := config.StaticPath

	//r.LoadHTMLGlob("dist/*.html")        // 添加入口index.html
	//r.LoadHTMLFiles("dist/*/*")          // 添加资源路径
	r.Static("/static", path.Join(staticPath, "static"))               // 添加资源路径
	r.StaticFile("/", path.Join(staticPath, "index.html"))             // 前端接口
	r.StaticFile("/favicon.ico", path.Join(staticPath, "favicon.ico")) // 前端接口

	//配置跨域
	r.Use(cors.New(cors.Config{
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type", "ACCESS_TOKEN"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/news", getNewsByPage)
	r.GET("/api/pull", pullNewsFromGithub(config))
	r.GET("/api/job/status", getJobStatus)
	r.Run(fmt.Sprintf(":%d", config.Port))
}

func getNewsByPage(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageNum, _ := strconv.ParseInt(page, 10, 64)
	size := c.DefaultQuery("size", "10")
	pageSize, _ := strconv.ParseInt(size, 10, 64)
	news, count, err := GetPagedNews(pageNum, pageSize)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": count,
		"per":   pageSize,
		"items": news,
	})
}

func pullNewsFromGithub(config *Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		job := NewJob()
		_, err := job.CacheJob()
		go func() {
			lock.Lock()
			defer lock.Unlock()
			err := InitDataPuller(config)
			if err != nil {
				job.SetErr(err)
			} else {
				job.SetFinish()
			}
			job.CacheJob()
		}()
		if err != nil {
			job.SetErr(err)
			c.JSON(http.StatusOK, gin.H{
				"jobid":  job.ID,
				"status": job.Status,
				"err":    job.Err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"jobid":  job.ID,
				"status": job.Status,
				"err":    nil,
			})
		}
	}
}

func getJobStatus(c *gin.Context) {
	jobid := c.DefaultQuery("id", "1")
	job := &Job{ID: jobid}
	err := job.GetCacheJob()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"jobid":  job.ID,
		"status": job.Status,
		"err":    job.Err.Error(),
	})
}
