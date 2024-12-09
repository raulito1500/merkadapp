package repository

import (
	"context"
	"errors"
	"time"

	"github.com/raulito1500/merkadapp/src/bill/entities"
	"github.com/raulito1500/merkadapp/src/bill/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BillMongoRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

const BILL_COLLECTION = "bills"

func NewBillMongoRepository(db *mongo.Database) BillRepository {
	return &BillMongoRepository{
		db:   db,
		coll: db.Collection(BILL_COLLECTION),
	}
}

func (m *BillMongoRepository) ListBills() []*models.Bill {
	opts := options.Find().SetSort(bson.D{{"date", -1}, {"_id", -1}})
	opts = opts.SetLimit(1000)

	cursor, err := m.coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return []*models.Bill{}
	}
	defer cursor.Close(context.TODO())

	results := []*models.Bill{}

	if err := cursor.All(context.TODO(), &results); err != nil {
		return []*models.Bill{}
	}
	return results
}
func (m *BillMongoRepository) ListBill(id string) (models.Bill, error) {
	var result models.Bill
	objId, _ := primitive.ObjectIDFromHex(id)
	err := m.coll.FindOne(context.TODO(), bson.D{{"_id", objId}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, errors.New("Not found")
		}
	}
	return result, nil
}

func (m *BillMongoRepository) InsertBill(bill *models.Bill) (string, error) {
	result, err := m.coll.InsertOne(context.TODO(), bill)
	if err != nil {
		return "", err
	}
	mongoId := result.InsertedID
	return mongoId.(primitive.ObjectID).Hex(), nil
}

func (m *BillMongoRepository) UdpateBill(id string, bill *models.Bill) error {
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objId}}

	_, err := m.coll.ReplaceOne(context.TODO(), filter, bill)

	if err != nil {
		return err
	}
	return nil
}

func (m *BillMongoRepository) MergeBills(idDestination string, idsOrigen []string) error {
	destId, _ := primitive.ObjectIDFromHex(idDestination)

	var destDoc models.Bill
	err := m.coll.FindOne(context.TODO(), bson.M{"_id": destId}).Decode(&destDoc)
	if err != nil {
		return err
	}

	var originIDs []primitive.ObjectID
	for _, id := range idsOrigen {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		originIDs = append(originIDs, objID)
	}

	var sourceDocs []models.Bill
	for _, id := range originIDs {
		var doc models.Bill
		err = m.coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&doc)
		if err != nil {
			return err
		}
		sourceDocs = append(sourceDocs, doc)
	}

	var combinedItems []models.BillItem
	combinedTotal := destDoc.Total

	if destDoc.Items != nil {
		for _, item := range destDoc.Items {
			if item != nil {
				combinedItems = append(combinedItems, *item)
			}
		}
	}

	for _, doc := range sourceDocs {
		if doc.Items != nil {
			for _, item := range doc.Items {
				if item != nil {
					combinedTotal += item.Total
					combinedItems = append(combinedItems, *item)
				}
			}
		}
	}

	_, err = m.coll.UpdateOne(
		context.TODO(),
		bson.M{"_id": destId},
		bson.M{
			"$set": bson.M{
				"items": combinedItems,
				"total": combinedTotal,
			},
		},
	)
	if err != nil {
		return err
	}

	_, err = m.coll.DeleteMany(
		context.TODO(),
		bson.M{"_id": bson.M{"$in": originIDs}},
	)
	if err != nil {
		return nil
	}

	return nil
}

func (m *BillMongoRepository) MarkSpentItem(idBill string, idItem string) error {
	objId, _ := primitive.ObjectIDFromHex(idBill)
	filter := bson.D{
		{"_id", objId},
		{"items", bson.M{"$elemMatch": bson.M{"_id": idItem}}},
	}

	update := bson.D{
		{"$set", bson.M{"items.$.spent_date": time.Now()}},
	}
	_, err := m.coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}
	return nil
}

func (m *BillMongoRepository) TotalByMonth(startDate time.Time, endDate time.Time) ([]*entities.BillTotal, error) {
	pipeline := bson.A{
		bson.D{
			{"$match", bson.D{
				{"date", bson.D{
					{"$gte", startDate},
					{"$lt", endDate},
				}},
			}},
		},
		bson.D{
			{"$addFields", bson.D{
				{"month", bson.D{{"$month", "$date"}}},
			}},
		},
		bson.D{
			{"$group", bson.D{
				{"_id", "$month"},
				{"total", bson.D{{"$sum", "$total"}}},
			}},
		},
		bson.D{{"$sort", bson.D{{"_id", -1}}}},
	}

	cursor, err := m.coll.Aggregate(context.TODO(), pipeline)
	var result []*entities.BillTotal
	if err = cursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}
	if len(result) > 0 {
		return result, nil
	} else {
		return []*entities.BillTotal{}, errors.New("Not found")
	}
}
