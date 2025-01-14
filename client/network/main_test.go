package network

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	azure "github.com/virtual-kubelet/azure-aci/client"
	"github.com/virtual-kubelet/azure-aci/client/resourcegroups"
)

var (
	location      = "westus"
	resourceGroup = "virtual-node-test-rg"
	testAuth      *azure.Authentication
)

var defaultRetryConfig = azure.HTTPRetryConfig{
	RetryWaitMin: azure.DefaultRetryIntervalMin,
	RetryWaitMax: azure.DefaultRetryIntervalMax,
	RetryMax:     azure.DefaultRetryMax,
}

func TestMain(m *testing.M) {
	uid := uuid.New()
	resourceGroup += "-" + uid.String()[0:6]

	if err := setupAuth(); err != nil {
		fmt.Fprintln(os.Stderr, "Error setting up auth:", err)
		os.Exit(1)
	}

	c, err := resourcegroups.NewClient(testAuth, "unit-test", defaultRetryConfig)
	if err != nil {
		os.Exit(1)
	}
	_, err = c.CreateResourceGroup(resourceGroup, resourcegroups.Group{
		Name:     resourceGroup,
		Location: location,
	})
	if err != nil {
		os.Exit(1)
	}

	code := m.Run()

	if err := c.DeleteResourceGroup(resourceGroup); err != nil {
		fmt.Fprintln(os.Stderr, "error removing resource group:", err)
	}

	os.Exit(code)
}

var authOnce sync.Once

func setupAuth() error {
	var err error
	authOnce.Do(func() {
		testAuth, err = azure.NewAuthenticationFromFile(os.Getenv("AZURE_AUTH_LOCATION"))
		if err != nil {
			testAuth, err = azure.NewAuthenticationFromFile(os.Getenv("AZURE_AUTH_LOCATION"))
		}
		if err != nil {
			err = errors.Wrap(err, "failed to load Azure authentication file")
		}
	})
	return err
}

func newTestClient(t *testing.T) *Client {
	if err := setupAuth(); err != nil {
		t.Fatal(err)
	}
	c, err := NewClient(testAuth, "unit-test", defaultRetryConfig)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
