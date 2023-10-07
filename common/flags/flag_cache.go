package flags

import "flag"

var flagCacheInstance *flagCache

type flagCache struct {

	//cache
	bigcacheShards             *int
	bigcacheLifeWindowS        *int
	bigcacheCleanWindowS       *int
	bigcacheHardMaxCacheSize   *int
	bigcacheMaxEntrySize       *int
	bigcacheMaxEntriesInWindow *int
	bigcacheVerbose            *bool
}

//singleton instance
func init() {
	flagCacheInstance = new(flagCache)
}

func getFlagCacheInstance() *flagCache {
	return flagCacheInstance
}

// bigcahe_shards
func (s *flagCache) setBigcacheShards() {
	conf := flag.Int("bigcahe_shards", 1024, "")
	s.bigcacheShards = conf
}

func (s *flagCache) GetBigcacheShards() *int {
	return s.bigcacheShards
}

// bigcahe_lifeWindowS
func (s *flagCache) setBigcacheLifeWindowS() {
	conf := flag.Int("bigcahe_lifeWindowS", 300, "")
	s.bigcacheLifeWindowS = conf
}

func (s *flagCache) GetBigcacheLifeWindowS() *int {
	return s.bigcacheLifeWindowS
}

// bigcache_cleanWindowS
func (s *flagCache) setBigcacheCleanWindowS() {
	conf := flag.Int("bigcache_cleanWindowS", 120, "")
	s.bigcacheCleanWindowS = conf
}

func (s *flagCache) GetBigcacheCleanWindowS() *int {
	return s.bigcacheCleanWindowS
}

// bigcache_hardMaxCacheSize
func (s *flagCache) setBigcacheHardMaxCacheSize() {
	conf := flag.Int("bigcache_hardMaxCacheSize", 409600, "MB")
	s.bigcacheHardMaxCacheSize = conf
}

func (s *flagCache) GetBigcacheHardMaxCacheSize() *int {
	return s.bigcacheHardMaxCacheSize
}

//bigcache_maxEntrySize
func (s *flagCache) setBigcacheMaxEntrySize() {
	conf := flag.Int("bigcache_maxEntrySize", 1024, "byte")
	s.bigcacheMaxEntrySize = conf
}

func (s *flagCache) GetBigcacheMaxEntrySize() *int {
	return s.bigcacheMaxEntrySize
}

// bigcache_maxEntriesInWindow
func (s *flagCache) setBigcacheMaxEntriesInWindow() {
	conf := flag.Int("bigcache_maxEntriesInWindow", 2000000, "depends on tps")
	s.bigcacheMaxEntriesInWindow = conf
}

func (s *flagCache) GetBigcacheMaxEntriesInWindow() *int {
	return s.bigcacheMaxEntriesInWindow
}

// bigcache_verbose
func (s *flagCache) setBigcacheVerbose() {
	conf := flag.Bool("bigcache_verbose", false, "")
	s.bigcacheVerbose = conf
}

func (s *flagCache) GetBigcacheVerbose() *bool {
	return s.bigcacheVerbose
}
