package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	pkgDir string
	outDir string
)

const codeTemplate = `
package main

import (
	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/plugin"
)

func PluginManifest() *plugin.Manifest {
	return plugin.ExportManifest(&plugin.{{.kind}}Manifest{
		Manifest: plugin.Manifest{
			Kind:           plugin.{{.kind}},
			Name:           "{{.name}}",
			Description:    "{{.description}}",
			BuildTime:      "{{.buildTime}}",
			{{if .onInit }}
				OnInit:     {{.onInit}},
			{{end}}
			{{if .onShutdown }}
				OnShutdown: {{.onShutdown}},
			{{end}}
		},
		{{range .export}}
		{{.extPoint}}: {{.impl}},
		{{end}}
	})
}
`

func init() {
	flag.StringVar(&pkgDir, "pkg-dir", "", "plugin package folder path")
	flag.StringVar(&outDir, "out-dir", "", "plugin packaged folder path")
	flag.Usage = usage
}

func usage() {
	log.Printf("Usage: %s --pkg-dir [plugin source pkg folder] --out-dir [plugin packaged folder path]\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Parse()
	if pkgDir == "" || outDir == "" {
		flag.Usage()
	}
	pkgDir, err := filepath.Abs(pkgDir)
	if err != nil {
		log.Printf("unable to resolve absolute representation of package path , %+v\n", err)
		flag.Usage()
	}
	outDir, err := filepath.Abs(outDir)
	if err != nil {
		log.Printf("unable to resolve absolute representation of output path , %+v\n", err)
		flag.Usage()
	}

	if err := Build(pkgDir, outDir); err != nil {
		log.Fatalln(err)
	}
}

func Build(pkgDir, outDir string) error {
	var manifest map[string]interface{}
	_, err := toml.DecodeFile(filepath.Join(pkgDir, "manifest.toml"), &manifest)
	if err != nil {
		return fmt.Errorf("read pkg %s's manifest failure, %w", pkgDir, err)
	}
	manifest["buildTime"] = time.Now().Format("2006.01.02 15:04:05")

	pluginName := manifest["name"].(string)
	tmpl, err := template.New("gen-plugin").Parse(codeTemplate)
	if err != nil {
		return fmt.Errorf("generate code failure during parse template, %w", err)
	}

	genFileName := filepath.Join(pkgDir, filepath.Base(pkgDir)+".gen.go")
	genFile, err := os.OpenFile(genFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700) // #nosec G302
	if err != nil {
		return fmt.Errorf("generate code failure during prepare output file, %w", err)
	}
	defer func() {
		if err := os.Remove(genFileName); err != nil {
			log.Printf("remove tmp file %s failure, please clean up manually at %v", genFileName, err)
		}
	}()

	err = tmpl.Execute(genFile, manifest)
	if err != nil {
		return fmt.Errorf("generate code failure during generating code, %w", err)
	}

	outputFile := filepath.Join(outDir, pluginName+".so")
	ctx := context.Background()
	buildCmd := exec.CommandContext(ctx, "go", "build",
		"-buildmode=plugin",
		"-o", outputFile, pkgDir)
	buildCmd.Dir = pkgDir
	buildCmd.Stderr = os.Stderr
	buildCmd.Stdout = os.Stdout
	buildCmd.Env = append(os.Environ(), "GO111MODULE=on")
	err = buildCmd.Run()
	if err != nil {
		return fmt.Errorf("compile plugin source code failure, %w", err)
	}
	fmt.Printf(`Package "%s" as plugin "%s" success.`+"\nManifest:\n", pkgDir, outputFile)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "\t")
	err = encoder.Encode(manifest)
	if err != nil {
		return fmt.Errorf("print manifest detail failure, err: %w", err)
	}
	return nil
}
