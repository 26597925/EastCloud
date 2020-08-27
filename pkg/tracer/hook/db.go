package hook

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

type DbHook struct{
	Ctx context.Context
}

func (dbHook DbHook) AfterFind(scope *gorm.Scope) {
	if !trace.SpanFromContext(dbHook.Ctx).IsRecording() {
		return
	}

	tracer := global.Tracer("github.com/jinzhu/gorm")
	_, span := tracer.Start(dbHook.Ctx, "Db::query" )
	span.SetAttributes(
		label.String("db.sql", scope.SQL),
		label.String("db.var", fmt.Sprintf("%v", scope.SQLVars)),
	)
	span.End()
}

func (dbHook DbHook) AfterCreate(scope *gorm.Scope) {
	if !trace.SpanFromContext(dbHook.Ctx).IsRecording() {
		return
	}

	tracer := global.Tracer("github.com/jinzhu/gorm")
	_, span := tracer.Start(dbHook.Ctx, "Db::create" )
	span.SetAttributes(
		label.String("db.sql", scope.SQL),
		label.String("db.var", fmt.Sprintf("%v", scope.SQLVars)),
	)
	span.End()
}

func (dbHook DbHook) AfterUpdate(scope *gorm.Scope) {
	if !trace.SpanFromContext(dbHook.Ctx).IsRecording() {
		return
	}

	tracer := global.Tracer("github.com/jinzhu/gorm")
	_, span := tracer.Start(dbHook.Ctx, "Db::update" )
	span.SetAttributes(
		label.String("db.sql", scope.SQL),
		label.String("db.var", fmt.Sprintf("%v", scope.SQLVars)),
	)
	span.End()
}

func (dbHook DbHook) AfterDelete(scope *gorm.Scope) {
	if !trace.SpanFromContext(dbHook.Ctx).IsRecording() {
		return
	}

	tracer := global.Tracer("github.com/jinzhu/gorm")
	_, span := tracer.Start(dbHook.Ctx, "Db::delete" )
	span.SetAttributes(
		label.String("db.sql", scope.SQL),
		label.String("db.var", fmt.Sprintf("%v", scope.SQLVars)),
	)
	span.End()
}
