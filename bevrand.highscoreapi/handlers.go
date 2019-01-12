package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	openlog "github.com/opentracing/opentracing-go/log"
	"gopkg.in/oauth2.v3/utils/uuid"
	"log"
	"net/http"
	"strconv"
)

func PingPong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ShowHighScores(user string, playlist string, c *gin.Context, ctx context.Context) ([]Score, int) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ShowHighScores")
	span.SetTag("Method", "ShowHighScores")
	defer span.Finish()
	key := user + ":" + playlist

	res, err := db.Cmd("EXISTS", key).Int()
	if err != nil {
		log.Fatal(err)
	}
	exists := res != 0
	var redisResult []Score

	if !exists {
		span.LogFields(
			openlog.String("http_status_code", "404"),
			openlog.String("body", "not found"),
		)
		span.Finish()
		b, _ := uuid.NewRandom()
		localUuid := fmt.Sprintf("%x-%x-%x-%x-%x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		notFoundError := ErrorModel{
			Message: "Could not find combination of user: " + user + " & list: " + playlist,
			UniqueCode: localUuid}
		respondWithJson(c, http.StatusNotFound, notFoundError, ctx)
		return nil, 0
	}

	result, err := db.Cmd("HGETALL", key).Map()
	if err != nil {
		log.Fatal(err)
	}
	for drink, score := range result{
		i, err := strconv.Atoi(score)
		if err != nil {
			log.Fatal("cannot convert string to int")
		}
		newScore := Score{
			Drink: drink,
			Rolled: i}
		redisResult = append(redisResult, newScore)
	}
	body, _ := json.Marshal(redisResult)

	span.LogFields(
		openlog.String("http_status_code", "200"),
		openlog.String("body", string(body)),
	)
	return redisResult, http.StatusOK
}

func CreateNewHighScore(user string, playlist string, drink string, ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CreateNewHighScore")
	span.SetTag("Method", "IncreaseHighScore")

	defer span.Finish()
	key := user + ":" + playlist

	res, err := db.Cmd("EXISTS", key).Int()
	if err != nil {
		log.Fatal(err)
	}
	exists := res != 0

	if !exists {
		span.LogFields(
			openlog.String("http_status_code", "200"),
			openlog.String("body", "New entry, setting new key: " + key),
		)
		err = db.Cmd("HMSET", key, drink, 1).Err
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = db.Cmd("HINCRBY", key, drink, 1).Err
	if err != nil {
		log.Fatal(err)
	}
	span.LogFields(
		openlog.String("http_status_code", "200"),
		openlog.String("body", "Increased " + drink + " by 1"),
	)

	IncreaseGlobalCount(drink, ctx)
	return
}

func IncreaseGlobalCount(drink string, ctx context.Context){
	span, _ := opentracing.StartSpanFromContext(ctx, "IncreaseGlobalCount")
	span.SetTag("Method", "GlobalHighscore")

	defer span.Finish()
	key := GLOBALNAME + ":" + GLOBALLIST

	res, err := db.Cmd("EXISTS", key).Int()
	if err != nil {
		log.Fatal(err)
	}
	exists := res != 0

	if !exists {
		span.LogFields(
			openlog.String("http_status_code", "200"),
			openlog.String("body", "New entry, setting new key: " + key),
		)
		err = db.Cmd("HMSET", key, drink, 1).Err
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = db.Cmd("HINCRBY", key, drink, 1).Err
	if err != nil {
		log.Fatal(err)
	}
	span.LogFields(
		openlog.String("http_status_code", "200"),
		openlog.String("body", "Increased " + drink + " by 1"),
	)
	return
}
