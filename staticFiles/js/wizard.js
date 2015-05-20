$( window ).on( 'hashchange', url_hash_changed );

function url_hash_changed() {
	var url_hash = window.location.hash.substring( 1, window.location.hash.length );
	if ( url_hash.length > 0 ) {
		if ( validate_url_hash ( url_hash ) ) {
			change_navigation_button( url_hash );
		}
	}
}

function validate_url_hash( url_hash ) {
	if ( url_hash.match( /tab-/ ) ) {
		return true;
	} else {
		return false;
	}
}

function change_navigation_button( url_hash ) {
	console.log( url_hash );
	$( '#' + url_hash + "-navigation" ).trigger( 'click' );
}

url_hash_changed();

$(document).foundation({
	tab: {
		callback : function (tab) {
			document.location = tab.context.hash;
		}
	}
});