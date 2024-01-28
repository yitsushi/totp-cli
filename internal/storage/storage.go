package storage

// Storage is an interface that describes a storage backend.
type Storage interface {
	SetPassword(password string)
	Prepare() error
	Decrypt() error

	FindNamespace(namespace string) (*Namespace, error)
	DeleteNamespace(namespace *Namespace)
	AddNamespace(namespace *Namespace) (*Namespace, error)
	ListNamespaces() []*Namespace

	Save() error
}
