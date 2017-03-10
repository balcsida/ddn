package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/djavorszky/ddn/common/model"
)

// Page is a struct holding the data to be displayed on the welcome page.
type Page struct {
	Connectors   *map[string]model.Connector
	AnyOnline    bool
	Title        string
	Pages        map[string]string
	ActivePage   string
	Message      string
	MessageType  string
	User         string
	HasUser      bool
	HasEntry     bool
	Databases    []model.DBEntry
	HasDatabases bool
	Ext62        model.PortalExt
	ExtDXP       model.PortalExt
}

func loadPage(w http.ResponseWriter, r *http.Request, pages ...string) {

	page := Page{
		Connectors: &registry,
		AnyOnline:  len(registry) > 0,
		Title:      getTitle(r.URL.Path),
		Pages:      getPages(),
		ActivePage: r.URL.Path,
	}

	userCookie, err := r.Cookie("user")
	if err != nil || userCookie.Value == "" {
		// if there's an err, it can only happen if there is no cookie.
		toLoad := []string{"base", "head", "nav", "login"}
		tmpl, err := buildTemplate(toLoad...)
		if err != nil {
			panic(err)
		}

		err = tmpl.ExecuteTemplate(w, "base", page)
		if err != nil {
			panic(err)
		}
		return
	}

	page.User = userCookie.Value
	page.HasUser = true

	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if flashes := session.Flashes("success"); len(flashes) > 0 {
		page.Message = flashes[0].(string)
		page.MessageType = "success"

		id := session.Values["id"].(int64)

		page.HasEntry = true
		entry := db.entryByID(id)

		page.ExtDXP = portalExt(entry, true)
		page.Ext62 = portalExt(entry, false)

	} else if flashes := session.Flashes("fail"); len(flashes) > 0 {
		page.Message = flashes[0].(string)
		page.MessageType = "danger"
	} else if flashes := session.Flashes("debug"); len(flashes) > 0 {
		page.Message = flashes[0].(string)
		page.MessageType = "success"
	} else {
		page.Message = ""
	}

	/*
		// DEBUG:
		if !page.HasEntry {
			page.HasEntry = true
			entry := db.entryByID(1)

			page.ExtDXP = portalExt(entry, true)
			page.Ext62 = portalExt(entry, false)
		}
	*/

	session.Save(r, w)

	if pages[0] == "home" {
		pages = append(pages, "databases")

		page.Databases, _ = db.listWhere(clause{"creator", page.User})
		if len(page.Databases) != 0 {
			page.HasDatabases = true
		}
	}

	toLoad := []string{"base", "head", "nav", "connectors", "properties"}
	toLoad = append(toLoad, pages...)

	tmpl, err := buildTemplate(toLoad...)
	if err != nil {
		panic(err)
	}

	err = tmpl.ExecuteTemplate(w, "base", page)
	if err != nil {
		panic(err)
	}
}

func buildTemplate(pages ...string) (*template.Template, error) {
	var templates []string
	for _, page := range pages {
		templates = append(templates, fmt.Sprintf("web/html/%s.html", page))
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return nil, fmt.Errorf("parsing templates failed: %s", err.Error())
	}

	return tmpl, nil
}

func getTitle(page string) string {
	return getPages()[page]
}

func getPages() map[string]string {
	pages := make(map[string]string)

	pages["/"] = "Home"
	pages["/createdb"] = "Create database"
	pages["/importdb"] = "Import database"

	return pages
}
