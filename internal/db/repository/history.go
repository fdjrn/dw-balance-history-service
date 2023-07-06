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
	Pagination  *entity.PaginatedRequest
}

func NewBalanceHistoryRepository() BalanceHistoryRepository {
	return BalanceHistoryRepository{
		Entity:      new(entity.BalanceHistory),
		Transaction: new(entity.BalanceTransaction),
		Pagination:  new(entity.PaginatedRequest),
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

//func (h *BalanceHistoryRepository) InsertBalanceHistory(history *entity.BalanceHistory) (int, interface{}, error) {
//
//	result, err := db.Mongo.Collection.BalanceHistory.InsertOne(context.TODO(), history)
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, err
//	}
//
//	return fiber.StatusCreated, result.InsertedID, nil
//}

func (h *BalanceHistoryRepository) Create() (interface{}, error) {

	result, err := db.Mongo.Collection.BalanceHistory.InsertOne(context.TODO(), h.Entity)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

/*
FindByLastTransaction
out-params:

	code: int,
	result: interface{},
	length: int,
	err: error
*/
//func (h *BalanceHistoryRepository) FindByLastTransaction(r *handlers.HistoryRequest) (int, interface{}, int, error) {
//
//	filter := bson.D{{"uniqueId", r.UID}}
//	opt := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(r.Limit)
//
//	cursor, err := db.Mongo.Collection.BalanceHistory.Find(context.TODO(), filter, opt)
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, 0, err
//	}
//
//	var histories []entity.BalanceHistory
//	if err = cursor.All(context.TODO(), &histories); err != nil {
//		return fiber.StatusInternalServerError, nil, 0, err
//	}
//
//	return fiber.StatusOK, histories, len(histories), nil
//}

/*
FindByPeriod
out-params:

	code: int,
	result: interface{},
	length: int,
	err: error
*/
//func (h *BalanceHistoryRepository) FindByPeriod(r *handlers.HistoryRequest) (int, interface{}, int, error) {
//
//	// 1. add fields
//	addFieldStage := bson.D{
//		{"$addFields", bson.D{
//			{"year", bson.D{{"$year", bson.D{{"$toDate", "$transDate"}}}}},
//			{"month", bson.D{{"$month", bson.D{{"$toDate", "$transDate"}}}}},
//		}},
//	}
//	// 2. match Stages
//	matchStage := bson.D{
//		{"$match", bson.D{
//			{"year", r.Period.Year},
//			{"month", r.Period.Month},
//		}},
//	}
//
//	cursor, err := db.Mongo.Collection.BalanceHistory.Aggregate(
//		context.TODO(),
//		mongo.Pipeline{
//			bson.D{{"$match", bson.D{{"uniqueId", r.UID}}}},
//			bson.D{{"$sort", bson.D{{"_id", -1}}}},
//			addFieldStage,
//			matchStage,
//		})
//
//	if err != nil {
//		return fiber.StatusInternalServerError, nil, 0, err
//	}
//
//	var histories []entity.BalanceHistory
//	for cursor.Next(context.TODO()) {
//		var item entity.BalanceHistory
//		err2 := cursor.Decode(&item)
//		if err2 != nil {
//			log.Println("error found on decoding cursor: ", err2.Error())
//			continue
//		}
//
//		histories = append(histories, item)
//	}
//
//	return fiber.StatusOK, histories, len(histories), nil
//}

/*
FindAllPaginated
function args:

	*payload.PaginatedRequest

return:

	code int,
	data interfaces{},
	totalDocument int64,
	totalPages int,
	err error
*/
func (h *BalanceHistoryRepository) FindAllPaginated() (interface{}, int64, int64, error) {
	filter := bson.D{}

	filter = append(filter, bson.D{
		{"partnerId", h.Pagination.PartnerID},
		{"merchantId", h.Pagination.MerchantID},
		{"terminalId", h.Pagination.TerminalID},
	}...)

	skipValue := (h.Pagination.Page - 1) * h.Pagination.Size

	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	cursor, err := db.Mongo.Collection.BalanceHistory.Find(
		ctx,
		filter,
		options.Find().
			SetSort(bson.D{{"_id", -1}}).
			SetSkip(skipValue).
			SetLimit(h.Pagination.Size),
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

	totalPages := math.Ceil(float64(totalDocs) / float64(h.Pagination.Size))

	return &histories, totalDocs, int64(totalPages), nil
}
