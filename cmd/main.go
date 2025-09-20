package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mksmstpck/spy_cat_agency/internal/config"
	"github.com/mksmstpck/spy_cat_agency/internal/db"
	"github.com/mksmstpck/spy_cat_agency/internal/events"
	"github.com/mksmstpck/spy_cat_agency/internal/handlers"
	"github.com/mksmstpck/spy_cat_agency/internal/services"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.Function), f.Line)
		},
	}
	logrus.SetFormatter(formatter)
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}

func main() {
	ctx := context.TODO()

	config := config.NewConfig()

	pgconn, err := pgxpool.New(ctx, config.PostgregUrl)
	if err != nil {
		logrus.Error(err)
	}

	db := db.NewDB(pgconn)

	services := services.NewServices(*db)

	if err := events.NewEvents(*services, config).LoadBreeds(ctx); err != nil {
		logrus.Error(err)
	}

	handlers := handlers.NewHandlers(config, services)
	handlers.HandleAll(ctx)
}
