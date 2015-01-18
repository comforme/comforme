package static

const settingsJS = `/* ---------- Communities Settings ---------- */
function registerCommunityCheckboxes() 
{
	$(".communityCheckbox").on
	(
		"click",
		function()
		{
			console.log( "Sending AJAX Request: Name: " + this.name + ", Value: " + this.value + ", Checked: " + this.checked );
			if(this.checked) {
				action = "addCommunity";
			} else {
				action = "removeCommunity";
			}
			
			$.post(
				"/ajax/" + action,
				{ "communityid": this.name }
			).done
			(
				function(data)
				{
					console.log(data);
					// TODO: Notify user of result.
				}
			);
		}
	);
}

function logoutOtherSessions(clickedButton)
{
	$.post("/ajax/logoutOtherSessions").done
	(
		function(data)
		{
			console.log(data);
			if(typeof data.number != "undefined")
			{
				$("#numOpenSessions").html("0");
			}
		}
	)
}

$(document).ready(registerCommunityCheckboxes);
`
