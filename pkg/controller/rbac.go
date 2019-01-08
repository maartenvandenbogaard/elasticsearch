package controller

import (
	core_util "github.com/appscode/kutil/core/v1"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/reference"
)

func (c *Controller) deleteServiceAccount(elasticsearch *api.Elasticsearch) error {
	// Delete existing ServiceAccount
	if err := c.Client.CoreV1().ServiceAccounts(elasticsearch.Namespace).Delete(elasticsearch.OffshootName(), nil); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (c *Controller) createServiceAccount(elasticsearch *api.Elasticsearch) error {
	ref, rerr := reference.GetReference(clientsetscheme.Scheme, elasticsearch)
	if rerr != nil {
		return rerr
	}
	// Create new ServiceAccount
	_, _, err := core_util.CreateOrPatchServiceAccount(
		c.Client,
		metav1.ObjectMeta{
			Name:      elasticsearch.OffshootName(),
			Namespace: elasticsearch.Namespace,
		},
		func(in *core.ServiceAccount) *core.ServiceAccount {
			core_util.EnsureOwnerReference(&in.ObjectMeta, ref)
			return in
		},
	)
	return err
}

func (c *Controller) deleteRoleBinding(elasticsearch *api.Elasticsearch) error {
	// Delete existing RoleBindings
	if err := c.Client.RbacV1beta1().RoleBindings(elasticsearch.Namespace).Delete(elasticsearch.OffshootName(), nil); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (c *Controller) ensureRBACStuff(elasticsearch *api.Elasticsearch) error {
	// Create New ServiceAccount
	if err := c.createServiceAccount(elasticsearch); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

func (c *Controller) deleteRBACStuff(elasticsearch *api.Elasticsearch) error {
	// Delete ServiceAccount
	if err := c.deleteServiceAccount(elasticsearch); err != nil {
		return err
	}

	return nil
}
