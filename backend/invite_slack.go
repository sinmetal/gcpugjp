package backend

import "net/http"

// SlackInviteHandler is Slack Invite URLへRedirectする
func SlackInviteHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://join.slack.com/t/gcpug/shared_invite/enQtNDMzNjM1Njc3NDkyLWViZmUyNDRhZTJmZTA1YjNhOTUzNDZiZDE5MjRiMzQxN2I1ZTkxMDZjMDcwMDEwOWE0NjFhOGIzZGQzYjhkNWQ", http.StatusFound)
}
