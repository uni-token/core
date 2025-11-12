package store

import (
	"encoding/json"
	"log"
	"time"

	"go.etcd.io/bbolt"
)

type Bucket[T any] struct {
	bucketName string
}

func CreateBucket[T any](bucketName string) (Bucket[T], error) {
	b := Bucket[T]{
		bucketName: bucketName,
	}
	err := Db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(b.bucketName))
		return err
	})
	return b, err
}

func DeleteBucket(bucketName string) error {
	return Db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
}

func InitBucket[T any](bucketName string) Bucket[T] {
	b, err := CreateBucket[T](bucketName)
	if err != nil {
		log.Fatal("Failed to create bucket:", err)
	}
	return b
}

func (b *Bucket[T]) List() ([]T, error) {
	result := make([]T, 0)
	return result, Db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		return b.ForEach(func(k, v []byte) error {
			var data T
			err := json.Unmarshal(v, &data)
			if err != nil {
				return err
			}
			result = append(result, data)
			return nil
		})
	})
}

func (b *Bucket[T]) Put(key string, data T) error {
	v, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return Db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		return b.Put([]byte(key), v)
	})
}

func (b *Bucket[T]) Get(key string) (T, error) {
	var data T
	err := Db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		v := b.Get([]byte(key))
		return json.Unmarshal(v, &data)
	})
	return data, err
}

func (b *Bucket[T]) Delete(key string) error {
	return Db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		return b.Delete([]byte(key))
	})
}

func (b *Bucket[T]) Clear() error {
	return Db.Update(func(tx *bbolt.Tx) error {
		err := tx.DeleteBucket([]byte(b.bucketName))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(b.bucketName))
		return err
	})
}

func (b *Bucket[T]) Count() (int, error) {
	var count int
	err := Db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		count = b.Inspect().KeyN
		return nil
	})
	return count, err
}

var (
	Db *bbolt.DB

	Users      Bucket[UserInfo]
	Apps       Bucket[AppInfo]
	LLMKeys    Bucket[LLMKey]
	AppPresets Bucket[AppPreset]
	Usage      Bucket[TokenUsage]
	Providers  Bucket[[]byte]
)

func Init(dbPath string) {
	var err error
	Db, err = bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	Users = InitBucket[UserInfo]("users")
	Usage = InitBucket[TokenUsage]("usage")
	Apps = InitBucket[AppInfo]("apps")
	LLMKeys = InitBucket[LLMKey]("llm_keys")
	AppPresets = InitBucket[AppPreset]("app_presets")
	Providers = InitBucket[[]byte]("providers")

	count, err := AppPresets.Count()
	if err != nil {
		panic("Failed to count app presets: " + err.Error())
	}
	if count == 0 {
		defaultPreset := AppPreset{
			ID:        "default",
			Name:      "Default Preset",
			Keys:      []string{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = AppPresets.Put("default", defaultPreset)
		if err != nil {
			panic("Failed to create default app preset: " + err.Error())
		}
	}
}
