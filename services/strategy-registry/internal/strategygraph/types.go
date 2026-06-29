package strategygraph

import (
	attribute "github.com/stratflow-labs/stratflow/services/strategy-registry/internal/attribute"

	"github.com/google/uuid"
)

const listAllPageSize = 100

type EntityRef struct {
	ID   *uuid.UUID
	Slug *string
}

type RelationTargetInput struct {
	AttributeRef EntityRef
	ValueRef     EntityRef
}

type Action struct {
	CreateAttribute  *CreateInput
	UpdateAttribute  *UpdateInput
	DeleteAttribute  *DeleteInput
	CreateValue      *CreateValueInput
	UpdateValue      *UpdateValueInput
	DeleteValue      *DeleteValueInput
	ReplaceRelations *ReplaceRelationsInput
}

type CreateInput struct {
	Slug        string
	Name        string
	Description string
}

type UpdateInput struct {
	AttributeRef EntityRef
	Slug         *string
	Name         *string
	Description  *string
}

type DeleteInput struct {
	AttributeRef EntityRef
}

type CreateValueInput struct {
	AttributeRef EntityRef
	Slug         string
	Value        string
}

type UpdateValueInput struct {
	AttributeRef EntityRef
	ValueRef     EntityRef
	Slug         *string
	Value        *string
}

type DeleteValueInput struct {
	AttributeRef EntityRef
	ValueRef     EntityRef
}

type ReplaceRelationsInput struct {
	AttributeRef EntityRef
	ValueRef     EntityRef
	Relations    []RelationTargetInput
}

type BatchActionInput struct {
	StrategyID uuid.UUID
	Actions    []Action
}

type BatchActionOutput struct {
	Strategy   StrategyView
	Attributes []attribute.AttributeView
}

type StrategyGraphEntityRef = EntityRef
type StrategyGraphRelationTargetInput = RelationTargetInput
type StrategyGraphAction = Action

type BatchCreateInput = CreateInput
type BatchUpdateInput = UpdateInput
type BatchDeleteInput = DeleteInput
type BatchCreateValueInput = CreateValueInput
type BatchUpdateValueInput = UpdateValueInput
type BatchDeleteValueInput = DeleteValueInput
type BatchReplaceRelationsInput = ReplaceRelationsInput

type BatchActionStrategyGraphInput = BatchActionInput
type BatchActionStrategyGraphOutput = BatchActionOutput
