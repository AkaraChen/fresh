package main

import (
	"os"
	"sync"

	"github.com/akarachen/taze-go/pkg/structs"
	"github.com/akarachen/taze-go/pkg/workspace"
	"github.com/pterm/pterm"
	"github.com/tidwall/gjson"
)

var (
	table = [][]string{
		{"Name", "Current Version", "Lastest Version", "Devdeps"},
	}
	wg            sync.WaitGroup
	UpdateSchemas []structs.UpdateSchema
)

func checkPackage(root string) {
	json, _ := os.ReadFile(root)
	name := gjson.GetBytes(json, "name").String()
	data := structs.CheckUpdate(json)
	pterm.Success.Printfln("Check " + name + " done.")
	table = append(table, data.ToTable()...)
	if data.Count > 0 {
		UpdateSchemas = append(UpdateSchemas, structs.UpdateSchema{File: root, Json: json, VersionMap: data})
	}
	wg.Done()
}

func main() {
	roots := workspace.CheckWorkSpace(currentPath)
	if len(roots) == 0 {
		pterm.Error.Printfln("No package founded.")
		os.Exit(0)
	}
	spinner, _ := pterm.DefaultSpinner.Start("Checking for update...")
	for _, root := range roots {
		wg.Add(1)
		go checkPackage(root)
	}
	wg.Wait()
	spinner.Success("All done.")
	if len(table) != 1 {
		pterm.DefaultTable.WithHasHeader().WithData(table).WithBoxed().Render()
	} else {
		pterm.Info.Printfln("Your deps are up to date.")
	}
	if len(UpdateSchemas) > 0 {
		result, _ := pterm.DefaultInteractiveConfirm.Show("Would you like to write to file?")
		if result {
			for _, schema := range UpdateSchemas {
				schema.Update()
			}
		}
	}
}
