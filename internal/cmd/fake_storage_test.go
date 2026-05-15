package cmd

import (
	"github.com/yitsushi/totp-cli/internal/storage"
)

// fakeStorage is an in-memory Storage implementation for tests.
type fakeStorage struct {
	namespaces []*storage.Namespace
	SaveCalled bool
	SaveErr    error
}

func newFakeStorage(namespaces ...*storage.Namespace) *fakeStorage {
	return &fakeStorage{namespaces: namespaces}
}

func (f *fakeStorage) SetPassword(_ string) {}

func (f *fakeStorage) Prepare() error { return nil }

func (f *fakeStorage) Decrypt() error { return nil }

func (f *fakeStorage) Save() error {
	f.SaveCalled = true

	return f.SaveErr
}

func (f *fakeStorage) ListNamespaces() []*storage.Namespace {
	return f.namespaces
}

func (f *fakeStorage) FindNamespace(name string) (*storage.Namespace, error) {
	for _, ns := range f.namespaces {
		if ns.Name == name {
			return ns, nil
		}
	}

	return nil, storage.NotFoundError{Type: "namespace", Name: name}
}

func (f *fakeStorage) AddNamespace(ns *storage.Namespace) (*storage.Namespace, error) {
	if existing, err := f.FindNamespace(ns.Name); err == nil {
		return existing, storage.BackendError{Message: "namespace already exists: " + ns.Name}
	}

	f.namespaces = append(f.namespaces, ns)

	return ns, nil
}

func (f *fakeStorage) DeleteNamespace(ns *storage.Namespace) {
	for i, item := range f.namespaces {
		if item == ns {
			f.namespaces = append(f.namespaces[:i], f.namespaces[i+1:]...)

			return
		}
	}
}
