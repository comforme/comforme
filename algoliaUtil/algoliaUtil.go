package algoliaUtil

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/comforme/comforme/common"
    "github.com/comforme/comforme/databaseActions"
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

func ExportPageRecords() error {
	if appId == "" || apiKey == "" {
		return errors.New("Missing Algolia API keys")
	}

	client := algoliasearch.NewClient(appId, apiKey)
    pages, err := databaseActions.GetPages()
    if err != nil {
        log.Println("algoliaUtil: Error retrieving pages from database")
        return err
    }

	// Check if we need to export all page records
	resp, err := client.ListIndexes()
	if err != nil {
		return errors.New(exportAbortError + err.Error())
	}
	indexBlob, ok := resp.(map[string]interface{})
    if !ok {
        fmt.Errorf("algoliaUtil: ExportPageRecords: Unable to extract index")
        return err
    }
	itemBlob, ok := indexBlob["items"].([]interface{})
    if !ok {
        fmt.Errorf("algoliaUtil: ExportPageRecords: Unable to extract items from index")
        return err
    }
	for _, value := range itemBlob {
		item, ok := value.(map[string]interface{})
        if !ok {
            fmt.Errorf("algoliaUtil: ExportPageRecords: Unable to extract item values")
            return err
        }
		if item["name"] == "Pages" {
			numOfEntries, ok := item["entries"].(float64)
            if !ok {
                fmt.Errorf("algoliaUtil: ExportPageRecords: Unable to extract number of entries in index")
                return err
            }
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
