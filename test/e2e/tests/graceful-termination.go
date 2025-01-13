// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

//go:build e2e

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/gateway-api/conformance/utils/http"
	"sigs.k8s.io/gateway-api/conformance/utils/kubernetes"
	"sigs.k8s.io/gateway-api/conformance/utils/roundtripper"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
	"sigs.k8s.io/gateway-api/conformance/utils/tlog"
)

func init() {
	ConformanceTests = append(ConformanceTests, GracefulTerminationTest)
}

// CapturedServerState is the response body from echo-basic, including the lifecycle parameters of the pod/container
type CapturedServerState struct {
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`

	// Ready is true when readiness probes would succeed on this server
	Ready bool `json:"ready"`
	// Stopping is true after the stop hook was called on this server
	Stopping bool `json:"stopping"`
	// Terminating is true after the SIGTERM signal was sent to this server
	Terminating bool `json:"terminating"`
}

var GracefulTerminationTest = suite.ConformanceTest{
	ShortName:   "GracefulTermination",
	Description: "Graceful pod termination does not result in disruptions",
	Manifests:   []string{"testdata/graceful-termination.yaml", "testdata/graceful-termination-sts.yaml"},
	Test: func(t *testing.T, suite *suite.ConformanceTestSuite) {
		t.Run("graceful termination of single pod", func(t *testing.T) {
			ctx := context.Background()
			ns := "gateway-conformance-infra"
			path := "/graceful-termination-sts"
			routeNN := types.NamespacedName{Name: "graceful-termination-sts", Namespace: ns}
			gwNN := types.NamespacedName{Name: "same-namespace", Namespace: ns}
			gwAddr := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gwNN), routeNN)

			kubernetes.HTTPRouteMustHaveResolvedRefsConditionsTrue(t, suite.Client, suite.TimeoutConfig, routeNN, gwNN)

			makeRequestEventually(t, gwAddr, path, ns, suite)

			err := suite.Client.Delete(ctx, &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "backend-termination-sts-0",
					Namespace: ns,
				},
			})
			if err != nil {
				t.Errorf("error deleting pod: %v", err)
			}

			client := &nethttp.Client{}

			// while preStop is running, the server should be reachable
			awaitPodState(t, 10*time.Second, gwAddr, path, client, "stopping", func(podState CapturedServerState) bool {
				return podState.Ready && podState.Stopping
			})
			makeRequestEventually(t, gwAddr, path, ns, suite)

			// the server should be reachable even after SIGTERM was already sent to the container
			awaitPodState(t, 15*time.Second, gwAddr, path, client, "terminating", func(podState CapturedServerState) bool {
				return podState.Ready && podState.Terminating
			})
			makeRequestEventually(t, gwAddr, path, ns, suite)

			// the readinessProbe is configured to allow 10 seconds of failures, we should be able to observe this state
			awaitPodState(t, 15*time.Second, gwAddr, path, client, "not ready", func(podState CapturedServerState) bool {
				return !podState.Ready && podState.Terminating
			})
			makeRequestEventually(t, gwAddr, path, ns, suite)

			// the pod is eventually replaced by the statefulset once deletion completes, and it should become accessible
			awaitPodState(t, 2*time.Minute, gwAddr, path, client, "replaced", func(podState CapturedServerState) bool {
				return !podState.Terminating && !podState.Stopping
			})
			makeRequestEventually(t, gwAddr, path, ns, suite)
		})
		t.Run("graceful update of a tiny slow deployment", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			ns := "gateway-conformance-infra"
			path := "/graceful-termination"
			routeNN := types.NamespacedName{Name: "graceful-termination", Namespace: ns}
			gwNN := types.NamespacedName{Name: "same-namespace", Namespace: ns}
			deploymentNN := types.NamespacedName{Name: "backend-termination", Namespace: ns}
			gwAddr := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gwNN), routeNN)

			kubernetes.HTTPRouteMustHaveResolvedRefsConditionsTrue(t, suite.Client, suite.TimeoutConfig, routeNN, gwNN)

			makeRequestEventually(t, gwAddr, path, ns, suite)

			cReq, err := makeRequest(t, gwAddr, path, ns, suite)
			require.NoError(t, err)

			observedPods := sets.New[string](cReq.Pod)

			deployment := appsv1.Deployment{}
			err = suite.Client.Get(ctx, deploymentNN, &deployment)
			require.NoError(t, err)

			deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
				Name: "AN_APP_UPDATE", Value: time.Now().Format(time.RFC3339),
			})

			err = suite.Client.Update(ctx, &deployment)
			require.NoError(t, err)

			awaitDeploymentAvailable(t, ctx, gwAddr, path, ns, deploymentNN, suite, observedPods)
		})
		t.Run("graceful recovery of a tiny slow deployment", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			ns := "gateway-conformance-infra"
			path := "/graceful-termination"
			routeNN := types.NamespacedName{Name: "graceful-termination", Namespace: ns}
			gwNN := types.NamespacedName{Name: "same-namespace", Namespace: ns}
			deploymentNN := types.NamespacedName{Name: "backend-termination", Namespace: ns}
			gwAddr := kubernetes.GatewayAndHTTPRoutesMustBeAccepted(t, suite.Client, suite.TimeoutConfig, suite.ControllerName, kubernetes.NewGatewayRef(gwNN), routeNN)

			kubernetes.HTTPRouteMustHaveResolvedRefsConditionsTrue(t, suite.Client, suite.TimeoutConfig, routeNN, gwNN)

			makeRequestEventually(t, gwAddr, path, ns, suite)

			cReq, err := makeRequest(t, gwAddr, path, ns, suite)
			require.NoError(t, err)

			observedPods := sets.New[string](cReq.Pod)

			err = suite.Client.Delete(ctx, &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: cReq.Pod, Namespace: ns},
			})
			if err != nil {
				t.Errorf("error deleting pod: %v", err)
			}

			awaitDeploymentAvailable(t, ctx, gwAddr, path, ns, deploymentNN, suite, observedPods)
		})
	},
}

func awaitDeploymentAvailable(t *testing.T, ctx context.Context, gwAddr string, path string, ns string, deploymentNN types.NamespacedName, suite *suite.ConformanceTestSuite, observedPods sets.Set[string]) {
	labelSelector := "app=backend-termination"

	podSelector, err := labels.Parse(labelSelector)
	require.NoError(t, err)

	podWatcher, err := suite.Clientset.CoreV1().Pods(ns).Watch(ctx, metav1.ListOptions{Watch: true, LabelSelector: labelSelector})
	require.NoError(t, err)
	defer podWatcher.Stop()

	// until the deployment and all pods are ready, check that no request ever fails
	for {
		cReq, err := makeRequest(t, gwAddr, path, ns, suite)
		require.NoError(t, err)
		observedPods.Insert(cReq.Pod)

		deployment := appsv1.Deployment{}
		err = suite.Client.Get(ctx, deploymentNN, &deployment)
		require.NoError(t, err)

		podList := corev1.PodList{}
		err = suite.Client.List(ctx, &podList, &client.ListOptions{Namespace: ns, LabelSelector: podSelector})
		require.NoError(t, err)

		if podsAndDeploymentFullyReady(t, podList.Items, deployment, observedPods) {
			break
		}

		select {
		case <-time.After(1 * time.Second):
		case evt := <-podWatcher.ResultChan():
			if evt.Type == watch.Error {
				require.NoError(t, errors.FromObject(evt.Object))
			}
			pod, ok := evt.Object.(*corev1.Pod)
			if !ok {
				t.Errorf("unexpected watch event %+v", evt)
			}
			tlog.Logf(t, "pod event name=%s, type=%s, status=%+v", pod.Name, evt.Type, pod.Status)
		case <-ctx.Done():
			require.NoError(t, ctx.Err())
			t.FailNow()
			return
		}
	}
}

func podsAndDeploymentFullyReady(t *testing.T, pods []corev1.Pod, deployment appsv1.Deployment, observedPods sets.Set[string]) bool {
	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			tlog.Logf(t, "a pod is not running: %s Phase=%s", pod.Name, pod.Status.Phase)
			return false
		}
		for _, podCondition := range pod.Status.Conditions {
			if podCondition.Type == corev1.PodReady && podCondition.Status != corev1.ConditionTrue {
				tlog.Logf(t, "a pod is not ready: %s ConditionStatus=%s", pod.Name, podCondition.Status)
				return false
			}
		}
		if pod.ObjectMeta.DeletionTimestamp != nil {
			tlog.Logf(t, "a pod is in deletion: %s DeletionTimestamp=%s", pod.Name, pod.ObjectMeta.DeletionTimestamp)
			return false
		}
	}
	if deployment.Status.Replicas != int32(len(pods)) {
		tlog.Logf(t, "found pods do not match expected replicas: %d != %d", deployment.Status.Replicas, len(pods))
		return false
	}

	if deployment.Status.UnavailableReplicas > 0 ||
		deployment.Status.ReadyReplicas != deployment.Status.Replicas ||
		deployment.Status.UpdatedReplicas != deployment.Status.Replicas ||
		deployment.Status.AvailableReplicas != deployment.Status.Replicas {
		tlog.Logf(t, "deployment status is not synchronized: %+v", deployment.Status)
		return false
	}

	if observedPods.Len() <= int(deployment.Status.Replicas) {
		tlog.Logf(t, "not all pods were reached yet : len(%s) <= %d", observedPods.UnsortedList(), deployment.Status.Replicas)
		return false
	}

	return true
}

func awaitPodState(t *testing.T, threshold time.Duration, gwAddr string, path string, client *nethttp.Client, description string, condition func(podState CapturedServerState) bool) { // nolint: unparam
	http.AwaitConvergence(t, 1, threshold, func(elapsed time.Duration) bool {
		podStateReq, err := makeRequestForPodState(gwAddr, path, client)
		if err != nil {
			tlog.Logf(t, "Request failed: %v (after %v)", err.Error(), elapsed)
			return false
		}
		ready := condition(podStateReq)
		if !ready {
			tlog.Logf(t, "Awaiting for pod to be %s: %+v (after %v)", description, podStateReq, elapsed)
		} else {
			tlog.Logf(t, "pod is %s: %+v (after %v)", description, podStateReq, elapsed)
		}
		return ready
	})
}

func makeRequestForPodState(gwAddr string, path string, client *nethttp.Client) (cReq CapturedServerState, err error) {
	req, err := nethttp.NewRequest(nethttp.MethodGet, fmt.Sprintf("http://%s%s", gwAddr, path), nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != nethttp.StatusOK {
		return cReq, fmt.Errorf("response status code = %d", resp.StatusCode)
	}

	// Parse response body.
	if resp.Header.Get("Content-Type") == "application/json" {
		err = json.Unmarshal(body, &cReq)
	} else {
		return cReq, fmt.Errorf("unsupported response content type")
	}

	return
}

func makeRequest(t *testing.T, gwAddr string, path string, ns string, suite *suite.ConformanceTestSuite) (*roundtripper.CapturedRequest, error) {
	expectedResponse := http.ExpectedResponse{
		Request:         http.Request{Path: path},
		ExpectedRequest: &http.ExpectedRequest{Request: http.Request{Path: "/"}},
		Response:        http.Response{StatusCode: 200},
		Namespace:       ns,
	}
	req := http.MakeRequest(t, &expectedResponse, gwAddr, "HTTP", "http")
	cReq, cResp, err := suite.RoundTripper.CaptureRoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get expected response: %w", err)
	}

	if err := http.CompareRequest(t, &req, cReq, cResp, expectedResponse); err != nil {
		return nil, fmt.Errorf("failed to compare request and response: %w", err)
	}
	return cReq, nil
}

func makeRequestEventually(t *testing.T, gwAddr string, path string, ns string, suite *suite.ConformanceTestSuite) { // nolint: unparam
	http.MakeRequestAndExpectEventuallyConsistentResponse(t, suite.RoundTripper, suite.TimeoutConfig, gwAddr, http.ExpectedResponse{
		Request:         http.Request{Path: path},
		ExpectedRequest: &http.ExpectedRequest{Request: http.Request{Path: "/"}},
		Response:        http.Response{StatusCode: 200},
		Namespace:       ns,
	})
}
