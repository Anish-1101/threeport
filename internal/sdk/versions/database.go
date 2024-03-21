package versions

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/dave/jennifer/jen"
)

// DatabaseInit generates code for database initialization.
func (gvc *GlobalVersionConfig) DatabaseInit() error {
	f := NewFile("database")
	f.HeaderComment("generated by 'threeport-sdk gen api-version' - do not edit")
	f.ImportAlias("github.com/threeport/threeport/pkg/log/v0", "log")

	autoMigrateCalls := &Statement{}
	autoMigrateCalls.Line()
	for _, version := range gvc.Versions {
		for _, name := range version.DatabaseInitNames {
			autoMigrateCalls.List(
				Op("&").Qual(
					fmt.Sprintf(
						"github.com/threeport/threeport/pkg/api/%s", version.VersionName,
					),
					name,
				).Values().Op(","),
			)
			autoMigrateCalls.Line()
		}
	}

	f.Comment("ZapLogger is a custom GORM logger that forwards log messages to a Zap logger.")
	f.Type().Id("ZapLogger").Struct(
		Id("Logger").Op("*").Qual(
			"go.uber.org/zap", "Logger",
		),
	)
	f.Line()

	f.Comment("Init initializes the API database.")
	f.Func().Id("Init").Params(
		Id("autoMigrate").Bool().Op(",").Id("logger").Op("*").Qual(
			"go.uber.org/zap", "Logger",
		),
	).Parens(Op("*").Qual(
		"gorm.io/gorm",
		"DB",
	).Op(",").Id("error")).Block(
		Id("dsn").Op(":=").Qual(
			"fmt", "Sprintf",
		).Call(
			Line().Lit("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC").Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_HOST")).Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_USER")).Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_PASSWORD")).Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_NAME")).Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_PORT")).Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_SSL_MODE")).Op(",").Line(),
		),
		Line(),
		Id("db").Op(",").Id("err").Op(":=").Qual(
			"gorm.io/gorm",
			"Open",
		).Call(Qual(
			"gorm.io/driver/postgres",
			"Open",
		).Call(Id("dsn")).Op(",").Op("&").Qual(
			"gorm.io/gorm",
			"Config",
		).Values(Dict{
			Id("Logger"): Op("&").Id("ZapLogger").Values(Dict{
				Id("Logger"): Id("logger").Op(","),
			}),
			//Id("Logger"): Qual(
			//	"gorm.io/gorm/logger",
			//	"Default",
			//).Dot("LogMode").Call(Qual(
			//	"gorm.io/gorm/logger",
			//	"Info",
			//).Op(",")),
			Id("NowFunc"): Func().Call().Qual(
				"time", "Time",
			).Block(
				Id("utc").Op(",").Id("_").Op(":=").Qual(
					"time", "LoadLocation",
				).Call(Lit("UTC")),
				Return().Qual(
					"time", "Now",
				).Call().Dot("In").Call(Id("utc")).Dot("Truncate").Call(Qual(
					"time", "Microsecond",
				)),
			),
		})),
		If(
			Id("err").Op("!=").Nil().Block(
				Return().Nil().Op(",").Id("err"),
			),
		),
		Line(),
		Return().Id("db").Op(",").Nil(),
	)

	f.Comment("LogMode overrides the standard GORM logger's LogMode method to set the logger mode.")
	f.Func().Parens(
		Id("zl").Op("*").Id("ZapLogger"),
	).Id("LogMode").Params(
		Id("level").Qual("gorm.io/gorm/logger", "LogLevel"),
	).Qual("gorm.io/gorm/logger", "Interface").Block(
		Return().Id("zl"),
	)
	f.Line()

	f.Comment("Info overrides the standard GORM logger's Info method to forward log messages")
	f.Comment("to the zap logger.")
	f.Func().Parens(
		Id("zl").Op("*").Id("ZapLogger"),
	).Id("Info").Params(
		Id("ctx").Qual("context", "Context").Op(",").Id("msg").String().Op(",").Id("data").Op("...").Interface(),
	).Block(
		Id("fields").Op(":=").Make(
			Index().Qual(
				"go.uber.org/zap", "Field",
			).Op(",").Lit(0).Op(",").Len(Id("data")),
		),
		For(Id("i").Op(":=").Lit(0).Op(";").Id("i").Op("<").Len(Id("data")).Op(";").Id("i").Op("+=").Lit(2)).Block(
			Id("fields").Op("=").Append(Id("fields").Op(",").Qual(
				"go.uber.org/zap", "Any",
			).Call(Id("data").Index(Id("i")).Assert(String()).Op(",").Id("data").Index(Id("i").Op("+").Lit(1)))),
		),
		Id("zl").Dot("Logger").Dot("Info").Call(Id("msg").Op(",").Id("fields").Op("...")),
	)
	f.Line()

	f.Comment("Warn overrides the standard GORM logger's Warn method to forward log messages")
	f.Comment("to the zap logger.")
	f.Func().Parens(
		Id("zl").Op("*").Id("ZapLogger"),
	).Id("Warn").Params(
		Id("ctx").Qual("context", "Context").Op(",").Id("msg").String().Op(",").Id("data").Op("...").Interface(),
	).Block(
		Id("fields").Op(":=").Make(
			Index().Qual(
				"go.uber.org/zap", "Field",
			).Op(",").Lit(0).Op(",").Len(Id("data")),
		),
		For(Id("i").Op(":=").Lit(0).Op(";").Id("i").Op("<").Len(Id("data")).Op(";").Id("i").Op("+=").Lit(2)).Block(
			Id("fields").Op("=").Append(Id("fields").Op(",").Qual(
				"go.uber.org/zap", "Any",
			).Call(Id("data").Index(Id("i")).Assert(String()).Op(",").Id("data").Index(Id("i").Op("+").Lit(1)))),
		),
		Id("zl").Dot("Logger").Dot("Warn").Call(Id("msg").Op(",").Id("fields").Op("...")),
	)
	f.Line()

	f.Comment("Error overrides the standard GORM logger's Error method to forward log messages")
	f.Comment("to the zap logger.")
	f.Func().Parens(
		Id("zl").Op("*").Id("ZapLogger"),
	).Id("Error").Params(
		Id("ctx").Qual("context", "Context").Op(",").Id("msg").String().Op(",").Id("data").Op("...").Interface(),
	).Block(
		Id("fields").Op(":=").Make(
			Index().Qual(
				"go.uber.org/zap", "Field",
			).Op(",").Lit(0).Op(",").Len(Id("data")),
		),
		For(Id("i").Op(":=").Lit(0).Op(";").Id("i").Op("<").Len(Id("data")).Op(";").Id("i").Op("+=").Lit(2)).Block(
			If(Qual(
				"reflect", "TypeOf",
			).Call(Id("data").Index(Id("i"))).Dot("Kind").Call().Op("==").Qual(
				"reflect", "Ptr",
			).Block(
				Id("data").Index(Id("i")).Op("=").Qual(
					"fmt", "Sprintf",
				).Call(Lit("%+v").Op(",").Id("data").Index(Id("i"))),
			)),
			Id("fields").Op("=").Append(Id("fields").Op(",").Qual(
				"go.uber.org/zap", "Any",
			).Call(Id("data").Index(Id("i")).Assert(String()).Op(",").Id("data").Index(Id("i").Op("+").Lit(1)))),
		),
		Id("zl").Dot("Logger").Dot("Error").Call(Id("msg").Op(",").Id("fields").Op("...")),
	)
	f.Line()

	f.Comment("Trace overrides the standard GORM logger's Trace method to forward log messages")
	f.Comment("to the zap logger.")
	f.Func().Parens(
		Id("zl").Op("*").Id("ZapLogger"),
	).Id("Trace").Params(
		Id("ctx").Qual(
			"context", "Context",
		).Op(",").Id("begin").Qual(
			"time", "Time",
		).Op(",").Id("fc").Func().Call().Parens(String().Op(",").Int64()).Op(",").Id("err").Error(),
	).Block(
		Comment("use the fc function to get the SQL statement and execution time"),
		Id("sql").Op(",").Id("rows").Op(":=").Id("fc").Call(),
		Line(),
		Comment("create a new logger with some additional fields"),
		Id("logger").Op(":=").Id("zl").Dot("Logger").Dot("With").Call(
			Line().Qual("go.uber.org/zap", "String").Call(Lit("type").Op(",").Lit("sql")),
			Line().Qual("go.uber.org/zap", "String").Call(Lit("sql").Op(",").Id("suppressSensitive").Call(Id("sql"))),
			Line().Qual("go.uber.org/zap", "Int64").Call(Lit("rows").Op(",").Id("rows")),
			Line().Qual("go.uber.org/zap", "Duration").Call(Lit("elapsed").Op(",").Qual(
				"time", "Since",
			).Call(Id("begin"))),
			Line(),
		),
		Line(),
		Comment("if an error occurred, add it as a field to the logger"),
		If(Id("err").Op("!=").Nil().Block(
			Id("logger").Op("=").Qual(
				"gorm.io/gorm/logger", "With",
			).Call(Qual(
				"go.uber.org/zap", "Error",
			).Call(Id("err"))),
		)),
		Line(),
		Comment("log the message using the logger"),
		Qual("gorm.io/gorm/logger", "Debug").Call(Lit("gorm query")),
	)
	f.Line()

	f.Comment("Return all database init object interfaces.")
	f.Func().Id("GetDbInterfaces").Parens(Empty()).Params(
		Index().Interface(),
	).Block(
		Return().Index().Interface().Block(autoMigrateCalls),
	)
	f.Line()

	f.Comment("suppressSensitive supresses messages containing sesitive strings.")
	f.Func().Id("suppressSensitive").Params(
		Id("msg").String(),
	).String().Block(
		For(Id("_").Op(",").Id("str").Op(":=").Range().Qual(
			"github.com/threeport/threeport/pkg/log/v0", "SensitiveStrings",
		).Call()).Block(
			If(Qual("strings", "Contains").Call(Id("msg").Op(",").Id("str"))).Block(
				Return().Qual("fmt", "Sprintf").Call(Lit("[log message containing %s supporessed]").Op(",").Id("str")),
			),
		),
		Line(),
		Return().Id("msg"),
	)

	// write code to file
	databaseInitFilepath := filepath.Join("pkg", "api-server", "v0", "database", "database_gen.go")
	file, err := os.OpenFile(databaseInitFilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file to write generated code for database initializer: %w", err)
	}
	defer file.Close()
	if err := f.Render(file); err != nil {
		return fmt.Errorf("failed to render generated source code for database initializer: %w", err)
	}
	fmt.Println("code generation complete for database initializer")

	return nil
}

