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
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.uber.org/zap"
)

// Must be called after TracerProvider initialization.
func Create() (*options.ClientOptions, string) {
	clientOptions := options.Client()
	clientOptions.Monitor = otelmongo.NewMonitor()
	clientOptions.ApplyURI(os.Getenv("MONGODB_SERVER_ADDR"))
	databaseName := os.Getenv("MONGODB_SERVER_DB")
	return clientOptions, databaseName
}

func Disconnect(client *mongo.Client, logger otelzap.LoggerWithCtx) {
	if err := client.Disconnect(logger.Context()); err != nil {
		logger.Error("Error during MongoDB disconnect", zap.Error(err))
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

func ExtractStringMap(value any) map[string]string {
	resMap := map[string]string{}
	switch casted := value.(type) {
	case bson.D:
		for _, elem := range casted {
			resMap[elem.Key], _ = elem.Value.(string)
		}
	case bson.M:
		for key, innerValue := range casted {
			resMap[key], _ = innerValue.(string)
		}
	}
	return resMap
}

func ConvertSlice[T any](docs []bson.M, converter func(bson.M) T) []T {
	resSlice := make([]T, 0, len(docs))
	for _, doc := range docs {
		resSlice = append(resSlice, converter(doc))
	}
	return resSlice
}
