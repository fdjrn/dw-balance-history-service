package repository

import (
	"context"
	"errors"
	"github.com/fdjrn/dw-balance-history-service/internal/db"
	"github.com/fdjrn/dw-balance-history-service/internal/db/entity"
	"github.com/fdjrn/dw-balance-history-service/internal/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

type BalanceHistoryRepository struct {
	Entity      *entity.BalanceHistory
	Transaction *entity.BalanceTransaction
	Request     *entity.PaginatedRequest
}

func NewBalanceHistoryRepository() BalanceHistoryRepository {
	return BalanceHistoryRepository{
		Entity:      new(entity.BalanceHistory),
		Transaction: new(entity.BalanceTransaction),
		Request:     new(entity.PaginatedRequest),
	}
}

func (h *BalanceHistoryRepository) IsExists(receiptNo string) bool {

	hist := new(entity.BalanceHistory)
	err := db.Mongo.Collection.BalanceHistory.FindOne(
		context.TODO(), bson.D{{"receiptNumber", receiptNo}},
	).Decode(hist)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}

		utilities.Log.Println(err.Error())
	}

	return true
}

func (h *BalanceHistoryRepository) FindByID(objectId interface{}) (*entity.BalanceHistory, error) {
	result := new(entity.BalanceHistory)
	err := db.Mongo.Collection.BalanceHistory.FindOne(context.TODO(), bson.D{
		{"_id", objectId},
	}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *BalanceHistoryRepository) Create() (interface{}, error) {

	result, err := db.Mongo.Collection.BalanceHistory.InsertOne(context.TODO(), h.Entity)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (h *BalanceHistoryRepository) FindByLastTransaction() (interface{}, int, error) {

	filter := bson.D{
		{"partnerId", h.Request.PartnerID},
		{"merchantId", h.Request.MerchantID},
		{"terminalId", h.Request.TerminalID},
	}

	opt := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(h.Request.Limit)

	cursor, err := db.Mongo.Collection.BalanceHistory.Find(context.TODO(), filter, opt)
	if err != nil {
		return nil, 0, err
	}

	var histories []entity.BalanceHistory
	if err = cursor.All(context.TODO(), &histories); err != nil {
		return nil, 0, err
	}

	return histories, len(histories), nil
}

func (h *BalanceHistoryRepository) FindAllPaginated() (interface{}, int64, int64, error) {
	filter := bson.D{}

	filter = append(filter, bson.D{
		{"partnerId", h.Request.PartnerID},
		{"merchantId", h.Request.MerchantID},
		{"terminalId", h.Request.TerminalID},
	}...)

	if h.Request.Periods != nil {
		filter = append(filter, bson.D{
			{"transDateNumeric", bson.D{
				{"$gte", h.Request.Periods.StartDate.UnixMilli()},
				{"$lte", h.Request.Periods.EndDate.UnixMilli()},
			}},
		}...)
	}

	skipValue := (h.Request.Page - 1) * h.Request.Size

	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	cursor, err := db.Mongo.Collection.BalanceHistory.Find(
		ctx,
		filter,
		options.Find().
			SetSort(bson.D{{"_id", -1}}).
			SetSkip(skipValue).
			SetLimit(h.Request.Size),
	)

	if err != nil {
		return nil, 0, 0, err
	}

	totalDocs, _ := db.Mongo.Collection.BalanceHistory.CountDocuments(ctx, filter)
	var histories []entity.BalanceHistory
	if err = cursor.All(context.TODO(), &histories); err != nil {
		return nil, 0, 0, err
	}

	if len(histories) == 0 {
		return nil, 0, 0, errors.New("empty results or last pages has been reached")
	}

	totalPages := math.Ceil(float64(totalDocs) / float64(h.Request.Size))

	return &histories, totalDocs, int64(totalPages), nil
}
