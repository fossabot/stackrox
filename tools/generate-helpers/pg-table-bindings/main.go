package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	_ "github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/tools/generate-helpers/pg-table-bindings/walker"
	"golang.org/x/tools/imports"
	_ "embed"
)

//go:embed store.go.tpl
var storeFile string

type properties struct {
	Type           string
	Table          string
	RegisteredType string
	NoKeyField     bool
	UniqKeyFunc    string
	Cache          bool
	TrackIndex     bool
}

func main() {
	c := &cobra.Command{
		Use: "generate store implementations",
	}

	var props properties
	c.Flags().StringVar(&props.Type, "type", "", "the (Go) name of the object")
	utils.Must(c.MarkFlagRequired("type"))

	c.Flags().StringVar(&props.RegisteredType, "registered-type", "", "the type this is registered in proto as storage.X")

	c.Flags().StringVar(&props.Table, "table", "", "the logical table of the objects")
	utils.Must(c.MarkFlagRequired("table"))

	c.Flags().BoolVar(&props.NoKeyField, "no-key-field", false, "whether or not object contains key field. If no, then to key function is not applied on object")
	c.Flags().StringVar(&props.UniqKeyFunc, "uniq-key-func", "", "when set, unique key constraint is added on the object field retrieved by the function")

	c.RunE = func(*cobra.Command, []string) error {
		typ := fmt.Sprintf("storage.%s", props.Type)
		if props.RegisteredType != "" {
			typ = fmt.Sprintf("storage.%s", props.RegisteredType)
		}
		fmt.Println("Generating for", typ)
		mt := proto.MessageType(typ)
		props.Table = strings.TrimPrefix(mt.Elem().String(), "storage.")
		table := walker.Walk(mt, props.Table)

		insertion := generateInsertFunctions(table)

		tableCreationQueries := createTables(table)
		var count int
		for _, t := range tableCreationQueries {
			if strings.HasPrefix(t, "create table") {
				count++
			}
		}
		fmt.Println("Number of tables", count)

		t := template.Must(template.New("insertion").Parse(insertion))
		buf := bytes.NewBuffer(nil)
		if err := t.Execute(buf, map[string]interface{}{"ExecutePrefix": "tx.Exec(context.Background(),"}); err != nil {
			return err
		}
		singleInsert := buf.String()

		t = template.Must(template.New("insertion").Parse(insertion))
		buf = bytes.NewBuffer(nil)
		if err := t.Execute(buf, map[string]interface{}{"ExecutePrefix": "batch.Queue(", "ExecuteUnchecked": true}); err != nil {
			return err
		}
		multiInsert := buf.String()

		templateMap := map[string]interface{}{
			"Type":       props.Type,
			"Bucket":     props.Table,
			"NoKeyField": props.NoKeyField,
			//"UniqKeyFunc": props.UniqKeyFunc,
			"Table": props.Table,

			"FlatInsertion":            singleInsert,
			"FlatTableCreationQueries": tableCreationQueries,
			"FlatMultiInsert":          multiInsert,
			"SingleTable":              count == 1,

			"TopLevelTable": table,
		}

		t = template.Must(template.New("gen").Funcs(funcMap).Funcs(sprig.TxtFuncMap()).Parse(autogenerated + storeFile))
		buf = bytes.NewBuffer(nil)
		if err := t.Execute(buf, templateMap); err != nil {
			return err
		}
		formatted, err := imports.Process("store.go", buf.Bytes(), nil)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile("store.go", formatted, 0644); err != nil {
			return err
		}

		return nil
	}
	if err := c.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
