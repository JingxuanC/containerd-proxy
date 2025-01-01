package executor

import (
	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/runtimeproxy/store"
	"github.com/koordinator-sh/koordinator/pkg/runtimeproxy/utils"
)

type Executor interface {
	GetMetaInfo() string
	GenerateHookRequest() interface{}
	// ParseRequest would be the first function after request intercepted, during which,
	// pod/container's meta/resource info would be parsed from request or loaded from local store,
	// and some hint info should also be offered during ParseRequest stage, e.g. to check if executor
	// should call hook plugins when pod/container is system component.
	ParseRequest(request interface{}) (string, error)
	ResourceCheckPoint(response interface{}) error
	DeleteCheckpointIfNeed(request interface{}) error
	DeleteCheckpointForce(id string)
	UpdateRequest(response interface{}, request interface{}) error
}

// NoopResourceExecutor means no-operation for cri request,
// where no hook exists like ListContainerStats/ExecSync etc.
type NoopResourceExecutor struct {
}

func (n *NoopResourceExecutor) GetMetaInfo() string {
	return ""
}

func (n *NoopResourceExecutor) GenerateResourceCheckpoint() interface{} {
	return v1alpha1.ContainerResourceHookRequest{}
}

func (n *NoopResourceExecutor) GenerateHookRequest() interface{} {
	return store.ContainerInfo{}
}

func (n *NoopResourceExecutor) ParseRequest(request interface{}) (utils.CallHookPluginOperation, error) {
	return utils.Unknown, nil
}
func (n *NoopResourceExecutor) ResourceCheckPoint(response interface{}) error {
	return nil
}

func (n *NoopResourceExecutor) DeleteCheckpointIfNeed(request interface{}) error {
	return nil
}

func (n *NoopResourceExecutor) UpdateRequest(response interface{}, request interface{}) error {
	return nil
}
