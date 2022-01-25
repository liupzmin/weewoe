// Code generated by pegomock. DO NOT EDIT.
// Source: github.com/liupzmin/weewoe/internal/client (interfaces: Connection)

package config_test

import (
	"reflect"
	"time"

	client "github.com/liupzmin/weewoe/internal/client"
	pegomock "github.com/petergtz/pegomock"
	v1 "k8s.io/api/core/v1"
	version "k8s.io/apimachinery/pkg/version"
	disk "k8s.io/client-go/discovery/cached/disk"
	dynamic "k8s.io/client-go/dynamic"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

type MockConnection struct {
	fail func(message string, callerSkip ...int)
}

func NewMockConnection(options ...pegomock.Option) *MockConnection {
	mock := &MockConnection{}
	for _, option := range options {
		option.Apply(mock)
	}
	return mock
}

func (mock *MockConnection) SetFailHandler(fh pegomock.FailHandler) { mock.fail = fh }
func (mock *MockConnection) FailHandler() pegomock.FailHandler      { return mock.fail }

func (mock *MockConnection) ActiveCluster() string {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("ActiveCluster", params, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem()})
	var ret0 string
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(string)
		}
	}
	return ret0
}

func (mock *MockConnection) ActiveNamespace() string {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("ActiveNamespace", params, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem()})
	var ret0 string
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(string)
		}
	}
	return ret0
}

func (mock *MockConnection) CachedDiscovery() (*disk.CachedDiscoveryClient, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("CachedDiscovery", params, []reflect.Type{reflect.TypeOf((**disk.CachedDiscoveryClient)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 *disk.CachedDiscoveryClient
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(*disk.CachedDiscoveryClient)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockConnection) CanI(_param0 string, _param1 string, _param2 []string) (bool, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{_param0, _param1, _param2}
	result := pegomock.GetGenericMockFrom(mock).Invoke("CanI", params, []reflect.Type{reflect.TypeOf((*bool)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 bool
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(bool)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockConnection) CheckConnectivity() bool {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("CheckConnectivity", params, []reflect.Type{reflect.TypeOf((*bool)(nil)).Elem()})
	var ret0 bool
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(bool)
		}
	}
	return ret0
}

func (mock *MockConnection) Config() *client.Config {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("Config", params, []reflect.Type{reflect.TypeOf((**client.Config)(nil)).Elem()})
	var ret0 *client.Config
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(*client.Config)
		}
	}
	return ret0
}

func (mock *MockConnection) ConnectionOK() bool {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("ConnectionOK", params, []reflect.Type{reflect.TypeOf((*bool)(nil)).Elem()})
	var ret0 bool
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(bool)
		}
	}
	return ret0
}

