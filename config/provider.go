package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var providerMap = make(map[string]Provider)

func init() {
	RegisterProvider(newFileProvider())
}

// ProviderCallback is callback function for provider to handle config change
type ProviderCallback func(string, []byte)

type Provider interface {
	// Name provider name
	Name() string

	// Read reads the specific path file
	Read(string) ([]byte, error)

	// Watch watches config changing
	Watch(ProviderCallback)
}

func newFileProvider() *FileProvider {
	fp := &FileProvider{
		cb:              make(chan ProviderCallback),
		disabledWatcher: true,
		cache:           make(map[string]string),
		modtime:         make(map[string]int64),
	}
	if watcher, err := fsnotify.NewWatcher(); err == nil {
		fp.disabledWatcher = false
		fp.watcher = watcher
		go fp.run()
	}
	return fp
}

// RegisterProvider registers a provider by its name.
func RegisterProvider(p Provider) {
	providerMap[p.Name()] = p
}

// GetProvider get a provider by its name.
func GetProvider(name string) Provider {
	return providerMap[name]
}

// FileProvider is a config provider which gets config from file system.
type FileProvider struct {
	disabledWatcher bool
	watcher         *fsnotify.Watcher
	cb              chan ProviderCallback
	cache           map[string]string
	modtime         map[string]int64
	mu              sync.RWMutex
}

// Name returns file provider's name.
func (*FileProvider) Name() string {
	return "file"
}

// Read reads the specific path file, returns
// it content as bytes.
func (fp *FileProvider) Read(path string) ([]byte, error) {
	if !fp.disabledWatcher {
		if err := fp.watcher.Add(filepath.Dir(path)); err != nil {
			return nil, err
		}
		fp.mu.Lock()
		fp.cache[filepath.Clean(path)] = path
		fp.mu.Unlock()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Watch watches config changing. The change will
// be handled by callback function.
func (fp *FileProvider) Watch(cb ProviderCallback) {
	if !fp.disabledWatcher {
		fp.cb <- cb
	}
}

func (fp *FileProvider) run() {
	fn := make([]ProviderCallback, 0)
	for {
		select {
		case i := <-fp.cb:
			fn = append(fn, i)
		case e := <-fp.watcher.Events:
			if t, ok := fp.isModified(e); ok {
				fp.trigger(e, t, fn)
			}
		}
	}
}

func (fp *FileProvider) isModified(e fsnotify.Event) (int64, bool) {
	if e.Op&fsnotify.Write != fsnotify.Write {
		return 0, false
	}
	fp.mu.RLock()
	defer fp.mu.RUnlock()
	if _, ok := fp.cache[filepath.Clean(e.Name)]; !ok {
		return 0, false
	}
	fi, err := os.Stat(e.Name)
	if err != nil {
		return 0, false
	}
	if fi.ModTime().Unix() > fp.modtime[e.Name] {
		return fi.ModTime().Unix(), true
	}
	return 0, false
}

func (fp *FileProvider) trigger(e fsnotify.Event, t int64, fn []ProviderCallback) {
	data, err := os.ReadFile(e.Name)
	if err != nil {
		return
	}
	fp.mu.Lock()
	path := fp.cache[filepath.Clean(e.Name)]
	fp.modtime[e.Name] = t
	fp.mu.Unlock()
	for _, f := range fn {
		go f(path, data)
	}
}
