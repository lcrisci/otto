package directory

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

// TestBackend is a public test helper that verifies a backend
// functions properly.
func TestBackend(t *testing.T, b Backend) {
	var buf bytes.Buffer
	var err error

	//---------------------------------------------------------------
	// Blob
	//---------------------------------------------------------------

	// GetBlob (doesn't exist)
	data, err := b.GetBlob("foo")
	if err != nil {
		t.Fatalf("GetBlob error: %s", err)
	}
	if data != nil {
		data.Close()
		t.Fatalf("GetBlob should be nil data")
	}

	// PutBlob
	err = b.PutBlob("foo", &BlobData{Data: strings.NewReader("bar")})
	if err != nil {
		t.Fatalf("PutBlob error: %s", err)
	}

	// GetBlob (exists)
	data, err = b.GetBlob("foo")
	if err != nil {
		t.Fatalf("GetBlob error: %s", err)
	}
	_, err = io.Copy(&buf, data.Data)
	data.Close()
	if err != nil {
		t.Fatalf("GetBlob error: %s", err)
	}
	if buf.String() != "bar" {
		t.Fatalf("GetBlob bad data: %s", buf.String())
	}

	//---------------------------------------------------------------
	// Infra
	//---------------------------------------------------------------

	// GetInfra (doesn't exist)
	infra := &Infra{Lookup: Lookup{Infra: "foo"}}
	actualInfra, err := b.GetInfra(infra)
	if err != nil {
		t.Fatalf("GetInfra (non-exist) error: %s", err)
	}
	if actualInfra != nil {
		t.Fatal("GetInfra (non-exist): infra should be nil")
	}

	// PutInfra (doesn't exist)
	infra.Outputs = map[string]string{"foo": "bar"}
	if infra.ID != "" {
		t.Fatalf("PutInfra: ID should be empty before set")
	}
	if err := b.PutInfra(infra); err != nil {
		t.Fatalf("PutInfra err: %s", err)
	}
	if infra.ID == "" {
		t.Fatalf("PutInfra: infra ID not set")
	}

	// GetInfra (exists)
	actualInfra, err = b.GetInfra(infra)
	if err != nil {
		t.Fatalf("GetInfra (exist) error: %s", err)
	}
	if !reflect.DeepEqual(actualInfra, infra) {
		t.Fatalf("GetInfra (exist) bad: %#v", actualInfra)
	}

	// GetInfra with foundation (doesn't exist)
	infra = &Infra{Lookup: Lookup{Infra: "foo", Foundation: "bar"}}
	actualInfra, err = b.GetInfra(infra)
	if err != nil {
		t.Fatalf("GetInfra (non-exist) error: %s", err)
	}
	if actualInfra != nil {
		t.Fatal("GetInfra (non-exist): infra should be nil")
	}

	// PutInfra with foundation (doesn't exist)
	infra.Outputs = map[string]string{"foo": "bar"}
	if infra.ID != "" {
		t.Fatalf("PutInfra: ID should be empty before set")
	}
	if err := b.PutInfra(infra); err != nil {
		t.Fatalf("PutInfra err: %s", err)
	}
	if infra.ID == "" {
		t.Fatalf("PutInfra: infra ID not set")
	}

	// GetInfra with foundation (exists)
	actualInfra, err = b.GetInfra(infra)
	if err != nil {
		t.Fatalf("GetInfra (exist) error: %s", err)
	}
	if !reflect.DeepEqual(actualInfra, infra) {
		t.Fatalf("GetInfra (exist) bad: %#v", actualInfra)
	}

	//---------------------------------------------------------------
	// Deploy
	//---------------------------------------------------------------

	// GetDeploy (doesn't exist)
	deploy := &Deploy{Lookup: Lookup{
		AppID: "foo", Infra: "bar", InfraFlavor: "baz"}}
	deployResult, err := b.GetDeploy(deploy)
	if err != nil {
		t.Fatalf("GetDeploy (non-exist) error: %s", err)
	}
	if deployResult != nil {
		t.Fatal("GetDeploy (non-exist): result should be nil")
	}

	// PutDeploy (doesn't exist)
	if deploy.ID != "" {
		t.Fatalf("PutDeploy: ID should be empty before set")
	}
	if err := b.PutDeploy(deploy); err != nil {
		t.Fatalf("PutDeploy err: %s", err)
	}
	if deploy.ID == "" {
		t.Fatalf("PutDeploy: deploy ID not set")
	}

	// GetDeploy (exists)
	deployResult, err = b.GetDeploy(deploy)
	if err != nil {
		t.Fatalf("GetDeploy (exist) error: %s", err)
	}
	if !reflect.DeepEqual(deployResult, deploy) {
		t.Fatalf("GetDeploy (exist) bad: %#v", deployResult)
	}

	//---------------------------------------------------------------
	// Dev
	//---------------------------------------------------------------

	// GetDev (doesn't exist)
	dev := &Dev{Lookup: Lookup{AppID: "foo"}}
	devResult, err := b.GetDev(dev)
	if err != nil {
		t.Fatalf("GetDev (non-exist) error: %s", err)
	}
	if devResult != nil {
		t.Fatal("GetDev (non-exist): result should be nil")
	}

	// PutDev (doesn't exist)
	if dev.ID != "" {
		t.Fatalf("PutDev: ID should be empty before set")
	}
	if err := b.PutDev(dev); err != nil {
		t.Fatalf("PutDev err: %s", err)
	}
	if dev.ID == "" {
		t.Fatalf("PutDev: dev ID not set")
	}

	// GetDev (exists)
	devResult, err = b.GetDev(dev)
	if err != nil {
		t.Fatalf("GetDev (exist) error: %s", err)
	}
	if !reflect.DeepEqual(devResult, dev) {
		t.Fatalf("GetDev (exist) bad: %#v", devResult)
	}

	// DeleteDev (exists)
	err = b.DeleteDev(dev)
	if err != nil {
		t.Fatalf("DeleteDev error: %s", err)
	}
	devResult, err = b.GetDev(dev)
	if err != nil {
		t.Fatalf("GetDev (non-exist) error: %s", err)
	}
	if devResult != nil {
		t.Fatal("GetDev (non-exist): result should be nil")
	}

	// DeleteDev (doesn't exist)
	err = b.DeleteDev(dev)
	if err != nil {
		t.Fatalf("DeleteDev error: %s", err)
	}
}