func (mock *MockConnection) Dial() (kubernetes.Interface, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("Dial", params, []reflect.Type{reflect.TypeOf((*kubernetes.Interface)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 kubernetes.Interface
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(kubernetes.Interface)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockConnection) DialLogs() (kubernetes.Interface, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("DialLogs", params, []reflect.Type{reflect.TypeOf((*kubernetes.Interface)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 kubernetes.Interface
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(kubernetes.Interface)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockConnection) DynDial() (dynamic.Interface, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("DynDial", params, []reflect.Type{reflect.TypeOf((*dynamic.Interface)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 dynamic.Interface
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(dynamic.Interface)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockConnection) HasMetrics() bool {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("HasMetrics", params, []reflect.Type{reflect.TypeOf((*bool)(nil)).Elem()})
	var ret0 bool
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(bool)
		}
	}
	return ret0
}

func (mock *MockConnection) IsActiveNamespace(_param0 string) bool {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{_param0}
	result := pegomock.GetGenericMockFrom(mock).Invoke("IsActiveNamespace", params, []reflect.Type{reflect.TypeOf((*bool)(nil)).Elem()})
	var ret0 bool
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(bool)
		}
	}
	return ret0
}

func (mock *MockConnection) MXDial() (*versioned.Clientset, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("MXDial", params, []reflect.Type{reflect.TypeOf((**versioned.Clientset)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 *versioned.Clientset
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(*versioned.Clientset)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockConnection) RestConfig() (*rest.Config, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("RestConfig", params, []reflect.Type{reflect.TypeOf((**rest.Config)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 *rest.Config
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(*rest.Config)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockConnection) ServerVersion() (*version.Info, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("ServerVersion", params, []reflect.Type{reflect.TypeOf((**version.Info)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 *version.Info
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(*version.Info)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockConnection) SwitchContext(_param0 string) error {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{_param0}
	result := pegomock.GetGenericMockFrom(mock).Invoke("SwitchContext", params, []reflect.Type{reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(error)
		}
	}
	return ret0
}

func (mock *MockConnection) ValidNamespaces() ([]v1.Namespace, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockConnection().")
	}
	params := []pegomock.Param{}
	result := pegomock.GetGenericMockFrom(mock).Invoke("ValidNamespaces", params, []reflect.Type{reflect.TypeOf((*[]v1.Namespace)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 []v1.Namespace
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].([]v1.Namespace)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockConnection) VerifyWasCalledOnce() *VerifierMockConnection {
	return &VerifierMockConnection{
		mock:                   mock,
		invocationCountMatcher: pegomock.Times(1),
	}
}

func (mock *MockConnection) VerifyWasCalled(invocationCountMatcher pegomock.Matcher) *VerifierMockConnection {
	return &VerifierMockConnection{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
	}
}

func (mock *MockConnection) VerifyWasCalledInOrder(invocationCountMatcher pegomock.Matcher, inOrderContext *pegomock.InOrderContext) *VerifierMockConnection {
	return &VerifierMockConnection{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		inOrderContext:         inOrderContext,
	}
}

func (mock *MockConnection) VerifyWasCalledEventually(invocationCountMatcher pegomock.Matcher, timeout time.Duration) *VerifierMockConnection {
	return &VerifierMockConnection{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		timeout:                timeout,
	}
}

type VerifierMockConnection struct {
	mock                   *MockConnection
	invocationCountMatcher pegomock.Matcher
	inOrderContext         *pegomock.InOrderContext
	timeout                time.Duration
}

func (verifier *VerifierMockConnection) ActiveCluster() *MockConnection_ActiveCluster_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "ActiveCluster", params, verifier.timeout)
	return &MockConnection_ActiveCluster_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_ActiveCluster_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_ActiveCluster_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_ActiveCluster_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) ActiveNamespace() *MockConnection_ActiveNamespace_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "ActiveNamespace", params, verifier.timeout)
	return &MockConnection_ActiveNamespace_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_ActiveNamespace_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_ActiveNamespace_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_ActiveNamespace_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) CachedDiscovery() *MockConnection_CachedDiscovery_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "CachedDiscovery", params, verifier.timeout)
	return &MockConnection_CachedDiscovery_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_CachedDiscovery_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_CachedDiscovery_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_CachedDiscovery_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) CanI(_param0 string, _param1 string, _param2 []string) *MockConnection_CanI_OngoingVerification {
	params := []pegomock.Param{_param0, _param1, _param2}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "CanI", params, verifier.timeout)
	return &MockConnection_CanI_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_CanI_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_CanI_OngoingVerification) GetCapturedArguments() (string, string, []string) {
	_param0, _param1, _param2 := c.GetAllCapturedArguments()
	return _param0[len(_param0)-1], _param1[len(_param1)-1], _param2[len(_param2)-1]
}

func (c *MockConnection_CanI_OngoingVerification) GetAllCapturedArguments() (_param0 []string, _param1 []string, _param2 [][]string) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]string, len(c.methodInvocations))
		for u, param := range params[0] {
			_param0[u] = param.(string)
		}
		_param1 = make([]string, len(c.methodInvocations))
		for u, param := range params[1] {
			_param1[u] = param.(string)
		}
		_param2 = make([][]string, len(c.methodInvocations))
		for u, param := range params[2] {
			_param2[u] = param.([]string)
		}
	}
	return
}

func (verifier *VerifierMockConnection) CheckConnectivity() *MockConnection_CheckConnectivity_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "CheckConnectivity", params, verifier.timeout)
	return &MockConnection_CheckConnectivity_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_CheckConnectivity_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_CheckConnectivity_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_CheckConnectivity_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) Config() *MockConnection_Config_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Config", params, verifier.timeout)
	return &MockConnection_Config_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_Config_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_Config_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_Config_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) ConnectionOK() *MockConnection_ConnectionOK_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "ConnectionOK", params, verifier.timeout)
	return &MockConnection_ConnectionOK_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_ConnectionOK_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_ConnectionOK_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_ConnectionOK_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) Dial() *MockConnection_Dial_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Dial", params, verifier.timeout)
	return &MockConnection_Dial_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_Dial_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_Dial_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_Dial_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) DynDial() *MockConnection_DynDial_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "DynDial", params, verifier.timeout)
	return &MockConnection_DynDial_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_DynDial_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_DynDial_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_DynDial_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) HasMetrics() *MockConnection_HasMetrics_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "HasMetrics", params, verifier.timeout)
	return &MockConnection_HasMetrics_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_HasMetrics_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_HasMetrics_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_HasMetrics_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) IsActiveNamespace(_param0 string) *MockConnection_IsActiveNamespace_OngoingVerification {
	params := []pegomock.Param{_param0}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "IsActiveNamespace", params, verifier.timeout)
	return &MockConnection_IsActiveNamespace_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_IsActiveNamespace_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_IsActiveNamespace_OngoingVerification) GetCapturedArguments() string {
	_param0 := c.GetAllCapturedArguments()
	return _param0[len(_param0)-1]
}

