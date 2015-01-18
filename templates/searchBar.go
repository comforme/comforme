package templates

const SearchBar = `
	<form method="post" action="/search">
		<div class="row collapse">
			<div class="small-10 columns">
				<input type="text" placeholder="Page Search" name="page-search" id="page-search-textbox">
			</div>
			<div class="small-2 columns">
				<button type="submit" class="button postfix">Search</button>
			</div>
		</div>
	</form>
`
