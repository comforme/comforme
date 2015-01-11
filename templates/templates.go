package templates

const SiteLayout = `<!DOCTYPE html>
<html>
<head>
	<link href="https://cdnjs.cloudflare.com/ajax/libs/foundation/5.5.0/css/normalize.min.css" rel="stylesheet" type="text/css" />
	<link href="https://cdnjs.cloudflare.com/ajax/libs/foundation/5.5.0/css/foundation.min.css" rel="stylesheet" type="text/css" />
	<script src="https://cdnjs.cloudflare.com/ajax/libs/foundation/5.5.0/js/vendor/jquery.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/foundation/5.5.0/js/foundation.min.js"></script>
	<meta charset="utf-8" />
	<title>ComFor.Me - {{.pageTitle}}</title>
    <link rel="stylesheet" href="http://spyrosoft.bitbucket.org/css/style.css" />
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/foundicons/3.0.0/foundation-icons.css" />
	<script scr="https://cdnjs.cloudflare.com/ajax/libs/foundation/5.5.0/js/login.js"></script>
</head>
<body>
    {{ template "nav" .}}
    {{ template "content" .}}
	<script>$(document).foundation();</script>
</body>
</html>
`

const NavBar = `
    <nav class="top-bar" data-topbar>
		<ul class="title-area">
			<li class="name"></li>
			<li class="toggle-topbar menu-icon">
				<a href="#">Menu <span class="icon-menu"></span></a>
			</li>
		</ul>
		<section class="top-bar-section">
			<ul class="right">
				<li>
					<a href="/"><i class="fi-home"></i></a>
				</li>
				<li>
					<a href="/page"><i class="fi-page-add"></i></a>
				</li>
				<li>
					<a href="/profile"><i class="fi-widget"></i></a>
				</li>
				<li>
					<a href="/"><i class="fi-power"></i></a>
				</li>
			</ul>
		</section>
	</nav>
`

const SearchBar = `
	<div class="row">
		<div class="column">
			<form method="post" action="/">
				<div class="row collapse">
					<div class="small-10 columns">
						<input type="text" placeholder="Page Search" name="page-search" id="page-search-textbox">
					</div>
				<div class="small-2 columns">
					<button type="submit" class="button postfix">Submit</button>
				</div>
				</div>
			</form>
            <script>document.getElementById('page-search-textbox').select();</script>
		</div>
	</div>
`
