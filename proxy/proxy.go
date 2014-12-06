// package proxy defines the Proxy interface which should be implemented by
// package user.
//
// It also provides Saver/Loader to do the load/save concurrently.
package proxy

// Proxy is the interface do real load/save.
// Proxy can be typically implemented to load/save data from a remote
// key-value database.
type Proxy interface {
	ProxyLoader
	ProxySaver
}
