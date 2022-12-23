package structs

import (
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/akarachen/fresh/pkg/request"
	"github.com/pterm/pterm"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Package struct {
	Name    string
	Version semver.Version
}

var cacheMap map[string]semver.Version = make(map[string]semver.Version)

func GetByName(name string) (Package, error) {
	if value, exist := cacheMap[name]; exist {
		return Package{Name: name, Version: value}, nil
	}
	var builder strings.Builder
	builder.WriteString("http://registry.npmjs.com/")
	builder.WriteString(name)
	resp, err := request.Client.R().Get(builder.String())
	if err != nil {
		return Package{}, errors.New(err.Error())
	}
	body := string(resp.Body())
	pkgVersionRaw := gjson.Get(body, "dist-tags.latest").Str
	version, _ := semver.NewVersion(pkgVersionRaw)
	pkg := Package{Name: name, Version: *version}
	cacheMap[name] = *version
	return pkg, nil
}

type UpdateMap struct {
	Deps    map[string]versionMap
	Devdeps map[string]versionMap
	Count   int
}

func (updateMap UpdateMap) ToTable() pterm.TableData {
	table := pterm.TableData{}
	for key, value := range updateMap.Deps {
		table = append(table, []string{key, value.Current.Original(), value.Lastest.Original(), "False"})
	}
	for key, value := range updateMap.Devdeps {
		table = append(table, []string{key, value.Current.Original(), value.Lastest.Original(), "True"})
	}
	return table
}

func CheckUpdate(json []byte) UpdateMap {
	depsMap := checkUpdateByfield(json, "dependencies")
	devdepsMap := checkUpdateByfield(json, "devDependencies")
	length := len(depsMap) + len(devdepsMap)
	return UpdateMap{Deps: depsMap, Devdeps: devdepsMap, Count: length}
}

type versionMap struct {
	Current semver.Version
	Lastest semver.Version
}

func checkUpdateByfield(json []byte, field string) map[string]versionMap {
	deps := gjson.GetBytes(json, field).Map()
	updateMap := make(map[string]versionMap)
	var wg sync.WaitGroup
	for key, value := range deps {
		current, _ := semver.NewVersion(value.String())
		wg.Add(1)
		go getUpdateByName(key, *current, updateMap, &wg)
	}
	wg.Wait()
	return updateMap
}

func getUpdateByName(name string, current semver.Version, m map[string]versionMap, wg *sync.WaitGroup) {
	info, _ := GetByName(name)
	if info.Version.GreaterThan(&current) {
		m[name] = versionMap{Current: current, Lastest: info.Version}
	}
	wg.Done()
}

type UpdateSchema struct {
	File       string
	Json       []byte
	VersionMap UpdateMap
}

func (updateSchema UpdateSchema) Update() {
	var bytes []byte = updateSchema.Json
	for key, value := range updateSchema.VersionMap.Deps {
		bytes, _ = sjson.SetBytes(bytes, "dependencies."+key, value.Lastest.Original())
	}
	for key, value := range updateSchema.VersionMap.Devdeps {
		bytes, _ = sjson.SetBytes(bytes, "devDependencies."+key, value.Lastest.Original())
	}
	os.WriteFile(updateSchema.File, bytes, os.ModeDevice)
}
