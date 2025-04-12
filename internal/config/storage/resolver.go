package storage

import "strings"

type Resolver struct {
	prefixes map[string]Storage
	def      Storage
}

func NewResolver(
	def Storage,
	prefixes map[string]Storage,
) *Resolver {
	return &Resolver{
		prefixes: prefixes,
		def:      def,
	}
}

func (r *Resolver) Resolve(path string) Storage {
	for prefix, storage := range r.prefixes {
		if strings.Contains(path, prefix) {
			return storage
		}
	}

	return r.def
}
