package snappy

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func Test2(t *testing.T) { TestingT(t) }

type HTestSuite struct{}

var _ = Suite(&HTestSuite{})

func (ts *HTestSuite) TestUnpack(c *C) {

	// setup tmpdir
	tmpdir, err := ioutil.TempDir(os.TempDir(), "meep")
	c.Assert(err, IsNil)
	defer os.RemoveAll(tmpdir)
	tmpfile := filepath.Join(tmpdir, "foo.tar.gz")

	// ok, slightly silly
	path := "/etc/fstab"

	// create test data
	cmd := exec.Command("tar", "cvzf", tmpfile, path)
	output, err := cmd.CombinedOutput()
	c.Assert(err, IsNil)
	if !strings.Contains(string(output), "/etc/fstab") {
		c.Error("Can not find expected output from tar")
	}

	// unpack
	unpackdir := filepath.Join(tmpdir, "t")
	err = unpackTar(tmpfile, unpackdir)
	c.Assert(err, IsNil)

	_, err = os.Open(filepath.Join(tmpdir, "t/etc/fstab"))
	c.Assert(err, IsNil)
}

func (ts *HTestSuite) TestGetMapFromValidYaml(c *C) {
	m, err := getMapFromYaml([]byte("name: value"))
	c.Assert(err, IsNil)
	me := map[string]interface{}{"name": "value"}
	if !reflect.DeepEqual(m, me) {
		c.Error(fmt.Sprintf("Unexpected map %v != %v", m, me))
	}
}

func (ts *HTestSuite) TestGetMapFromInvalidYaml(c *C) {
	_, err := getMapFromYaml([]byte("%lala%"))
	c.Assert(err, NotNil)
}
