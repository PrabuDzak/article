package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/olivere/elastic"
	"github.com/subosito/gotenv"

	"github.com/prabudzak/article/app/restapi"
	"github.com/prabudzak/article/event"
	"github.com/prabudzak/article/event/memory"
	"github.com/prabudzak/article/service/article"
	articleindexer "github.com/prabudzak/article/service/article/elasticsearch"
	articledb "github.com/prabudzak/article/service/article/mysql"
	articlecache "github.com/prabudzak/article/service/article/redis"
)

func main() {
	gotenv.Load()
	ctx := context.Background()

	sqlCfg := mysql.NewConfig()
	sqlCfg.Addr = fmt.Sprintf("%s:%s", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))
	sqlCfg.User = os.Getenv("MYSQL_USERNAME")
	sqlCfg.Passwd = os.Getenv("MYSQL_PASSWORD")
	sqlCfg.DBName = os.Getenv("MYSQL_DATABASE")
	sqlCfg.ParseTime = true

	dbDriver, err := mysql.NewConnector(sqlCfg)
	if err != nil {
		log.Fatalln(err)
	}

	conn := sql.OpenDB(dbDriver)
	err = conn.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	redisClient := redis.NewClient(&redis.Options{})
	err = redisClient.Ping().Err()
	if err != nil {
		log.Fatalln(err)
	}

	esClient, err := elastic.NewClient(
		elastic.SetURL(os.Getenv("ELASTICSEARCH_URL")),
		elastic.SetHttpClient(&http.Client{}),
	)
	if err != nil {
		log.Fatalln(err)
	}

	articleDatabase := articledb.NewArticleDatabase(conn)
	articleCache := articlecache.NewArticleCache(redisClient)
	articleIndexer := articleindexer.NewArticleIndexer(esClient, os.Getenv("ELASTICSEARCH_ARTICLE_INDEX"))
	articleService := article.NewArticleService(articleDatabase, articleCache, articleIndexer)

	dispatcher := memory.NewDispatcher()
	dispatcher.Start()
	event.SetDispatcher(dispatcher)

	dispatcher.AddSubscriber(ctx, event.ArticleCreated{}, articleService.SubscriberCacheArticle)

	router := restapi.New(articleService)

	log.Printf("listening in %s\n", os.Getenv("PORT"))
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")), router.Router())
	if err != nil {
		log.Fatalln(err)
	}

}
