package kube

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SUSE/fissile/helm"
	"github.com/SUSE/fissile/model"
	"github.com/SUSE/fissile/testhelpers"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func jobTestLoadRole(assert *assert.Assertions, roleName, manifestName string) *model.Role {
	workDir, err := os.Getwd()
	assert.NoError(err)

	manifestPath := filepath.Join(workDir, "../test-assets/role-manifests", manifestName)
	releasePath := filepath.Join(workDir, "../test-assets/tor-boshrelease")
	releasePathBoshCache := filepath.Join(releasePath, "bosh-cache")

	release, err := model.NewDevRelease(releasePath, "", "", releasePathBoshCache)
	if !assert.NoError(err) {
		return nil
	}

	manifest, err := model.LoadRoleManifest(manifestPath, []*model.Release{release})
	if !assert.NoError(err) {
		return nil
	}

	role := manifest.LookupRole(roleName)
	if !assert.NotNil(role, "Failed to find role %s", roleName) {
		return nil
	}
	return role
}

func TestJobPreFlight(t *testing.T) {
	assert := assert.New(t)
	role := jobTestLoadRole(assert, "pre-role", "jobs.yml")
	if role == nil {
		return
	}

	job, err := NewJob(role, ExportSettings{
		Opinions: model.NewEmptyOpinions(),
	})
	if !assert.NoError(err, "Failed to create job from role pre-role") {
		return
	}
	assert.NotNil(job)

	yamlConfig := &bytes.Buffer{}
	if err := helm.NewEncoder(yamlConfig).Encode(job); !assert.NoError(err) {
		return
	}

	var expected, actual interface{}
	if !assert.NoError(yaml.Unmarshal(yamlConfig.Bytes(), &actual)) {
		return
	}

	expectedYAML := strings.Replace(`---
		apiVersion: extensions/v1beta1
		kind: Job
		metadata:
			name: pre-role
		spec:
			template:
				metadata:
					name: pre-role
				spec:
					containers:
					-
						name: pre-role
					restartPolicy: OnFailure
	`, "\t", "    ", -1)
	if !assert.NoError(yaml.Unmarshal([]byte(expectedYAML), &expected)) {
		return
	}
	testhelpers.IsYAMLSubset(assert, expected, actual)
}

func TestJobPostFlight(t *testing.T) {
	assert := assert.New(t)
	role := jobTestLoadRole(assert, "post-role", "jobs.yml")
	if role == nil {
		return
	}

	job, err := NewJob(role, ExportSettings{
		Opinions: model.NewEmptyOpinions(),
	})
	if !assert.NoError(err, "Failed to create job from role post-role") {
		return
	}
	assert.NotNil(job)

	yamlConfig := &bytes.Buffer{}
	if err := helm.NewEncoder(yamlConfig).Encode(job); !assert.NoError(err) {
		return
	}

	var expected, actual interface{}
	if !assert.NoError(yaml.Unmarshal(yamlConfig.Bytes(), &actual)) {
		return
	}
	expectedYAML := strings.Replace(`---
		apiVersion: extensions/v1beta1
		kind: Job
		metadata:
			name: post-role
		spec:
			template:
				metadata:
					name: post-role
				spec:
					containers:
					-
						name: post-role
					restartPolicy: OnFailure
	`, "\t", "    ", -1)
	if !assert.NoError(yaml.Unmarshal([]byte(expectedYAML), &expected)) {
		return
	}
	testhelpers.IsYAMLSubset(assert, expected, actual)
}

func TestJobWithAnnotations(t *testing.T) {
	assert := assert.New(t)

	role := jobTestLoadRole(assert, "role", "job-with-annotation.yml")
	if role == nil {
		return
	}

	job, err := NewJob(role, ExportSettings{
		Opinions: model.NewEmptyOpinions(),
	})
	if !assert.NoError(err, "Failed to create job from role pre-role") {
		return
	}
	assert.NotNil(job)

	yamlConfig := &bytes.Buffer{}
	if err := helm.NewEncoder(yamlConfig).Encode(job); !assert.NoError(err) {
		return
	}

	var expected, actual interface{}
	if !assert.NoError(yaml.Unmarshal(yamlConfig.Bytes(), &actual)) {
		return
	}
	expectedYAML := strings.Replace(`---
		apiVersion: extensions/v1beta1
		kind: Job
		metadata:
			name: role
			annotations:
				helm.sh/hook: post-install
	`, "\t", "    ", -1)
	if !assert.NoError(yaml.Unmarshal([]byte(expectedYAML), &expected)) {
		return
	}

	testhelpers.IsYAMLSubset(assert, expected, actual)
}
