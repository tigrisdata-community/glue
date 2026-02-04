package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
)

//go:embed site.css
var siteCSS string

type cssWriter struct {
	styles string
}

func (cw cssWriter) Render(_ context.Context, w io.Writer) error {
	fmt.Fprintln(w, "<style>")
	fmt.Fprintln(w, cw.styles)
	fmt.Fprintln(w, "</style>")

	return nil
}
