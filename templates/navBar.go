package templates

const NavBar = `
	<nav class="top-bar" data-topbar>
		<ul class="title-area">
			<li class="name"></li>
			<li class="toggle-topbar menu-icon">
				<a href="#">Menu <span class="icon-menu"></span></a>
			</li>
		</ul>
		<section class="top-bar-section">
			<ul class="title-area">
				<li class="name">
					<h1><a href="/">{{.siteName}}</a></h1>
				</li>
			</ul>
			<ul class="right">
				<li>
					<a href="/" title="Home"><i class="fi-home"><span class="show-for-small-only"> Home</span></i></a>
				</li>
				<li>
					<a href="/newPage" title="Add Resource"><i class="fi-page-add"><span class="show-for-small-only"> Add Resource</span></i></a>
				</li>
				<li>
					<a href="/settings" title="Settings"><i class="fi-widget"><span class="show-for-small-only"> Settings</span></i></a>
				</li>
				<li>
					<a href="/logout" title="Log Out"><i class="fi-power"><span class="show-for-small-only"> Logout</span></i></a>
				</li>
			</ul>
		</section>
	</nav>
`

const NavlessBar = `
	<nav class="top-bar" data-topbar>
		<ul class="title-area">
			<li class="name"></li>
			<li class="toggle-topbar menu-icon">
				<a href="#">Menu <span class="icon-menu"></span></a>
			</li>
		</ul>
		<section class="top-bar-section">
			<ul class="title-area">
				<li class="name">
					<h1><a href="/">{{.siteName}}</a></h1>
				</li>
			</ul>
		</section>
	</nav>
`
