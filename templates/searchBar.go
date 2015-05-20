package templates

const SearchBar = `
	<form method="get" action="/search">
		<div class="row collapse">
			<div class="small-10 columns">
				<input type="text" placeholder="Page Search" name="q" id="page-search-textbox">
			</div>
			<div class="small-2 columns">
				<button type="submit" class="button postfix">Search</button>
			</div>
		</div>
	</form>
`
