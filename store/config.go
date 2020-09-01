package store

type HostParams struct {
    Host string
    Port int
}
type HostsParams []HostParams

// CacheStoreConfig is a struct for store configuration variables
type CacheStoreConfig struct {
    Hosts HostsParams
    Port int
    MaxIdleConnections int
    ReadWriteTimeOut int
    Prefix string
}