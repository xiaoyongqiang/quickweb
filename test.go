package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mog "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://admin:123456@127.0.0.1:27017/admin")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检测连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// 获取collection
	detection_coll := client.Database("wisteria_detect").Collection("ids_detection")

	// 保存单条数据
	saveOne(detection_coll)

	// 批量保存数据
	batchSave(detection_coll)

	// 查询单条数据
	findOne(detection_coll)

	// 批量查询数据
	findMany(detection_coll)

	// 更新多条数据
	updateMany(detection_coll)

	// 删除单条数据
	deleteOne(detection_coll)
}

type DetectionDO struct {
	Id          string `bson:"_id"`
	ComId       string `bson:"comId"`
	DetectionId string `bson:"detectionId"`
}

/**
 * 插入单个
 */
func saveOne(detection_coll *mongo.Collection) {
	detection := DetectionDO{
		Id:          "100001",
		ComId:       "7070714a613939797533",
		DetectionId: "abc123456",
	}
	objId, err := detection_coll.InsertOne(context.TODO(), detection)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("_id:", objId.InsertedID)
}

/**
 * 批量插入
 */
func batchSave(detection_coll *mongo.Collection) {
	detection1 := DetectionDO{
		Id:          "100002",
		ComId:       "7070714a613939797533",
		DetectionId: "abc123456-2",
	}
	detection2 := DetectionDO{
		Id:          "100003",
		ComId:       "7070714a613939797533",
		DetectionId: "abc123456-3",
	}

	models := []mog.WriteModel{
		mog.NewInsertOneModel().SetDocument(detection1),
		mog.NewInsertOneModel().SetDocument(detection2),
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := detection_coll.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("insertCount:", res.InsertedCount)
}

/**
 * 查询单个
 */
func findOne(detection_coll *mongo.Collection) {
	filter := bson.D{{Key: "comId", Value: "7070714a613939797533"}, {Key: "detectionId", Value: "684b8788169c481e9dff73889b170072zh1pxms0"}}

	var result map[string]interface{}
	err := detection_coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(result)
}

/**
 * 批量查询
 */
func findMany(detection_coll *mongo.Collection) {
	filter := bson.D{{Key: "comId", Value: "7070714a613939797533"}}

	cur, err := detection_coll.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	var list []*DetectionDO
	err = cur.All(context.Background(), &list)
	if err != nil {
		log.Fatal(err)
	}

	_ = cur.Close(context.Background())

	for _, one := range list {
		fmt.Printf("_id:%s, comId:%s, detectionId:%s\n", one.Id, one.ComId, one.DetectionId)
	}
}

/**
 * 更新
 */
func updateMany(detection_coll *mongo.Collection) {
	detectionIds := []string{"d42161747d9547b3aaee8f6464792c40vcnzd8q6", "684b8788169c481e9dff73889b170072zh1pxms0", "bed744bcae2d47d5a682369d76acda52jgpphi5y"}
	filter := bson.D{{Key: "comId", Value: "7070714a613939797533"}, {Key: "detectionId", Value: bson.D{{Key: "$in", Value: detectionIds}}}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "detectionName", Value: "使用golang测试批量更新咯"},
		}},
	}
	result, err := detection_coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("modifiedCount:%d\n", result.ModifiedCount)
}

/**
 * 删除
 */
func deleteOne(detection_coll *mongo.Collection) {
	filter := bson.D{{Key: "comId", Value: "7070714a613939797533"}, {Key: "detectionId", Value: "684b8788169c481e9dff73889b170072zh1pxms0"}}
	result, err := detection_coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleteCount:%d\n", result.DeletedCount)
}
