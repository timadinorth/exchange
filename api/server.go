package api

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	redisStore "github.com/gofiber/storage/redis/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/timadinorth/bet-exchange/docs"
	"github.com/timadinorth/bet-exchange/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBHost         string `mapstructure:"DATABASE_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	CacheUrl       string `mapstructure:"CACHE_URL"`
	CachePassword  string `mapstructure:"CACHE_PASSWORD"`
	CacheDB        string `mapstructure:"CACHE_DB"`
	SessionDB      string `mapstructure:"CACHE_SESSION_DB"`
	HttpsEndpoint  string `mapstructure:"HTTPS_ENDPOINT"`
	HttpsCrt       string `mapstructure:"HTTPS_CRT"`
	HttpsKey       string `mapstructure:"HTTPS_KEY"`
}

type Server struct {
	Log       *logrus.Logger
	DB        *gorm.DB
	Web       *fiber.App
	Config    *Config
	Cache     *redis.Client
	Session   *session.Store
	validator *validator.Validate
}

func (s *Server) InitLogger() {
	s.Log = logrus.New()
}

func (s *Server) InitWeb() {
	s.validator = validator.New()
	db, err := strconv.Atoi(s.Config.SessionDB)
	if err != nil {
		s.Log.Fatal("Redis config wrong database")
	}

	s.Web = fiber.New()
	s.Session = session.New(session.Config{
		Storage: redisStore.New(redisStore.Config{
			Addrs:    []string{s.Config.CacheUrl},
			Password: s.Config.CachePassword,
			Database: db,
		}),
	})
}

func (s *Server) LoadConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		s.Log.Fatal("Failed to read config")
	}

	s.Config = new(Config)
	err = viper.Unmarshal(s.Config)
	if err != nil {
		s.Log.Fatal("Failed to parse config")
	}
}

func (s *Server) ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", s.Config.DBHost, s.Config.DBUserName, s.Config.DBUserPassword, s.Config.DBName)

	s.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		s.Log.Fatal("Failed to connect to the Database")
	}
	s.Log.Info("Connected Successfully to the Database")
}

func (s *Server) SetupModels() error {
	return s.DB.AutoMigrate(&model.Category{}, &model.Competition{}, &model.User{})
}

func (s *Server) CleanupModels() error {
	return s.DB.Migrator().DropTable(&model.Category{}, &model.Competition{}, &model.User{})
}

func (s *Server) ConnectCache() {
	db, err := strconv.Atoi(s.Config.CacheDB)
	if err != nil {
		s.Log.Fatal("Redis config wrong database")
	}
	s.Cache = redis.NewClient(&redis.Options{
		Addr:     s.Config.CacheUrl,
		Password: s.Config.CachePassword,
		DB:       db,
	})
	status := s.Cache.Ping()
	s.Log.Info("Connected to Cache ", status)
}

func (s *Server) Start() {
	s.Log.Fatal(s.Web.Listen(":8080"))
}
