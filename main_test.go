package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	t.Log("Build without app path")
	{

		build, err := build("", SAMPLE_TEST_SUITE, "username", "password")
		require.Equal(t, "", build)
		require.Error(t, err, FILE_NOT_AVAILABLE_ERROR)
	}

	t.Log("Build without test_suite_app path")
	{

		build, err := build(SAMPLE_APP, "", "", "")
		require.Equal(t, "", build)
		require.Error(t, err, FILE_NOT_AVAILABLE_ERROR)
	}

	t.Log("Build with invalid credentials")
	{
		build, err := build(SAMPLE_APP, SAMPLE_TEST_SUITE, "a", "a")

		require.Equal(t, build, "{\"message\":\"Unauthorized\"}")

		require.NoError(t, err)
	}

}

func TestUpload(t *testing.T) {
	t.Log("It should throw file not found error with empty path")
	{

		build, err := upload("", APP_UPLOAD_ENDPOINT, "username", "password")
		t.Log(build, err)
		require.Equal(t, "", build)
		require.Error(t, err)
	}

	t.Log("It should throw file not found error with invalid path")
	{

		build, err := upload("invalidpath", APP_UPLOAD_ENDPOINT, "username", "password")

		t.Log(build, err)
		require.Equal(t, "", build)
		require.Error(t, err)
	}

}
