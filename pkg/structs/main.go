package structs

import (
	"errors"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/akarachen/taze-go/pkg/request"
	"github.com/tidwall/gjson"
)

type Package struct {
	Name    string
	Version semver.Version
}

func GetByName(ch chan Package, name string) error {
	var builder strings.Builder
	builder.WriteString("http://registry.npmjs.com/")
	builder.WriteString(name)
	resp, err := request.Client.R().Get(builder.String())
	if err != nil {
		return errors.New(err.Error())
	}
	body := string(resp.Body())
	pkgName := gjson.Get(body, "name").Str
	pkgVersionRaw := gjson.Get(body, "dist-tags.latest").Str
	version, _ := semver.NewVersion(pkgVersionRaw)
	ch <- Package{Name: pkgName, Version: *version}
	return nil
}