// ExtensionDatabaseInit generates code for database initialization.
func (gvc *GlobalVersionConfig) ExtensionDatabaseInit(modulePath string) error {
	f := NewFile("database")
	f.HeaderComment("generated by 'threeport-sdk codegen api-version' - do not edit")
	f.ImportAlias("github.com/threeport/threeport/pkg/log/v0", "log")

	autoMigrateCalls := &Statement{}
	autoMigrateCalls.Line()
	for _, version := range gvc.Versions {
		for _, name := range version.DatabaseInitNames {
			autoMigrateCalls.List(
				Op("&").Qual(
					fmt.Sprintf(
						"%s/pkg/api/%s", modulePath, version.VersionName,
					),
					name,
				).Values().Op(","),
			)
			autoMigrateCalls.Line()
		}
	}

	f.Comment("ZapLogger is a custom GORM logger that forwards log messages to a Zap logger.")
	f.Type().Id("ZapLogger").Struct(
		Id("Logger").Op("*").Qual(
			"go.uber.org/zap", "Logger",
		),
	)
	f.Line()

	f.Comment("Init initializes the API database.")
	f.Func().Id("Init").Params(
		Id("autoMigrate").Bool().Op(",").Id("logger").Op("*").Qual(
			"go.uber.org/zap", "Logger",
		),
	).Parens(Op("*").Qual(
		"gorm.io/gorm",
		"DB",
	).Op(",").Id("error")).Block(
		Id("dsn").Op(":=").Qual(
			"fmt", "Sprintf",
		).Call(
			Line().Lit("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC").Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_HOST")).Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_USER")).Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_PASSWORD")).Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_NAME")).Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_PORT")).Op(",").Line().
				Qual(
					"os", "Getenv",
				).Call(Lit("DB_SSL_MODE")).Op(",").Line(),
		),
		Line(),
		Id("db").Op(",").Id("err").Op(":=").Qual(
			"gorm.io/gorm",
			"Open",
		).Call(Qual(
			"gorm.io/driver/postgres",
			"Open",
		).Call(Id("dsn")).Op(",").Op("&").Qual(
			"gorm.io/gorm",
			"Config",
		).Values(Dict{
			Id("Logger"): Op("&").Id("ZapLogger").Values(Dict{
				Id("Logger"): Id("logger").Op(","),
			}),
			//Id("Logger"): Qual(
			//	"gorm.io/gorm/logger",
			//	"Default",
			//).Dot("LogMode").Call(Qual(
			//	"gorm.io/gorm/logger",
			//	"Info",
			//).Op(",")),
			Id("NowFunc"): Func().Call().Qual(
				"time", "Time",
			).Block(
				Id("utc").Op(",").Id("_").Op(":=").Qual(
					"time", "LoadLocation",
				).Call(Lit("UTC")),
				Return().Qual(
					"time", "Now",
				).Call().Dot("In").Call(Id("utc")).Dot("Truncate").Call(Qual(
					"time", "Microsecond",
				)),
			),
		})),
		If(
			Id("err").Op("!=").Nil().Block(
				Return().Nil().Op(",").Id("err"),
			),
		),
		Line(),

		Id("dbAutoMigrate").Op(":=").Index().Interface().Block(autoMigrateCalls),
		Id("threeportDbAutoMigrate").Op(":=").Qual(
			"github.com/threeport/threeport/pkg/api-server/v0/database",
			"GetDbInterfaces",
		).Parens(Empty()),
		For(Id("_").Op(",").Id("obj")).Op(":=").Range().Id("dbAutoMigrate").Block(
			Id("threeportDbAutoMigrate").Op("=").Append(Id("threeportDbAutoMigrate").Op(",").Id("obj")),
		),
		If(
			Id("autoMigrate").Block(
				If(Id("err").Op(":=").Id("db").Dot("AutoMigrate").Call(
					Id("threeportDbAutoMigrate").Op("..."),
				).Op(";").Id("err").Op("!=").Nil().Block(
					Return().Nil().Op(",").Id("err")),
				),
			),
		),
		Line(),
		Return().Id("db").Op(",").Nil(),
	)

	f.Comment("LogMode overrides the standard GORM logger's LogMode method to set the logger mode.")
	f.Func().Parens(
		Id("zl").Op("*").Id("ZapLogger"),
	).Id("LogMode").Params(
		Id("level").Qual("gorm.io/gorm/logger", "LogLevel"),
	).Qual("gorm.io/gorm/logger", "Interface").Block(
		Return().Id("zl"),
	)
	f.Line()

	f.Comment("Info overrides the standard GORM logger's Info method to forward log messages")
	f.Comment("to the zap logger.")
	f.Func().Parens(
		Id("zl").Op("*").Id("ZapLogger"),
	).Id("Info").Params(
		Id("ctx").Qual("context", "Context").Op(",").Id("msg").String().Op(",").Id("data").Op("...").Interface(),
	).Block(
		Id("fields").Op(":=").Make(
			Index().Qual(
				"go.uber.org/zap", "Field",
			).Op(",").Lit(0).Op(",").Len(Id("data")),
		),
		For(Id("i").Op(":=").Lit(0).Op(";").Id("i").Op("<").Len(Id("data")).Op(";").Id("i").Op("+=").Lit(2)).Block(
			Id("fields").Op("=").Append(Id("fields").Op(",").Qual(
				"go.uber.org/zap", "Any",
			).Call(Id("data").Index(Id("i")).Assert(String()).Op(",").Id("data").Index(Id("i").Op("+").Lit(1)))),
		),
		Id("zl").Dot("Logger").Dot("Info").Call(Id("msg").Op(",").Id("fields").Op("...")),
	)
	f.Line()

	f.Comment("Warn overrides the standard GORM logger's Warn method to forward log messages")
	f.Comment("to the zap logger.")
	f.Func().Parens(
		Id("zl").Op("*").Id("ZapLogger"),
	).Id("Warn").Params(
		Id("ctx").Qual("context", "Context").Op(",").Id("msg").String().Op(",").Id("data").Op("...").Interface(),
	).Block(
		Id("fields").Op(":=").Make(
			Index().Qual(
				"go.uber.org/zap", "Field",
			).Op(",").Lit(0).Op(",").Len(Id("data")),
		),
		For(Id("i").Op(":=").Lit(0).Op(";").Id("i").Op("<").Len(Id("data")).Op(";").Id("i").Op("+=").Lit(2)).Block(
			Id("fields").Op("=").Append(Id("fields").Op(",").Qual(
				"go.uber.org/zap", "Any",
			).Call(Id("data").Index(Id("i")).Assert(String()).Op(",").Id("data").Index(Id("i").Op("+").Lit(1)))),
		),
		Id("zl").Dot("Logger").Dot("Warn").Call(Id("msg").Op(",").Id("fields").Op("...")),
	)
	f.Line()

	f.Comment("Error overrides the standard GORM logger's Error method to forward log messages")
	f.Comment("to the zap logger.")
	f.Func().Parens(
		Id("zl").Op("*").Id("ZapLogger"),
	).Id("Error").Params(
		Id("ctx").Qual("context", "Context").Op(",").Id("msg").String().Op(",").Id("data").Op("...").Interface(),
	).Block(
		Id("fields").Op(":=").Make(
			Index().Qual(
				"go.uber.org/zap", "Field",
			).Op(",").Lit(0).Op(",").Len(Id("data")),
		),
		For(Id("i").Op(":=").Lit(0).Op(";").Id("i").Op("<").Len(Id("data")).Op(";").Id("i").Op("+=").Lit(2)).Block(
			If(Qual(
				"reflect", "TypeOf",
			).Call(Id("data").Index(Id("i"))).Dot("Kind").Call().Op("==").Qual(
				"reflect", "Ptr",
			).Block(
				Id("data").Index(Id("i")).Op("=").Qual(
					"fmt", "Sprintf",
				).Call(Lit("%+v").Op(",").Id("data").Index(Id("i"))),
			)),
			Id("fields").Op("=").Append(Id("fields").Op(",").Qual(
				"go.uber.org/zap", "Any",
			).Call(Id("data").Index(Id("i")).Assert(String()).Op(",").Id("data").Index(Id("i").Op("+").Lit(1)))),
		),
		Id("zl").Dot("Logger").Dot("Error").Call(Id("msg").Op(",").Id("fields").Op("...")),
	)
	f.Line()

	f.Comment("Trace overrides the standard GORM logger's Trace method to forward log messages")
	f.Comment("to the zap logger.")
	f.Func().Parens(
		Id("zl").Op("*").Id("ZapLogger"),
	).Id("Trace").Params(
		Id("ctx").Qual(
			"context", "Context",
		).Op(",").Id("begin").Qual(
			"time", "Time",
		).Op(",").Id("fc").Func().Call().Parens(String().Op(",").Int64()).Op(",").Id("err").Error(),
	).Block(
		Comment("use the fc function to get the SQL statement and execution time"),
		Id("sql").Op(",").Id("rows").Op(":=").Id("fc").Call(),
		Line(),
		Comment("create a new logger with some additional fields"),
		Id("logger").Op(":=").Id("zl").Dot("Logger").Dot("With").Call(
			Line().Qual("go.uber.org/zap", "String").Call(Lit("type").Op(",").Lit("sql")),
			Line().Qual("go.uber.org/zap", "String").Call(Lit("sql").Op(",").Id("suppressSensitive").Call(Id("sql"))),
			Line().Qual("go.uber.org/zap", "Int64").Call(Lit("rows").Op(",").Id("rows")),
			Line().Qual("go.uber.org/zap", "Duration").Call(Lit("elapsed").Op(",").Qual(
				"time", "Since",
			).Call(Id("begin"))),
			Line(),
		),
		Line(),
		Comment("if an error occurred, add it as a field to the logger"),
		If(Id("err").Op("!=").Nil().Block(
			Id("logger").Op("=").Qual(
				"gorm.io/gorm/logger", "With",
			).Call(Qual(
				"go.uber.org/zap", "Error",
			).Call(Id("err"))),
		)),
		Line(),
		Comment("log the message using the logger"),
		Qual("gorm.io/gorm/logger", "Debug").Call(Lit("gorm query")),
	)
	f.Line()

	f.Comment("suppressSensitive supresses messages containing sesitive strings.")
	f.Func().Id("suppressSensitive").Params(
		Id("msg").String(),
	).String().Block(
		For(Id("_").Op(",").Id("str").Op(":=").Range().Qual(
			"github.com/threeport/threeport/pkg/log/v0", "SensitiveStrings",
		).Call()).Block(
			If(Qual("strings", "Contains").Call(Id("msg").Op(",").Id("str"))).Block(
				Return().Qual("fmt", "Sprintf").Call(Lit("[log message containing %s supporessed]").Op(",").Id("str")),
			),
		),
		Line(),
		Return().Id("msg"),
	)

	// write code to file
	databaseInitFilepath := filepath.Join("..", "..", "pkg", "api-server", "v0", "database", "database_gen.go")
	file, err := os.OpenFile(databaseInitFilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file to write generated code for database initializer: %w", err)
	}
	defer file.Close()
	if err := f.Render(file); err != nil {
		return fmt.Errorf("failed to render generated source code for database initializer: %w", err)
	}
	fmt.Println("code generation complete for database initializer")

	return nil
}
