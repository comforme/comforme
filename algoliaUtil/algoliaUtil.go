package algoliaUtil

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/comforme/comforme/common"
)

const exportAbortError string = `Export aborted: `

var (
	apiKey string
	appId  string
)

func init() {
	apiKey = os.Getenv("ALGOLIASEARCH_API_KEY")
	appId = os.Getenv("ALGOLIASEARCH_APPLICATION_ID")
}

func ExportPageRecords(pages []common.Page) error {
	if appId == "" || apiKey == "" {
		return errors.New("Missing Algolia API keys")
	}

	client := algoliasearch.NewClient(appId, apiKey)

	// Check if we need to export all page records
	resp, err := client.ListIndexes()
	if err != nil {
		return errors.New(exportAbortError + err.Error())
	}
	indexBlob := resp.(map[string]interface{})
	itemBlob := indexBlob["items"].([]interface{})
	for _, value := range itemBlob {
		item := value.(map[string]interface{})
		if item["name"] == "Pages" {
			numOfEntries := item["entries"].(float64)
			if numOfEntries == float64(len(pages)) {
				log.Println("Index 'Pages' already exists, aborting export.")
				return nil
			}
		}
	}

	// Start export
	pageIndex := client.InitIndex("Pages")

	if len(pages) == 0 {
		log.Println("No pages to export.")
		return nil
	}

	log.Println("Contructing page objects for transport...")
	objects := make([]interface{}, len(pages))
	for ind, page := range pages {
		object := pageToObject(page)
		objects[ind] = object
	}

	fmt.Println("Adding objects to 'Pages' index")
	resp, err = pageIndex.AddObjects(interface{}(objects))
	if err != nil {
		return errors.New(exportAbortError + err.Error())
	}
	pageIndex.WaitTask(resp)

	// Set ranking information

	settings := make(map[string]interface{})
	settings["attributesToIndex"] = []string{"title", "category"}
	settings["ranking"] = []string{"words", "desc(title)", "desc(category)"}
	resp, err = pageIndex.SetSettings(settings)
	if err != nil {
		return errors.New(exportAbortError + err.Error())
	}
	pageIndex.WaitTask(resp)

	return err
}

func ExportPageRecord(page common.Page) (err error) {
	if appId == "" || apiKey == "" {
		return errors.New("Missing Algolia API keys")
	}

	client := algoliasearch.NewClient(appId, apiKey)

	log.Println("Exporting page:" + page.Title + " to algolia servers..")
	object := pageToObject(page)
	pageIndex := client.InitIndex("Pages")
	resp, err := pageIndex.AddObject(object)
	if err != nil {
		return errors.New(exportAbortError + err.Error())
	}
	pageIndex.WaitTask(resp)
	log.Println("Done exporting.")
	return
}

func DeleteExportedPage(objectId string) error {
	client := algoliasearch.NewClient(appId, apiKey)
	pageIndex := client.InitIndex("Pages")
	resp, err := pageIndex.DeleteObject(objectId)
	if err != nil {
		return errors.New(exportAbortError + err.Error())
	}
	pageIndex.WaitTask(resp)
	return nil
}

func pageToObject(page common.Page) map[string]interface{} {
	object := make(map[string]interface{}, 4)
	object["objectID"] = page.PageSlug + "-" + page.Category
	object["title"] = page.Title
	object["category"] = page.Category
	object["dateCreated"] = page.DateCreated
	return object
}
