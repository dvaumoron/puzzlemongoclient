/*
 *
 * Copyright 2023 puzzlemongoclient authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package puzzlemongoclient

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create() (*options.ClientOptions, string) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_SERVER_ADDR"))
	databaseName := os.Getenv("MONGODB_SERVER_DB")
	return clientOptions, databaseName
}

func Disconnect(client *mongo.Client, ctx context.Context) {
	if err := client.Disconnect(ctx); err != nil {
		log.Print("Error during MongoDB disconnect :", err)
	}
}

func ExtractCreateDate(doc bson.M) time.Time {
	id, _ := doc["_id"].(primitive.ObjectID)
	return id.Timestamp()
}

func ExtractUint64(value any) uint64 {
	switch casted := value.(type) {
	case int32:
		return uint64(casted)
	case int64:
		return uint64(casted)
	}
	return 0
}

func ExtractBinary(value any) []byte {
	binary, _ := value.(primitive.Binary)
	return binary.Data
}
