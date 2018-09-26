package stub

import (
	"context"

	"github.com/kubevirt/kubevirt-operator/pkg/apis/virt/v1alpha1"
	"github.com/kubevirt/kubevirt-operator/pkg/kubevirt"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Virt:
		return kubevirt.Reconcile(o)
	}
	return nil
}