func (c *MockConnection_IsActiveNamespace_OngoingVerification) GetAllCapturedArguments() (_param0 []string) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]string, len(c.methodInvocations))
		for u, param := range params[0] {
			_param0[u] = param.(string)
		}
	}
	return
}

func (verifier *VerifierMockConnection) MXDial() *MockConnection_MXDial_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "MXDial", params, verifier.timeout)
	return &MockConnection_MXDial_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_MXDial_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_MXDial_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_MXDial_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) RestConfig() *MockConnection_RestConfig_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "RestConfig", params, verifier.timeout)
	return &MockConnection_RestConfig_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_RestConfig_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_RestConfig_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_RestConfig_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) ServerVersion() *MockConnection_ServerVersion_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "ServerVersion", params, verifier.timeout)
	return &MockConnection_ServerVersion_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_ServerVersion_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_ServerVersion_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_ServerVersion_OngoingVerification) GetAllCapturedArguments() {
}

func (verifier *VerifierMockConnection) SwitchContext(_param0 string) *MockConnection_SwitchContext_OngoingVerification {
	params := []pegomock.Param{_param0}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "SwitchContext", params, verifier.timeout)
	return &MockConnection_SwitchContext_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_SwitchContext_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_SwitchContext_OngoingVerification) GetCapturedArguments() string {
	_param0 := c.GetAllCapturedArguments()
	return _param0[len(_param0)-1]
}

func (c *MockConnection_SwitchContext_OngoingVerification) GetAllCapturedArguments() (_param0 []string) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]string, len(c.methodInvocations))
		for u, param := range params[0] {
			_param0[u] = param.(string)
		}
	}
	return
}

func (verifier *VerifierMockConnection) ValidNamespaces() *MockConnection_ValidNamespaces_OngoingVerification {
	params := []pegomock.Param{}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "ValidNamespaces", params, verifier.timeout)
	return &MockConnection_ValidNamespaces_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockConnection_ValidNamespaces_OngoingVerification struct {
	mock              *MockConnection
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockConnection_ValidNamespaces_OngoingVerification) GetCapturedArguments() {
}

func (c *MockConnection_ValidNamespaces_OngoingVerification) GetAllCapturedArguments() {
}
