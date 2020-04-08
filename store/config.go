package store

// CacheStoreConfig is a struct for store configuration variables
type CacheStoreConfig struct {
    Host string
    Port int
    MaxIdleConnections int
    ReadWriteTimeOut int
    Prefix string
}