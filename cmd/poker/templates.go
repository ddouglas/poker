package main

import "survey/internal/templates"

var templateFiles = []*templates.Template{
	{
		Name: "style-css",
		Path: "style.css",
	},
	{
		Name: "partial-top",
		Path: "partials/top.html",
	},
	{
		Name: "partial-bottom",
		Path: "partials/bottom.html",
	},
	{
		Name: "partial-navbar",
		Path: "partials/navbar.html",
	},

	{
		Name: "page-home",
		Path: "pages/home.html",
	},
}
