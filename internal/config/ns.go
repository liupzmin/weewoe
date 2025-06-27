// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package config

const (
	// MaxFavoritesNS number # favorite namespaces to keep in the configuration.
	MaxFavoritesNS = 9
	defaultNS      = "default"
	allNS          = "all"
)

// Namespace tracks active and favorites namespaces.
type Namespace struct {
	Active    string   `yaml:"active"`
	Favorites []string `yaml:"favorites"`
}

// NewNamespace create a new namespace configuration.
func NewNamespace() *Namespace {
	return &Namespace{
		Active:    defaultNS,
		Favorites: []string{defaultNS},
	}
}

// SetActive set the active namespace.
func (n *Namespace) SetActive(ns string) error {
	n.Active = ns
	if ns != "" {
		n.addFavNS(ns)
	}
	return nil
}

func (n *Namespace) isAllNamespaces() bool {
	return n.Active == allNS || n.Active == ""
}

func (n *Namespace) addFavNS(ns string) {
	if InList(n.Favorites, ns) {
		return
	}

	nfv := make([]string, 0, MaxFavoritesNS)
	nfv = append(nfv, ns)
	for i := 0; i < len(n.Favorites); i++ {
		if i+1 < MaxFavoritesNS {
			nfv = append(nfv, n.Favorites[i])
		}
	}
	n.Favorites = nfv
}

func (n *Namespace) rmFavNS(ns string) {
	victim := -1
	for i, f := range n.Favorites {
		if f == ns {
			victim = i
			break
		}
	}
	if victim < 0 {
		return
	}

	n.Favorites = append(n.Favorites[:victim], n.Favorites[victim+1:]...)
}
