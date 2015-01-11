package search

import (
	"html/template"
	"net/http"

	"github.com/comforme/comforme/common"
	// "github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var searchTemplate *template.Template

func init() {
	searchTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(searchTemplate.New("nav").Parse(templates.NavBar))
	template.Must(searchTemplate.New("searchBar").Parse(templates.SearchBar))
	template.Must(searchTemplate.New("content").Parse(searchTemplateText))
}

func SearchHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	if req.Method == "POST" {
		// TODO uncomment when put to use
		//query := req.PostFormValue("query")
	}

	// TODO: Add template and compile it.
	common.ExecTemplate(searchTemplate, res, data)
}

const searchTemplateText = `
    <div class="content">
        <div class="row">
            <div class="column">
                <h1>Search</h1>
                {{template "searchBar" .}}
            </div>
        </div>
        <div class="row">
            <div class="columns">
                <h2>Lorem Hipsum</h2>
                <div>
                    <p>Odd Future Bushwick irony, Neutra artisan chambray forage Banksy skateboard Schlitz hoodie cold-pressed sustainable brunch. Freegan Etsy mixtape, selvage small batch pop-up distillery VHS. IPhone flexitarian tousled, letterpress Pitchfork readymade cornhole. Shabby chic irony skateboard, swag lumbersexual DIY Portland ethical Williamsburg forage farm-to-table meditation. Intelligentsia quinoa Odd Future semiotics hella Wes Anderson fap, typewriter Austin cliche meh lomo tattooed. Trust fund yr cronut, fap mumblecore viral Blue Bottle readymade. Sriracha street art Thundercats, PBR deep v trust fund fashion axe... <a href="">Continue Reading</a></p>
                </div>
            </div>
            <div class="columns">
                <h2>So Obscure, You Probably Haven't Heard of It</h2>
                <div>
                    <p>Photo booth Portland hoodie, retro sartorial ugh Thundercats tofu selfies Williamsburg meditation PBR pour-over bespoke. Meh heirloom kogi, trust fund pug messenger bag migas. Bicycle rights four dollar toast kale chips biodiesel. Chia umami Helvetica Brooklyn. Wolf iPhone Helvetica keffiyeh, hoodie keytar pop-up normcore Neutra mlkshk. Pour-over messenger bag Thundercats, swag mumblecore plaid 90's sustainable wolf mixtape hashtag. Pork belly fap occupy, Wes Anderson polaroid migas keffiyeh mustache single-origin coffee Intelligentsia actually meggings Thundercats pug... <a href="">Continue Reading</a></p>
                </div>
            </div>
        </div>
    </div>
</div>
`
