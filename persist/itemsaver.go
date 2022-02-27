package persist

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"jasper.com/gospider/config"
	"jasper.com/gospider/engine"
)

var (
	esClient *elastic.Client
	DB       *gorm.DB

	defaultWorkerCount = 1

	saveDBType = ""
)

type Option struct {
	// which db to save
	SaveDBType string
	// how many workers to save items
	WorkerCount int
}

func initESClient() {
	var err error
	esClient, err = elastic.NewClient(elastic.SetSniff(false)) // Must turn off sniff in docker
	if err != nil {
		panic(err)
	}
}

func initMySQL() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MYSQL_USER,
		config.MYSQL_PASSWORD,
		config.MYSQL_HOST,
		config.MYSQL_PORT,
		config.MYSQL_DB)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error,
			Colorful:      true,
		},
	)

	DB, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: newLogger,
		})
	if err != nil {
		panic(err)
	}
}

func ItemSaver(opt Option) chan engine.Item {
	saveDBType = opt.SaveDBType
	if opt.SaveDBType == config.DB_TYPE_ES {
		initESClient()
	} else if opt.SaveDBType == config.DB_TYPE_MYSQL {
		initMySQL()
	}

	out := make(chan engine.Item)

	if opt.WorkerCount <= 0 {
		opt.WorkerCount = defaultWorkerCount
	}

	for i := 0; i < opt.WorkerCount; i++ {
		worker(out, i)
	}

	return out
}

func worker(out chan engine.Item, workerIndex int) {
	go func() {
		for {
			item := <-out
			err := save(item)
			if err != nil {
				log.Printf("[ERROR][worder:%02d] ItemSaver: error saving item -> index: %s, id: %s, err: %v\n", workerIndex, item.Index, item.ID, err)
				return
			}
			log.Printf("[worder:%02d] ItemSaver: saved [%s] item %v\n", workerIndex, item.Index, item.ID)
		}
	}()
}

func save(item engine.Item) error {
	switch saveDBType {
	case config.DB_TYPE_ES:
		return saveToES(item)
	case config.DB_TYPE_MYSQL:
		return saveToMySQL(item)
	}
	return nil
}

func saveToMySQL(item engine.Item) error {
	// TODO: save to mysql

	return nil
}

func saveToES(item engine.Item) error {
	_, err := esClient.Index().
		Index(item.Index).
		Id(item.ID).
		BodyJson(item).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
