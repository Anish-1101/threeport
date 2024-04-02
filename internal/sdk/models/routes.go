package models

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/dave/jennifer/jen"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"

	"github.com/threeport/threeport/internal/sdk"
)

// apiRoutesPath returns the path from the models to the API's internal routes
// package.
func apiRoutesPath(apiVersion string) string {
	return filepath.Join("pkg", "api-server", apiVersion, "routes")
}

// ModelRoutes generates the REST routes and maps them to their handlers.  These
// routes are assembled in route functions which are added to the echo server
// when the API starts.
func (cc *ControllerConfig) ModelRoutes() error {
	pluralize := pluralize.NewClient()
	f := NewFile("routes")
	f.HeaderComment("generated by 'threeport-sdk gen' for API routes boilerplate' - do not edit")
	f.ImportAlias("github.com/labstack/echo/v4", "echo")
	for i, mc := range cc.ModelConfigs {
		routeFuncName := fmt.Sprintf("%sRoutes", mc.TypeName)
		mc.GetVersionHandlerName = fmt.Sprintf("Get%sVersions", mc.TypeName)
		mc.AddHandlerName = fmt.Sprintf("Add%s", mc.TypeName)
		mc.AddMiddlewareFuncName = fmt.Sprintf("Add%sMiddleware", mc.TypeName)
		mc.GetAllHandlerName = fmt.Sprintf("Get%s", pluralize.Pluralize(mc.TypeName, 2, false))
		mc.GetOneHandlerName = fmt.Sprintf("Get%s", mc.TypeName)
		mc.GetMiddlewareFuncName = fmt.Sprintf("Get%sMiddleware", mc.TypeName)
		mc.PatchHandlerName = fmt.Sprintf("Update%s", mc.TypeName)
		mc.PatchMiddlewareFuncName = fmt.Sprintf("Patch%sMiddleware", mc.TypeName)
		mc.PutHandlerName = fmt.Sprintf("Replace%s", mc.TypeName)
		mc.PutMiddlewareFuncName = fmt.Sprintf("Put%sMiddleware", mc.TypeName)
		mc.DeleteHandlerName = fmt.Sprintf("Delete%s", mc.TypeName)
		mc.DeleteMiddlewareFuncName = fmt.Sprintf("Delete%sMiddleware", mc.TypeName)
		cc.ModelConfigs[i] = mc

		addMiddleware := Null()
		getMiddleware := Null()
		patchMiddleware := Null()
		putMiddleware := Null()
		deleteMiddleware := Null()

		if mc.AllowCustomMiddleware {
			addMiddleware = Id("h").Dot(mc.AddMiddlewareFuncName).Call().Op("...")
			getMiddleware = Id("h").Dot(mc.GetMiddlewareFuncName).Call().Op("...")
			patchMiddleware = Id("h").Dot(mc.PatchMiddlewareFuncName).Call().Op("...")
			putMiddleware = Id("h").Dot(mc.PutMiddlewareFuncName).Call().Op("...")
			deleteMiddleware = Id("h").Dot(mc.DeleteMiddlewareFuncName).Call().Op("...")
		}

		f.Comment(fmt.Sprintf(
			"%s sets up all routes for the %s handlers.", routeFuncName, mc.TypeName,
		))
		f.Func().Id(routeFuncName).Params(
			Id("e").Op("*").Qual(
				"github.com/labstack/echo/v4",
				"Echo",
			),
			Id("h").Op("*").Qual(
				fmt.Sprintf(
					"github.com/threeport/threeport/pkg/api-server/%s/handlers",
					cc.ApiVersion,
				),
				"Handler",
			),
		).Block(
			Id("e").Dot("GET").Call(
				Lit(
					fmt.Sprintf("/%s/versions", pluralize.Pluralize(strcase.ToKebab(mc.TypeName), 2, false)),
				),
				Id("h").Dot(mc.GetVersionHandlerName),
			),
			Line(),
			Id("e").Dot("POST").Call(
				Qual(
					fmt.Sprintf(
						"github.com/threeport/threeport/pkg/api/%s",
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				),
				Id("h").Dot(mc.AddHandlerName),
				addMiddleware,
			),
			Id("e").Dot("GET").Call(
				Qual(
					fmt.Sprintf(
						"github.com/threeport/threeport/pkg/api/%s",
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				),
				Id("h").Dot(mc.GetAllHandlerName),
				getMiddleware,
			),
			Id("e").Dot("GET").Call(
				Qual(
					fmt.Sprintf(
						"github.com/threeport/threeport/pkg/api/%s",
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				).Op("+").Lit("/:id"),
				Id("h").Dot(mc.GetOneHandlerName),
				getMiddleware,
			),
			Id("e").Dot("PATCH").Call(
				Qual(
					fmt.Sprintf(
						"github.com/threeport/threeport/pkg/api/%s",
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				).Op("+").Lit("/:id"),
				Id("h").Dot(mc.PatchHandlerName),
				patchMiddleware,
			),
			Id("e").Dot("PUT").Call(
				Qual(
					fmt.Sprintf(
						"github.com/threeport/threeport/pkg/api/%s",
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				).Op("+").Lit("/:id"),
				Id("h").Dot(mc.PutHandlerName),
				putMiddleware,
			),
			Id("e").Dot("DELETE").Call(
				Qual(
					fmt.Sprintf(
						"github.com/threeport/threeport/pkg/api/%s",
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				).Op("+").Lit("/:id"),
				Id("h").Dot(mc.DeleteHandlerName),
				deleteMiddleware,
			),
		)
	}

	// write code to file
	genFilename := fmt.Sprintf("%s_gen.go", sdk.FilenameSansExt(cc.ModelFilename))
	genFilepath := filepath.Join(apiRoutesPath(cc.ApiVersion), genFilename)
	file, err := os.OpenFile(genFilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file to write generated code for model routes: %w", err)
	}
	defer file.Close()
	if err := f.Render(file); err != nil {
		return fmt.Errorf("failed to render generated source code for model routes: %w", err)
	}
	fmt.Printf("code generation complete for %s model routes\n", cc.ControllerDomainLower)

	return nil
}

// ExtensionModelRoutes generates the REST routes and maps them to their handlers in an extension.  These
// routes are assembled in route functions which are added to the echo server
// when the API starts.
func (cc *ControllerConfig) ExtensionModelRoutes(modulePath string) error {
	pluralize := pluralize.NewClient()
	f := NewFile("routes")
	f.HeaderComment("generated by 'threeport-sdk gen' for API routes boilerplate' - do not edit")
	f.ImportAlias("github.com/labstack/echo/v4", "echo")
	for i, mc := range cc.ModelConfigs {
		routeFuncName := fmt.Sprintf("%sRoutes", mc.TypeName)
		mc.GetVersionHandlerName = fmt.Sprintf("Get%sVersions", mc.TypeName)
		mc.AddHandlerName = fmt.Sprintf("Add%s", mc.TypeName)
		mc.AddMiddlewareFuncName = fmt.Sprintf("Add%sMiddleware", mc.TypeName)
		mc.GetAllHandlerName = fmt.Sprintf("Get%s", pluralize.Pluralize(mc.TypeName, 2, false))
		mc.GetOneHandlerName = fmt.Sprintf("Get%s", mc.TypeName)
		mc.GetMiddlewareFuncName = fmt.Sprintf("Get%sMiddleware", mc.TypeName)
		mc.PatchHandlerName = fmt.Sprintf("Update%s", mc.TypeName)
		mc.PatchMiddlewareFuncName = fmt.Sprintf("Patch%sMiddleware", mc.TypeName)
		mc.PutHandlerName = fmt.Sprintf("Replace%s", mc.TypeName)
		mc.PutMiddlewareFuncName = fmt.Sprintf("Put%sMiddleware", mc.TypeName)
		mc.DeleteHandlerName = fmt.Sprintf("Delete%s", mc.TypeName)
		mc.DeleteMiddlewareFuncName = fmt.Sprintf("Delete%sMiddleware", mc.TypeName)
		cc.ModelConfigs[i] = mc

		addMiddleware := Null()
		getMiddleware := Null()
		patchMiddleware := Null()
		putMiddleware := Null()
		deleteMiddleware := Null()

		if mc.AllowCustomMiddleware {
			addMiddleware = Id("h").Dot(mc.AddMiddlewareFuncName).Call().Op("...")
			getMiddleware = Id("h").Dot(mc.GetMiddlewareFuncName).Call().Op("...")
			patchMiddleware = Id("h").Dot(mc.PatchMiddlewareFuncName).Call().Op("...")
			putMiddleware = Id("h").Dot(mc.PutMiddlewareFuncName).Call().Op("...")
			deleteMiddleware = Id("h").Dot(mc.DeleteMiddlewareFuncName).Call().Op("...")
		}

		f.Comment(fmt.Sprintf(
			"%s sets up all routes for the %s handlers.", routeFuncName, mc.TypeName,
		))
		f.Func().Id(routeFuncName).Params(
			Id("e").Op("*").Qual(
				"github.com/labstack/echo/v4",
				"Echo",
			),
			Id("h").Op("*").Qual(
				fmt.Sprintf("%s/pkg/api-server/v0/handlers", modulePath),
				"Handler",
			),
		).Block(
			Id("e").Dot("GET").Call(
				Lit(
					fmt.Sprintf("/%s/versions", pluralize.Pluralize(strcase.ToKebab(mc.TypeName), 2, false)),
				),
				Id("h").Dot(mc.GetVersionHandlerName),
			),
			Line(),
			Id("e").Dot("POST").Call(
				Qual(
					fmt.Sprintf(
						"%s/pkg/api/%s",
						modulePath,
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				),
				Id("h").Dot(mc.AddHandlerName),
				addMiddleware,
			),
			Id("e").Dot("GET").Call(
				Qual(
					fmt.Sprintf(
						"%s/pkg/api/%s",
						modulePath,
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				),
				Id("h").Dot(mc.GetAllHandlerName),
				getMiddleware,
			),
			Id("e").Dot("GET").Call(
				Qual(
					fmt.Sprintf(
						"%s/pkg/api/%s",
						modulePath,
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				).Op("+").Lit("/:id"),
				Id("h").Dot(mc.GetOneHandlerName),
				getMiddleware,
			),
			Id("e").Dot("PATCH").Call(
				Qual(
					fmt.Sprintf(
						"%s/pkg/api/%s",
						modulePath,
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				).Op("+").Lit("/:id"),
				Id("h").Dot(mc.PatchHandlerName),
				patchMiddleware,
			),
			Id("e").Dot("PUT").Call(
				Qual(
					fmt.Sprintf(
						"%s/pkg/api/%s",
						modulePath,
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				).Op("+").Lit("/:id"),
				Id("h").Dot(mc.PutHandlerName),
				putMiddleware,
			),
			Id("e").Dot("DELETE").Call(
				Qual(
					fmt.Sprintf(
						"%s/pkg/api/%s",
						modulePath,
						cc.Version,
					),
					fmt.Sprintf("Path%s", pluralize.Pluralize(mc.TypeName, 2, false)),
				).Op("+").Lit("/:id"),
				Id("h").Dot(mc.DeleteHandlerName),
				deleteMiddleware,
			),
		)
	}

	// write code to file
	genFilename := fmt.Sprintf("%s_gen.go", sdk.FilenameSansExt(cc.ModelFilename))
	genFilepath := filepath.Join(apiRoutesPath(cc.ApiVersion), genFilename)
	file, err := os.OpenFile(genFilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file to write generated code for model routes: %w", err)
	}
	defer file.Close()
	if err := f.Render(file); err != nil {
		return fmt.Errorf("failed to render generated source code for model routes: %w", err)
	}
	fmt.Printf("code generation complete for %s model routes\n", cc.ControllerDomainLower)

	return nil
}
