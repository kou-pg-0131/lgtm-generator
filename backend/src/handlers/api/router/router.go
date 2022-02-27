package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/koki-develop/lgtm-generator/backend/src/adapters/controllers"
	"github.com/koki-develop/lgtm-generator/backend/src/adapters/controllers/middlewares"
	imgsrepo "github.com/koki-develop/lgtm-generator/backend/src/adapters/gateways/images"
	lgtmsrepo "github.com/koki-develop/lgtm-generator/backend/src/adapters/gateways/lgtms"
	"github.com/koki-develop/lgtm-generator/backend/src/adapters/gateways/notifier"
	rptsrepo "github.com/koki-develop/lgtm-generator/backend/src/adapters/gateways/reports"
	"github.com/koki-develop/lgtm-generator/backend/src/entities"
	"github.com/koki-develop/lgtm-generator/backend/src/infrastructures"
	infiface "github.com/koki-develop/lgtm-generator/backend/src/infrastructures/iface"
	"github.com/koki-develop/lgtm-generator/backend/src/usecases"
)

func withContext(h func(ctx infiface.Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h(infrastructures.NewContextFromGin(ctx))
	}
}

func New() *gin.Engine {
	r := gin.Default()

	s := infrastructures.NewSlackClient(&infrastructures.SlackClientConfig{
		AccessToken: os.Getenv("NOTIFICATION_SLACK_ACCESS_TOKEN"),
		HTTPAPI:     http.DefaultClient,
	})
	rdr := infrastructures.NewRenderer(&infrastructures.RendererConfig{SlackAPI: s})
	imgse := infrastructures.NewGoogleImageSearchEngine(&infrastructures.GoogleImageSearchEngineConfig{
		APIKey:         os.Getenv("GOOGLE_API_KEY"),
		SearchEngineID: os.Getenv("GOOGLE_CUSTOM_SEARCH_ENGINE_ID"),
	})
	db := infrastructures.NewGureguDynamoDB()
	s3lgtms := infrastructures.NewS3(&infrastructures.S3Config{
		Bucket: fmt.Sprintf("lgtm-generator-backend-%s-images", os.Getenv("STAGE")),
	})
	lgtmgen := infrastructures.NewLGTMGenerator()

	n := notifier.New(&notifier.Config{
		Slack:       s,
		Channel:     fmt.Sprintf("lgtm-generator-backend-%s-reports", os.Getenv("STAGE")),
		FileStorage: s3lgtms,
	})
	lgtmsrepo := lgtmsrepo.NewRepository(&lgtmsrepo.RepositoryConfig{
		LGTMGenerator: lgtmgen,
		DynamoDB:      db,
		DBPrefix:      fmt.Sprintf("lgtm-generator-backend-%s", os.Getenv("STAGE")),
		FileStorage:   s3lgtms,
		HTTPAPI:       http.DefaultClient,
	})
	imgsrepo := imgsrepo.NewRepository(&imgsrepo.RepositoryConfig{
		ImageSearchEngine: imgse,
	})
	rptsrepo := rptsrepo.NewRepository(&rptsrepo.RepositoryConfig{
		DynamoDB: db,
		DBPrefix: fmt.Sprintf("lgtm-generator-backend-%s", os.Getenv("STAGE")),
	})

	lgtmsuc := usecases.NewLGTMsUsecase(&usecases.LGTMsUsecaseConfig{
		LGTMsRepository: lgtmsrepo,
	})
	imgsuc := usecases.NewImagesUsecase(&usecases.ImagesUsecaseConfig{
		ImagesRepository: imgsrepo,
	})
	rptsuc := usecases.NewReportsUsecase(&usecases.ReportsUsecaseConfig{
		ReportsRepository: rptsrepo,
		LGTMsRepository:   lgtmsrepo,
		Notifier:          n,
	})

	{
		cfg := &middlewares.CORSMiddlewareConfig{
			Renderer:     rdr,
			AllowOrigins: []string{},
			AllowMethods: []string{"GET", "POST", "OPTIONS"},
			AllowHeaders: []string{
				"Origin",
				"Content-Length",
				"Content-Type",
				"Accept-Encoding",
				"Sentry-Trace",
			},
		}
		stg := os.Getenv("STAGE")
		switch stg {
		case "local":
			cfg.AllowOrigins = []string{"http://localhost:3000"}
		case "dev":
			cfg.AllowOrigins = []string{"https://lgtm-generator-*-koki-develop.vercel.app"}
		case "prod":
			cfg.AllowOrigins = []string{"https://lgtmgen.org"}
		default:
			panic(fmt.Sprintf("unknown stage: %s", stg))
		}
		cors := middlewares.NewCORSMiddleware(cfg)
		r.Use(withContext(cors.Apply))
	}

	v1 := r.Group("/v1")
	{
		ctrl := controllers.NewHealthController(&controllers.HealthControllerConfig{
			Renderer: rdr,
		})
		v1.GET("/h", withContext(ctrl.Standard))
	}
	{
		ctrl := controllers.NewImagesController(&controllers.ImagesControllerConfig{
			Renderer:      rdr,
			ImagesUsecase: imgsuc,
		})
		v1.GET("/images", withContext(ctrl.Search))
	}
	{
		ctrl := controllers.NewLGTMsController(&controllers.LGTMsControllerConfig{
			Renderer:     rdr,
			LGTMsUsecase: lgtmsuc,
		})
		v1.GET("/lgtms", withContext(ctrl.Index))
		v1.POST("/lgtms", withContext(ctrl.Create))
	}
	{
		ctrl := controllers.NewReportsController(&controllers.ReportsControllerConfig{
			Renderer:       rdr,
			ReportsUsecase: rptsuc,
		})
		v1.POST("/reports", withContext(ctrl.Create))
	}
	r.NoRoute(withContext(func(ctx infiface.Context) {
		rdr.NotFound(ctx, entities.ErrCodeNotFound, errors.New("no route"))
	}))

	return r
}
