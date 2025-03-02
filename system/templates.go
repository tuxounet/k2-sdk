package system

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func UnTemplateWithGoTemplate(tpl string, data any) (string, error) {

	hashBytes := md5.Sum([]byte(tpl))
	hash := hex.EncodeToString(hashBytes[:])

	var buf bytes.Buffer
	tmpl, err := template.New(hash).Funcs(sprig.FuncMap()).Parse(tpl)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}
