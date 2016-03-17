package algoliaUtil

import (
  "errors"
  "log"
  "os"

  "github.com/algolia/algoliasearch-client-go/algoliasearch"
  "github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
)

const exportAbortError string = `Export aborted: `

var (
  apiKey    string
  appId     string
  client    *algoliasearch.Client
  pageIndex *algoliasearch.Index
)

func init() {
  apiKey    = os.Getenv("ALGOLIASEARCH_API_KEY")
  appId     = os.Getenv("ALGOLIASEARCH_APPLICATION_ID")
}

func ExportPageRecords() error {
  if appId == "" || apiKey == "" {
    return errors.New("Missing Algolia API keys")
  }

  client := algoliasearch.NewClient(appId, apiKey)

  // Check if we need to export all page records (only checks to see if Algolia
  // has a page index set up, does not check for differences in postgres db and algolia index)
  resp, err := client.ListIndexes()
  if err != nil { return errors.New(exportAbortError + err.Error()) }
  indexBlob := resp.(map[string]interface{})
  itemBlob := indexBlob["items"].([]interface{})
  found := false
  for _, value := range itemBlob {
    item := value.(map[string]interface{})
    if item["name"] == "Pages" {
      found = true
    }
  }

  if found {
    return nil
  }

  // Start export
  pageIndex := client.InitIndex("Pages")

  log.Println("Exporting records to Algolia servers...")

  pages, err := databaseActions.GetPages()
  if err != nil { return errors.New(exportAbortError + err.Error()) }

  objects := make([]interface{}, len(pages))
  for ind, page := range pages {
    object := make(map[string]interface{})
    object["title"] = page.Title
    object["category"] = page.Category
    object["objectID"] = page.PageSlug
    object["address"] = page.Address
    object["website"] = page.Website
    object["dateCreated"] = page.DateCreated
    objects[ind] = object
  }
  resp, err = pageIndex.AddObjects(objects)
  if err != nil { return errors.New(exportAbortError + err.Error()); }
  pageIndex.WaitTask(resp)

  // Set ranking information
  settings := make(map[string]interface{})
  settings["attributesToIndex"] = []string{"title", "category"}
  settings["ranking"] = []string{"words", "desc(title)", "desc(category)"}
  resp, err = pageIndex.SetSettings(settings)
  if err != nil { return errors.New(exportAbortError + err.Error()); }
  pageIndex.WaitTask(resp)

  log.Println("Finished export")
  return err
}

func ExportNewPageRecord(page common.Page) (err error) {
  resp, err := pageIndex.AddObject(page)
  if err != nil { return errors.New(exportAbortError + err.Error()); }
  pageIndex.WaitTask(resp)
  return
}

func ExportUpdatedPage(page common.Page) (err error) {
  resp, err := pageIndex.UpdateObject(page)
  if err != nil { return errors.New(exportAbortError + err.Error()); }
  pageIndex.WaitTask(resp)
  return
}

func DeleteExportedPage(objectId string) error {
  resp, err := pageIndex.DeleteObject(objectId)
  if err != nil { return errors.New(exportAbortError + err.Error()); }
  pageIndex.WaitTask(resp)
  return nil
}
