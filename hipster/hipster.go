package hipster

import (
	"html/template"
	//"log"
	"net/http"

	"github.com/comforme/comforme/common"
	//"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var hipsterTemplate *template.Template

func init() {
	hipsterTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(hipsterTemplate.New("nav").Parse(templates.NavBar))
	template.Must(hipsterTemplate.New("hipster").Parse(templates.SearchBar))
	template.Must(hipsterTemplate.New("content").Parse(hipsterTemplateText))
}

func HipsterHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	data["formAction"] = req.URL.Path

	common.ExecTemplate(hipsterTemplate, res, data)
}

const hipsterTemplateText = `
	<div class="content">
		<div class="row">
			<div class="column">
				<h1><a href="/hipster/lorem-hipsum">Lorem Hipsum</a></h1>
				<p>
					Odd Future Bushwick irony, Neutra artisan chambray forage Banksy skateboard Schlitz			hoodie cold-pressed sustainable brunch. Freegan Etsy mixtape, selvage small batch pop-up					 distillery VHS. IPhone flexitarian tousled, letterpress Pitchfork readymade cornhole. Shabby chic	irony skateboard, swag lumbersexual DIY Portland ethical Williamsburg forage farm-to-table				 meditation. Intelligentsia quinoa Odd Future semiotics hella Wes Anderson fap, typewriter Austin	 cliche meh lomo tattooed. Trust fund yr cronut, fap mumblecore viral Blue Bottle readymade.				Sriracha street art Thundercats, PBR deep v trust fund fashion axe.
				</p>
				<p>
					American Apparel bespoke photo booth, ennui distillery Truffaut heirloom 90's brunch.		Stumptown kogi heirloom ugh lo-fi, Pinterest blog Schlitz XOXO cliche slow-carb kale chips master	cleanse Tumblr. Forage food truck Thundercats Blue Bottle vegan polaroid artisan next level				taxidermy. Fashion axe semiotics authentic hashtag. Ennui seitan lomo, pop-up fixie brunch fanny	 pack semiotics readymade Pitchfork Williamsburg forage try-hard. Pork belly Truffaut normcore			organic Brooklyn readymade, meh Austin asymmetrical. Bicycle rights Odd Future synth, butcher			tousled chambray selfies Williamsburg.
				</p>
				<p>
					Tattooed direct trade lo-fi four loko, photo booth mlkshk single-origin coffee cliche		chillwave bespoke dreamcatcher. Pitchfork mlkshk plaid 3 wolf moon small batch umami. High Life		selvage yr selfies YOLO. You probably haven't heard of them Williamsburg +1 Etsy fashion axe			 wayfarers. Photo booth tattooed literally, asymmetrical health goth Austin skateboard slow-carb.	 Fingerstache bitters pug wolf. Distillery gluten-free aesthetic cardigan normcore sriracha.
				</p>
			</div>
		</div>
	</div>
`
