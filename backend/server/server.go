package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go/valkeycompat"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/db"
	"github.com/shafikshaon/shortlink/logger"
)

type Server struct {
	Cfg       *config.Config
	Gin       *gin.Engine
	MasterDB  *gorm.DB
	ReplicaDB *gorm.DB
	Cache     valkeycompat.Cmdable
}

func NewServer(ctx context.Context, cfg *config.Config) *Server {
	logger.Info("Initializing server")

	srv := &Server{Cfg: cfg}

	masterDB, err := db.NewPostgresDB(ctx, cfg.DB, false)
	if err != nil {
		logger.Error("Failed to connect to master database", zap.Error(err))
	} else {
		srv.MasterDB = masterDB
		logger.Info("Connected to master database",
			zap.String("host", cfg.DB.Host),
			zap.String("database", cfg.DB.Name))
	}

	if srv.MasterDB != nil {
		if err := db.RunAutoMigration(srv.MasterDB); err != nil {
			logger.Error("Auto migration failed", zap.Error(err))
		}
	}

	replicaDB, err := db.NewPostgresDB(ctx, cfg.DB, true)
	if err != nil {
		logger.Warn("Failed to connect to replica database, falling back to master for reads", zap.Error(err))
		srv.ReplicaDB = srv.MasterDB
	} else {
		srv.ReplicaDB = replicaDB
		logger.Info("Connected to replica database")
	}

	cache, err := db.NewValkeyClient(ctx, cfg.Cache)
	if err != nil {
		logger.Error("Failed to connect to Valkey", zap.Error(err))
	} else {
		srv.Cache = cache
		logger.Info("Connected to Valkey",
			zap.String("host", cfg.Cache.Host),
			zap.String("port", cfg.Cache.Port))
	}

	logger.Info("Server initialization completed")
	return srv
}

func (s *Server) GetDSN(isReplica bool) string {
	if isReplica {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			s.Cfg.DB.ReadHost, s.Cfg.DB.ReadPort, s.Cfg.DB.ReadUser, s.Cfg.DB.ReadPassword, s.Cfg.DB.Name)
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		s.Cfg.DB.Host, s.Cfg.DB.Port, s.Cfg.DB.User, s.Cfg.DB.Password, s.Cfg.DB.Name)
}
