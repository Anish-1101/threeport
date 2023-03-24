package models

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/dave/jennifer/jen"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"

	"github.com/threeport/threeport/internal/codegen/name"
)

// apiHandlersPath returns the path from the models to the API's internal handlers
// package.
func clientLibPath(packageName string) string {
	return filepath.Join("..", "..", "..", "pkg", "client", packageName)
}

// ClientLib generates the client library code for the API models in a model
// file.
func (cc *ControllerConfig) ClientLib() error {
	pluralize := pluralize.NewClient()
	f := NewFile(cc.PackageName)
	f.HeaderComment("generated by 'threeport-codegen api-model' - do not edit")

	for _, mc := range cc.ModelConfigs {
		// get object by ID
		getByIDFuncName := fmt.Sprintf("Get%sByID", mc.TypeName)
		f.Comment(fmt.Sprintf(
			"%s feteches a %s by ID",
			getByIDFuncName,
			strcase.ToDelimited(mc.TypeName, ' '),
		))
		f.Func().Id(getByIDFuncName).Params(
			Id("id").Uint(),
			Id("apiAddr").Op(",").Id("apiToken").String(),
		).Parens(List(
			Op("*").Id(cc.PackageName).Dot(mc.TypeName),
			Error(),
		)).Block(
			Var().Id(strcase.ToLowerCamel(mc.TypeName)).Qual(
				"github.com/threeport/threeport/pkg/api/v0",
				mc.TypeName,
			),
			Line(),
			Id("response").Op(",").Id("err").Op(":=").Id("GetResponse").Call(
				Line().Qual("fmt", "Sprintf").Call(
					Lit(fmt.Sprintf(
						"%%s/%%s/%s/%%d", pluralize.Pluralize(strcase.ToSnake(mc.TypeName), 2, false),
					)).Op(",").
						Id("apiAddr").Op(",").Id("ApiVersion").Op(",").Id("id"),
				),
				Line().Id("apiToken"),
				Line().Qual("net/http", "MethodGet"),
				Line().New(Qual("bytes", "Buffer")),
				Line().Qual("net/http", "StatusOK"),
				Line(),
			),
			If(Id("err").Op("!=").Nil().Block(
				Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Id("err"),
			)),
			Line(),
			Id("jsonData").Op(",").Id("err").Op(":=").Qual("encoding/json", "Marshal").Call(
				Id("response").Dot("Data").Index(Lit(0)),
			),
			If(Id("err").Op("!=").Nil().Block(
				Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Id("err"),
			)),
			Line(),
			If(
				Id("err").Op("=").Qual("encoding/json", "Unmarshal").Call(
					Id("jsonData").Op(",").Op("&").Id(strcase.ToLowerCamel(mc.TypeName)),
				).Op(";").Id("err").Op("!=").Nil().Block(
					Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Id("err"),
				),
			),
			Line(),
			Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Nil(),
		)
		f.Line()
		// get object by name
		getByNameFuncName := fmt.Sprintf("Get%sByName", mc.TypeName)
		f.Comment(fmt.Sprintf(
			"%s feteches a %s by name",
			getByNameFuncName,
			strcase.ToDelimited(mc.TypeName, ' '),
		))
		f.Func().Id(getByNameFuncName).Params(
			Id("name").Op(",").Id("apiAddr").Op(",").Id("apiToken").String(),
		).Parens(List(
			Op("*").Qual(
				"github.com/threeport/threeport/pkg/api/v0",
				mc.TypeName,
			),
			Error(),
		)).Block(
			Var().Id(
				pluralize.Pluralize(strcase.ToLowerCamel(mc.TypeName), 2, false),
			).Index().Id(cc.PackageName).Dot(mc.TypeName),
			Line(),
			Id("response").Op(",").Id("err").Op(":=").Id("GetResponse").Call(
				Line().Qual("fmt", "Sprintf").Call(
					Lit(fmt.Sprintf(
						"%%s/%%s/%s?name=%%s", pluralize.Pluralize(strcase.ToSnake(mc.TypeName), 2, false),
					)).Op(",").
						Id("apiAddr").Op(",").Id("ApiVersion").Op(",").Id("name"),
				),
				Line().Id("apiToken"),
				Line().Qual("net/http", "MethodGet"),
				Line().New(Qual("bytes", "Buffer")),
				Line().Qual("net/http", "StatusOK"),
				Line(),
			),
			If(Id("err").Op("!=").Nil().Block(
				Return().Op("&").Qual(
					"github.com/threeport/threeport/pkg/api/v0",
					mc.TypeName,
				).Values().Op(",").Id("err"),
			)),
			Line(),
			Id("jsonData").Op(",").Id("err").Op(":=").Qual("encoding/json", "Marshal").Call(
				Id("response").Dot("Data"),
			),
			If(Id("err").Op("!=").Nil().Block(
				Return().Op("&").Qual(
					"github.com/threeport/threeport/pkg/api/v0",
					mc.TypeName,
				).Values().Op(",").Id("err"),
			)),
			Line(),
			If(
				Id("err").Op("=").Qual("encoding/json", "Unmarshal").Call(
					Id("jsonData").Op(",").Op("&").Id(pluralize.Pluralize(strcase.ToLowerCamel(mc.TypeName), 2, false)),
				).Op(";").Id("err").Op("!=").Nil().Block(
					Return().Op("&").Qual(
						"github.com/threeport/threeport/pkg/api/v0",
						mc.TypeName,
					).Values().Op(",").Id("err"),
				),
			),
			Line(),
			Switch().Block(
				Case(Len(Id(pluralize.Pluralize(strcase.ToLowerCamel(mc.TypeName), 2, false))).Op("<").Lit(1)).Block(
					Return().Op("&").Qual(
						"github.com/threeport/threeport/pkg/api/v0",
						mc.TypeName,
					).Values().Op(",").Qual("errors", "New").Call(
						Qual("fmt", "Sprintf").Call(
							Lit("no workload definitions with name %s").Op(",").Id("name"),
						),
					),
				),
				Case(Len(Id(pluralize.Pluralize(strcase.ToLowerCamel(mc.TypeName), 2, false))).Op(">").Lit(1)).Block(
					Return().Op("&").Qual(
						"github.com/threeport/threeport/pkg/api/v0",
						mc.TypeName,
					).Values().Op(",").Qual("errors", "New").Call(
						Qual("fmt", "Sprintf").Call(
							Lit("more than one workload definition with name %s returned").Op(",").Id("name"),
						),
					),
				),
			),
			Line(),
			Return().Op("&").Id(pluralize.Pluralize(strcase.ToLowerCamel(mc.TypeName), 2, false)).
				Index(Lit(0)).Op(",").Nil(),
		)
		f.Line()
		// create object
		createFuncName := fmt.Sprintf("Create%s", mc.TypeName)
		f.Comment(fmt.Sprintf(
			"%s creates a new %s",
			createFuncName,
			strcase.ToDelimited(mc.TypeName, ' '),
		))
		f.Func().Id(createFuncName).Params(
			Id(fmt.Sprintf("json%s", mc.TypeName)).Index().Byte().Op(",").Id("apiAddr").Op(",").Id("apiToken").String(),
		).Parens(List(
			Op("*").Qual(
				"github.com/threeport/threeport/pkg/api/v0",
				mc.TypeName,
			),
			Error(),
		)).Block(
			Var().Id(strcase.ToLowerCamel(mc.TypeName)).Qual(
				"github.com/threeport/threeport/pkg/api/v0",
				mc.TypeName,
			),
			Line(),
			Id("response").Op(",").Id("err").Op(":=").Id("GetResponse").Call(
				Line().Qual("fmt", "Sprintf").Call(
					Lit(fmt.Sprintf(
						"%%s/%%s/%s", pluralize.Pluralize(strcase.ToSnake(mc.TypeName), 2, false),
					)).Op(",").
						Id("apiAddr").Op(",").Id("ApiVersion"),
				),
				Line().Id("apiToken"),
				Line().Qual("net/http", "MethodGet"),
				Line().Qual("bytes", "NewBuffer").Call(Id(
					fmt.Sprintf("json%s", mc.TypeName),
				)),
				Line().Qual("net/http", "StatusCreated"),
				Line(),
			),
			If(Id("err").Op("!=").Nil().Block(
				Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Id("err"),
			)),
			Line(),
			Id("jsonData").Op(",").Id("err").Op(":=").Qual("encoding/json", "Marshal").Call(
				Id("response").Dot("Data").Index(Lit(0)),
			),
			If(Id("err").Op("!=").Nil().Block(
				Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Id("err"),
			)),
			Line(),
			If(
				Id("err").Op("=").Qual("encoding/json", "Unmarshal").Call(
					Id("jsonData").Op(",").Op("&").Id(strcase.ToLowerCamel(mc.TypeName)),
				).Op(";").Id("err").Op("!=").Nil().Block(
					Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Id("err"),
				),
			),
			Line(),
			Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Nil(),
		)
		f.Line()
		// update object
		updateFuncName := fmt.Sprintf("Update%s", mc.TypeName)
		f.Comment(fmt.Sprintf(
			"%s creates a new %s",
			updateFuncName,
			strcase.ToDelimited(mc.TypeName, ' '),
		))
		f.Func().Id(updateFuncName).Params(
			Id("id").Uint().Op(",").Id(fmt.Sprintf("json%s", mc.TypeName)).Index().Byte().
				Op(",").Id("apiAddr").Op(",").Id("apiToken").String(),
		).Parens(List(
			Op("*").Qual(
				"github.com/threeport/threeport/pkg/api/v0",
				mc.TypeName,
			),
			Error(),
		)).Block(
			Var().Id(strcase.ToLowerCamel(mc.TypeName)).Qual(
				"github.com/threeport/threeport/pkg/api/v0",
				mc.TypeName,
			),
			Line(),
			Id("response").Op(",").Id("err").Op(":=").Id("GetResponse").Call(
				Line().Qual("fmt", "Sprintf").Call(
					Lit(fmt.Sprintf(
						"%%s/%%s/%s/%%d", pluralize.Pluralize(strcase.ToSnake(mc.TypeName), 2, false),
					)).Op(",").
						Id("apiAddr").Op(",").Id("ApiVersion").Op(",").Id("id"),
				),
				Line().Id("apiToken"),
				Line().Qual("net/http", "MethodPatch"),
				Line().Qual("bytes", "NewBuffer").Call(Id(
					fmt.Sprintf("json%s", mc.TypeName),
				)),
				Line().Qual("net/http", "StatusOK"),
				Line(),
			),
			If(Id("err").Op("!=").Nil().Block(
				Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Id("err"),
			)),
			Line(),
			Id("jsonData").Op(",").Id("err").Op(":=").Qual("encoding/json", "Marshal").Call(
				Id("response").Dot("Data").Index(Lit(0)),
			),
			If(Id("err").Op("!=").Nil().Block(
				Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Id("err"),
			)),
			Line(),
			If(
				Id("err").Op("=").Qual("encoding/json", "Unmarshal").Call(
					Id("jsonData").Op(",").Op("&").Id(strcase.ToLowerCamel(mc.TypeName)),
				).Op(";").Id("err").Op("!=").Nil().Block(
					Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Id("err"),
				),
			),
			Line(),
			Return().Op("&").Id(strcase.ToLowerCamel(mc.TypeName)).Op(",").Nil(),
		)
		f.Line()
	}

	// write code to file
	genFilename := fmt.Sprintf("%s_gen.go", name.FilenameSansExt(cc.ModelFilename))
	genFilepath := filepath.Join(clientLibPath(cc.PackageName), genFilename)
	file, err := os.OpenFile(genFilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed open file to write generated code for model client library: %w", err)
	}
	defer file.Close()
	if err := f.Render(file); err != nil {
		return fmt.Errorf("failed to render generated source code for model client library: %w", err)
	}
	fmt.Printf("code generation complete for %s client library\n", cc.ControllerDomainLower)

	return nil
}