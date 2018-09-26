package kube

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"code.cloudfoundry.org/fissile/testhelpers"
)

func TestMakeRegistryCredentialsKube(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	registryCredentials, err := MakeRegistryCredentials(ExportSettings{})

	if !assert.NoError(err) {
		return
	}

	actual, err := RoundtripKube(registryCredentials)
	if !assert.NoError(err) {
		return
	}
	testhelpers.IsYAMLEqualString(assert, `---
		apiVersion: "v1"
		data:
			.dockercfg: ""
		kind: "Secret"
		metadata:
			name: "registry-credentials"
			labels:
				skiff-role-name: "registry-credentials"
		type: "kubernetes.io/dockercfg"
	`, actual)
}

func TestMakeRegistryCredentialsHelm(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	registryCredentials, err := MakeRegistryCredentials(ExportSettings{
		CreateHelmChart: true,
	})
	if !assert.NoError(err) {
		return
	}

	user := "the-user"
	pass := "the-password"
	host := "the-host"

	auth64 := RenderEncodeBase64(fmt.Sprintf("%s:%s", user, pass))
	dcfg := RenderEncodeBase64(fmt.Sprintf(
		`{%q:{"username":%q,"password":%q,"auth":%q}}`,
		host, user, pass, auth64))

	config := map[string]interface{}{
		"Values.kube.registry.hostname": host,
		"Values.kube.registry.username": user,
		"Values.kube.registry.password": pass,
	}

	actual, err := RoundtripNode(registryCredentials, config)
	if !assert.NoError(err) {
		return
	}

	testhelpers.IsYAMLEqualString(assert, fmt.Sprintf(`---
		apiVersion: "v1"
		data:
			.dockercfg: %s
		kind: "Secret"
		metadata:
			name: "registry-credentials"
			labels:
				skiff-role-name: "registry-credentials"
		type: "kubernetes.io/dockercfg"
	`, dcfg), actual)
}
