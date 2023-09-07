package resources

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"cardboard.package-operator.run/internal"
	"github.com/fsnotify/fsnotify"
	"github.com/go-logr/logr"
)

const folderWatchDebounce = 50 * time.Millisecond

type folder struct {
	internal.Named
	Path string

	watchingPaths map[string]struct{}
}

func Folder(name, path string) *folder {
	return &folder{
		Named: internal.Named(name),
		Path:  path,
	}
}

func (f *folder) Get(ctx context.Context) (Volume, error) {
	return Volume{Path: f.Path}, nil
}

func (f *folder) Put(ctx context.Context, v Volume) error {
	return nil
}

func (f *folder) Watch(ctx context.Context) (
	<-chan WatchEvent, error,
) {
	log := logr.FromContextOrDiscard(ctx).
		WithName("Folder " + f.Path)
	ctx = logr.NewContext(ctx, log)

	out := make(chan WatchEvent)
	go f.watch(ctx, out)
	return out, nil
}

func (f *folder) watch(
	ctx context.Context,
	out chan WatchEvent,
) {
	defer close(out)

	f.watchingPaths = map[string]struct{}{}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		out <- WatchEvent{Err: err}
		return
	}
	defer watcher.Close()

	if err := f.refreshWatch(ctx, watcher); err != nil {
		out <- WatchEvent{Err: err}
		return
	}

	log := logr.FromContextOrDiscard(ctx)
	for {
		select {
		case _, ok := <-internal.DebounceCh(
			watcher.Events, folderWatchDebounce):
			if !ok {
				return
			}
			log.V(internal.Debug).Info("change!")
			if err := f.refreshWatch(ctx, watcher); err != nil {
				out <- WatchEvent{Err: err}
				return
			}
			out <- WatchEvent{}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			out <- WatchEvent{Err: err}
			return

		case <-ctx.Done():
			return
		}
	}
}

func (f *folder) refreshWatch(
	ctx context.Context,
	watcher *fsnotify.Watcher,
) error {
	log := logr.FromContextOrDiscard(ctx)

	watching := map[string]struct{}{}

	if err := filepath.WalkDir(
		f.Path,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				return nil
			}

			if strings.HasPrefix(filepath.Base(path), ".") && d.Name() != "." {
				return filepath.SkipDir
			}

			path = filepath.Clean(path)
			watching[path] = struct{}{}
			if _, ok := f.watchingPaths[path]; ok {
				return nil
			}

			log.V(internal.Debug).Info("start watch", "path", path)
			return watcher.Add(path)
		}); err != nil {
		return err
	}

	for path := range f.watchingPaths {
		if _, ok := watching[path]; ok {
			continue
		}

		log.V(internal.Debug).Info("stop watch", "path", path)
		if err := watcher.Remove(path); err != nil && !errors.Is(err, fsnotify.ErrNonExistentWatch) {
			return err
		}
	}
	f.watchingPaths = watching
	return nil
}
