package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	packageName = flag.String("package", "model", "name of package")
	fileName    = flag.String("file", "model/columns_gen.go", "name of file")
)

func main() {
	flag.Parse()
	if err := _main(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func _main() error {
	var buf bytes.Buffer

	synonyms := map[string]string{
		"Integer": "Int",
		"Numeric": "Decimal",
		"Real":    "Double",
	}

	types := []string{
		"Invalid",
		"Bit",
		"TinyInt",
		"SmallInt",
		"MediumInt",
		"Int",
		"Integer",
		"BigInt",
		"Real",
		"Double",
		"Float",
		"Decimal",
		"Numeric",
		"Date",
		"Time",
		"Timestamp",
		"DateTime",
		"Year",
		"Char",
		"VarChar",
		"Binary",
		"VarBinary",
		"TinyBlob",
		"Blob",
		"MediumBlob",
		"LongBlob",
		"TinyText",
		"Text",
		"MediumText",
		"LongText",
		"Enum",
		"Set",
	}

	buf.WriteString(`// generated by internal/cmd/gencoltypes/main.go. DO NOT EDIT`)
	buf.WriteString("\n\npackage " + *packageName)
	buf.WriteString("\n\n// ColumnType describes the possible types that a column may take")
	buf.WriteString("\ntype ColumnType int")
	buf.WriteString("\n\n// List of possible ColumnType values")
	buf.WriteString("\nconst (")
	for i, typ := range types {
		buf.WriteString("\nColumnType")
		buf.WriteString(typ)
		if i == 0 {
			buf.WriteString(" ColumnType = iota")
		}
	}
	buf.WriteString("\n\nColumnTypeMax")
	buf.WriteString("\n)")
	buf.WriteString("\n\nfunc (c ColumnType) String() string {")
	buf.WriteString("\nswitch c {")
	for _, typ := range types[1:] {
		buf.WriteString("\ncase ColumnType")
		buf.WriteString(typ)
		buf.WriteByte(':')
		buf.WriteString("\nreturn ")
		buf.WriteString(strconv.Quote(strings.ToUpper(typ)))
	}
	buf.WriteString("\ndefault:")
	buf.WriteString("\nreturn \"(invalid)\"")
	buf.WriteString("\n}")
	buf.WriteString("\n}")

	buf.WriteString("\n\n// SynonymType returns synonym for a given type.")
	buf.WriteString("\n// If the type does not have a synonym then this method returns the receiver itself")
	buf.WriteString("\nfunc (c ColumnType) SynonymType() ColumnType {")
	buf.WriteString("\nswitch c {")
	for from, to := range synonyms {
		buf.WriteString("\ncase ColumnType" + from + ":")
		buf.WriteString("\nreturn ColumnType" + to)
	}
	buf.WriteString("\n}")
	buf.WriteString("\nreturn c")
	buf.WriteString("\n}")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("%s", buf.Bytes())
		return err
	}

	f, err := os.Create(*fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(formatted)
	return nil
}
