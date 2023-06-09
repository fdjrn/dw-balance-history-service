package repository

import (
	"context"
	"github.com/fdjrn/dw-balance-history-service/internal/db"
	"github.com/fdjrn/dw-balance-history-service/internal/db/entity"
	"github.com/fdjrn/dw-balance-history-service/pkg/payload"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BalanceHistory struct {
}

var BalanceHistoryRepository = BalanceHistory{}

func (h *BalanceHistory) InsertBalanceHistory(history *entity.BalanceHistory) (int, interface{}, error) {

	result, err := db.Mongo.Collection.BalanceHistory.InsertOne(context.TODO(), history)
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusCreated, result.InsertedID, nil
}

func (h *BalanceHistory) FindByLastTransaction(r *payload.HistoryRequest) (int, interface{}, error) {

	filter := bson.D{{"uniqueId", r.UID}}
	opt := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(r.Limit)

	cursor, err := db.Mongo.Collection.BalanceHistory.Find(context.TODO(), filter, opt)
	if err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	var histories []entity.BalanceHistory
	if err = cursor.All(context.TODO(), &histories); err != nil {
		return fiber.StatusInternalServerError, nil, err
	}

	return fiber.StatusOK, histories, nil
}
