// Package swift implements the "swift" runtime.
package swift

import (
	"io"
	"net/http"
	"os"
	"os/user"

	"github.com/apex/apex/archive"
	"github.com/apex/apex/function"
	"github.com/apex/apex/plugins/nodejs"
)

const (
	// Runtime name used by Apex
	Runtime = "swift"
)

func init() {
	function.RegisterPlugin(Runtime, &Plugin{})
}

var functions = make(map[string]bool)

var libraries = []string {
	"libswiftCore.so",
	"libFoundation.so",
	"libicudata.so.52",
	"libicui18n.so.52",
	"libicuuc.so.52",
	"libswiftGlibc.so",
}

// Plugin implementation.
type Plugin struct{}

// Open adds the shim and swift defaults.
func (p *Plugin) Open(fn *function.Function) error {
	if fn.Runtime != Runtime {
		return nil
	}

	functions[fn.Name] = true

	fn.Shim = true
	fn.Runtime = nodejs.Runtime

	if fn.Hooks.Build == "" {
		fn.Hooks.Build = "swift build -c release; mv .build/release/main main; mv .build/release/*.so ."
	}

	if fn.Hooks.Clean == "" {
		fn.Hooks.Clean = "swift build --clean=dist; rm -f main; rm -f *.so"
	}

	return nil
}

// Build adds the swift libraries.
func (p *Plugin) Build(fn *function.Function, zip *archive.Zip) error {
	_, isSwiftFunction := functions[fn.Name]

	if isSwiftFunction {
		fn.Log.Debug("adding swift libraries")

		for _, library := range libraries {
			addFile(library, zip)
		}
	}

	return nil
}

func addFile(fileName string, zip *archive.Zip) error {
	file, err := getFile(fileName)

	if err != nil {
		return err
	}

	if err := zip.AddFile(fileName, file); err != nil {
		return err
	}

	return nil
}

func getFile(fileName string) (*os.File, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir + "/.apex/plugins/swift/"
	path := dir + fileName

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)

		if err != nil {
			return nil, err
		}

		file, err = download(fileName, file)

		if err != nil {
			return nil, err
		}

		return file, nil
	}

	return os.Open(path)
}

func download(fileName string, file *os.File) (*os.File, error) {
	url := "https://raw.githubusercontent.com/paulofaria/swift-binaries/master/" + fileName

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return nil, err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}

	return file, nil
}
