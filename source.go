package main

import (
    "os"
    "strings"
    "text/template"
)

type Source struct {
    Name     string
    Template template.Template
}

func (f Source) generate(path string, definition interface{}) error {
    wr, err := os.Create(strings.Join([]string{path, f.Name}, "/"))
    if err != nil {
        return err
    }

    defer wr.Close()
    return f.Template.Execute(wr, definition)
}
