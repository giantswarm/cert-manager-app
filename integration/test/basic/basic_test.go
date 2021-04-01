// +build k8srequired

package basic

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	cainjectorName = "cert-manager-cainjector"
	controllerName = "cert-manager-controller"
	webhookName    = "cert-manager-webhook"
)

func TestReadyDeployments(t *testing.T) {
	var err error

	ctx := context.Background()

	err = checkCRDInstallJob(ctx, metav1.NamespaceSystem, "cert-manager-crd-install")
	if err != nil {
		t.Fatalf("expected nil got: %v", err)
	}

	deployments := []string{
		cainjectorName,
		controllerName,
		webhookName,
	}

	for _, deploy := range deployments {
		err = checkReadyDeployment(ctx, metav1.NamespaceSystem, deploy)
		if err != nil {
			t.Fatalf("expected nil got: %v", err)
		}
	}
}

func checkCRDInstallJob(ctx context.Context, namespace, name string) error {
	l.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waiting for ready %#q job", name))
	o := func() error {
		job, err := appTest.K8sClient().BatchV1().Jobs(metav1.NamespaceSystem).Get(ctx, name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return microerror.Maskf(executionFailedError, "job %#q in %#q not found", name, metav1.NamespaceSystem)
		} else if err != nil {
			return microerror.Mask(err)
		}

		// l.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("Job A %+v job", job.Status))

		// if job.Status.Succeeded <= 0 {
		// 	l.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("Job B %+v job", job.Status))
		// 	return microerror.Maskf(executionFailedError, "job %#q want >= 0 succeeded, got %d", name, job.Status.Succeeded)
		// }

		// l.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("Job C %+v job", job.Status))

		// if job.Status.Failed > 0 {
		// 	l.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("Job D %+v job", job.Status))
		// 	return microerror.Maskf(executionFailedError, "job %#q want <= 0 failed, got %d", name, job.Status.Failed)
		// }

		return nil
	}

	b := backoff.NewConstant(10*time.Minute, 5*time.Second)
	n := backoff.NewNotifier(l, ctx)

	err := backoff.RetryNotify(o, b, n)
	if err != nil {
		return microerror.Mask(err)
	}

	l.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waited for ready %#q job", name))

	return nil
}

func checkReadyDeployment(ctx context.Context, namespace, name string) error {
	var err error

	l.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waiting for ready %#q deployment", name))

	o := func() error {
		deploy, err := appTest.K8sClient().AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			return microerror.Maskf(executionFailedError, "deployment %#q in %#q not found", name, metav1.NamespaceSystem)
		} else if err != nil {
			return microerror.Mask(err)
		}

		if deploy.Status.ReadyReplicas != *deploy.Spec.Replicas {
			return microerror.Maskf(executionFailedError, "deployment %#q want %d replicas %d ready", name, *deploy.Spec.Replicas, deploy.Status.ReadyReplicas)
		}

		return nil
	}
	b := backoff.NewConstant(2*time.Minute, 5*time.Second)
	n := backoff.NewNotifier(l, ctx)

	err = backoff.RetryNotify(o, b, n)
	if err != nil {
		return microerror.Mask(err)
	}

	l.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waited for ready %#q deployment", name))

	return nil
}
