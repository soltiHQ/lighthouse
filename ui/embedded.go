package ui

import "embed"

// Templates contain the embedded templates.
//
//go:embed templates/*
var Templates embed.FS

// Static contains the embedded static assets.
//
//go:embed static/*
var Static embed.FS
