package controllers

import "os"

// RepositoryController handles request to /repository
type RepositoryController struct {
	BaseController
}

// Get renders repository page
func (rc *RepositoryController) Get() {
	rc.Data["HarborRegUrl"] = os.Getenv("HARBOR_REG_URL")
	rc.Data["ConsoleWebUrl"] = os.Getenv("CONSOLE_WEB_URL")
	//rc.Data["ConsoleWebUrl"] = "http://localhost:3000"
	rc.Forward("page_title_repository", "repository.htm")
}
