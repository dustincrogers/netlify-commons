package ntoml

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMissingFile(t *testing.T) {
	_, err := LoadFrom("does-not-exist")
	assert.True(t, os.IsNotExist(err))
}

func TestLoadingExampleToml(t *testing.T) {
	tmp := testToml(t)
	defer os.Remove(t.Name())

	conf, err := LoadFrom(tmp.Name())
	require.NoError(t, err)

	expected := &NetlifyToml{
		Settings: Settings{
			ID:   "this-is-a-site",
			Path: ".",
		},
		Redirects: []Redirect{
			{Origin: "/other", Destination: "/otherpage.html", Force: true},
		},
		Build: &BuildConfig{
			Command: "echo 'not a thing'",
		},
		Context: map[string]DeployContext{
			"deploy-preview": {BuildConfig{
				Command: "hugo version && npm run build-preview",
			}},
			"branch-deploy": {BuildConfig{
				Command:     "hugo version && npm run build-branch",
				Environment: map[string]string{"HUGO_VERSION": "0.20.5"},
			}},
		},
	}

	assert.Equal(t, expected, conf)
}

func TestLoadingExampleJSON(t *testing.T) {
	tmp := testJSON(t)
	defer os.Remove(t.Name())

	conf, err := LoadFrom(tmp.Name())
	require.NoError(t, err)

	expected := &NetlifyToml{
		Settings: Settings{
			ID:   "this-is-a-site",
			Path: ".",
		},
		Redirects: []Redirect{
			{Origin: "/other", Destination: "/otherpage.html", Force: true},
		},
		Build: &BuildConfig{
			Command: "echo 'not a thing'",
		},
		Context: map[string]DeployContext{
			"deploy-preview": {BuildConfig{
				Command: "hugo version && npm run build-preview",
			}},
			"branch-deploy": {BuildConfig{
				Command:     "hugo version && npm run build-branch",
				Environment: map[string]string{"HUGO_VERSION": "0.20.5"},
			}},
		},
	}

	assert.Equal(t, expected, conf)
}

func TestLoadingExampleYAML(t *testing.T) {
	tmp := testYAML(t)
	defer os.Remove(t.Name())

	conf, err := LoadFrom(tmp.Name())
	require.NoError(t, err)

	expected := &NetlifyToml{
		Settings: Settings{
			ID:   "this-is-a-site",
			Path: ".",
		},
		Redirects: []Redirect{
			{Origin: "/other", Destination: "/otherpage.html", Force: true},
		},
		Build: &BuildConfig{
			Command: "echo 'not a thing'",
		},
		Context: map[string]DeployContext{
			"deploy-preview": {BuildConfig{
				Command: "hugo version && npm run build-preview",
			}},
			"branch-deploy": {BuildConfig{
				Command:     "hugo version && npm run build-branch",
				Environment: map[string]string{"HUGO_VERSION": "0.20.5"},
			}},
		},
	}

	assert.Equal(t, expected, conf)
}

func TestSaveTomlFile(t *testing.T) {
	conf := &NetlifyToml{
		Settings: Settings{ID: "This is something", Path: "/dist"},
	}

	tmp, err := ioutil.TempFile("", "netlify-ctl")
	require.NoError(t, err)

	require.NoError(t, SaveTo(conf, tmp.Name()))

	data, err := ioutil.ReadFile(tmp.Name())
	require.NoError(t, err)

	expected := "[settings]\n" +
		"  id = \"This is something\"\n" +
		"  path = \"/dist\"\n"

	assert.Equal(t, expected, string(data))
}

func testToml(t *testing.T) *os.File {
	tmp, err := ioutil.TempFile("", "netlify-ctl-*.toml")
	require.NoError(t, err)

	data := `
[Settings]
  id = "this-is-a-site"
  path = "."

[build]
	command = "echo 'not a thing'"

[[redirects]]
  origin = "/other"
	force = true
  destination = "/otherpage.html"

  [context.deploy-preview]
  command = "hugo version && npm run build-preview"


[context.branch-deploy]
  command = "hugo version && npm run build-branch"

  [context.branch-deploy.environment]
    HUGO_VERSION = "0.20.5"
`
	require.NoError(t, ioutil.WriteFile(tmp.Name(), []byte(data), 0664))

	return tmp
}

func testJSON(t *testing.T) *os.File {
	tmp, err := ioutil.TempFile("", "netlify-ctl-*.json")
	require.NoError(t, err)

	data := `
{
  "settings": {
    "id": "this-is-a-site",
    "path": "."
  },
  "redirects": [
    {
      "origin": "/other",
      "destination": "/otherpage.html",
      "force": true
    }
  ],
  "build": {
    "command": "echo 'not a thing'"
  },
  "context": {
    "deploy-preview": {
      "command": "hugo version && npm run build-preview"
    },
    "branch-deploy": {
      "command": "hugo version && npm run build-branch",
      "environment": {
        "HUGO_VERSION": "0.20.5"
      }
    }
  }
}
`

	require.NoError(t, ioutil.WriteFile(tmp.Name(), []byte(data), 0664))

	return tmp
}

func testYAML(t *testing.T) *os.File {
	tmp, err := ioutil.TempFile("", "netlify-ctl-*.yaml")
	require.NoError(t, err)

	data := `
settings:
  id: "this-is-a-site"
  path: "."

redirects:
  - origin: "/other"
    destination: "/otherpage.html"
    force: true

build:
  command: "echo 'not a thing'"

context:
  deploy-preview:
    command: "hugo version && npm run build-preview"
  branch-deploy:
    command: "hugo version && npm run build-branch"
    environment:
      HUGO_VERSION: "0.20.5"
`

	require.NoError(t, ioutil.WriteFile(tmp.Name(), []byte(data), 0664))

	return tmp
}

func TestFindOnlyOneExistingPath(t *testing.T) {
	_, err := findOnlyOneExistingPath("", "does-not-exist")
	assert.IsType(t, &FoundNoConfigPathError{}, err)

	tmp, err := ioutil.TempFile("", "netlify-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmp.Name())
	path, err := findOnlyOneExistingPath("", tmp.Name())
	assert.NoError(t, err)
	assert.Equal(t, tmp.Name(), path)
}
