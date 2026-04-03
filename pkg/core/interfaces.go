package core

import (
	"context"

	"gateway/pkg/models"
)

// IngestPlugin 定义数据接入插件接口
type IngestPlugin interface {
	Start(ctx context.Context, ch chan<- models.NormalizedData) error
	Stop() error
}

// EgressPlugin 定义数据输出插件接口
type EgressPlugin interface {
	Start(ctx context.Context) error
	Send(data models.NormalizedData) error
	Stop() error
}

// TransformProfile 定义数据转换规则接口
type TransformProfile interface {
	Transform(raw interface{}) (models.NormalizedData, error)
}
