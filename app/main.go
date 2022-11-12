package main

import (
	"auth/ent"
	"context"
	"fmt"

	"auth/domain"
	handler "auth/pkg/controller/http"
	redisCli "auth/pkg/repository/redis"
	"auth/pkg/service"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("get configuration error")
	}
}

func main() {

	client := connectRedis()

	token := domain.JwtToken{
		AccessSecret: viper.GetString(`token.secret`),
		AccessTtl:    viper.GetDuration(`token.ttl`) * time.Minute,
	}
	redis := redisCli.NewRedisRepo(client)
	timeout := viper.GetDuration(`timeout`) * time.Second

	db := connectDB()
	defer db.Close()
	jwtUsecase := service.NewJWTUseCase(token, redis)
	userUsecase := service.NewUserService(db, jwtUsecase, timeout)

	e := echo.New()
	handler.NewUserHandler(e, userUsecase, jwtUsecase)

	err := e.Start(viper.GetString(`addr`))
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg(`shutting down the server`)
	}
}

func connectRedis() *redis.Client {

	address := viper.GetString(`redis.address`)
	password := viper.GetString(`redis.password`)

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal().Err(err).Msg("redis ping error")
	}
	return client
}

func connectDB() *ent.Client {

	username := viper.GetString(`postgres.user`)
	password := viper.GetString(`postgres.password`)
	hostname := viper.GetString(`postgres.host`)
	port := viper.GetInt(`postgres.port`)
	dbname := viper.GetString(`postgres.dbname`)

	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", hostname, port, username, dbname, password))
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening connection to postgres")
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed creating schema resources")
	}

	ctx := context.Background()
	user, err := client.User.Create().
		SetEmail("gajayip@gmail.com").
		SetPassword("pass").
		SetUsername("nazerke").Save(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("create user error")
	}
	_, err = client.Task.Create().
		SetTitle("First task").
		SetUser(user).
		Save(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("create task error")
	}
	return client
}
