package templates

const Communities = `
					<h2>Your Communities</h2>
					<h6>Check all that apply.</h6>
					<div class="row">{{range $col_number, $communitiesCol := $.communitiesCols}}
						<div class="large-3 medium-6 small-12 columns left">{{range $line_number, $community := $communitiesCol}}
							<div>
								<label>
									<input class="communityCheckbox" type="checkbox" name="{{$community.Id}}"{{if eq $community.IsMember true}} checked="checked"{{end}} value="{{$community.Name}}">
									{{$community.Name}}
								</label>
							</div>{{end}}
						</div>{{end}}
					</div>
`
