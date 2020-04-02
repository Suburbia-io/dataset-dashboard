package dbgen

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

var (
	fpJoin = filepath.Join
)

func Generate(
	tables []Table,
	pkgName string,
	tmplDir string,
	outDir string,
) {
	// Parse templates.
	tmpl := template.Must(
		template.ParseGlob(
			filepath.Join(tmplDir, "*.tmpl")))

	execTmpl(tmpl, "common.go", pkgName, fpJoin(outDir, "common.go"))
	execTmpl(tmpl, "main_test.go", pkgName, fpJoin(outDir, "main_test.go"))

	for _, table := range tables {
		t := table.t
		t.Package = pkgName
		log.Printf("Processing %s...", t.Table)
		execTmpl(tmpl, "table.go", t, fpJoin(outDir, "table-"+t.Table+".go"))
		execTmpl(tmpl, "table_test.go", t, fpJoin(outDir, "table-"+t.Table+"_test.go"))
	}
}

func execTmpl(
	tmpl *template.Template,
	name string,
	ctx interface{},
	outPath string,
) {
	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	if err := tmpl.ExecuteTemplate(f, name, ctx); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}

	if err := exec.Command("goimports", "-w", outPath).Run(); err != nil {
		panic(err)
	}

	if err := exec.Command("gofmt", "-w", outPath).Run(); err != nil {
		panic(err)
	}
}
