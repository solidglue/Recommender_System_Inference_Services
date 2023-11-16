package flags

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

// singleton instance
func init() {
	flagCacheInstance = new(flagCache)
}

func getFlagCacheInstance() *flagCache {
	return flagCacheInstance
}

// bigcahe_shards
func (s *flagCache) setBigcacheShards(bigcacheShards *int) {
	s.bigcacheShards = bigcacheShards
}

func (s *flagCache) GetBigcacheShards() *int {
	return s.bigcacheShards
}

// bigcahe_lifeWindowS
func (s *flagCache) setBigcacheLifeWindowS(bigcacheLifeWindowS *int) {
	s.bigcacheLifeWindowS = bigcacheLifeWindowS
}

func (s *flagCache) GetBigcacheLifeWindowS() *int {
	return s.bigcacheLifeWindowS
}

// bigcache_cleanWindowS
func (s *flagCache) setBigcacheCleanWindowS(bigcacheCleanWindowS *int) {
	s.bigcacheCleanWindowS = bigcacheCleanWindowS
}

func (s *flagCache) GetBigcacheCleanWindowS() *int {
	return s.bigcacheCleanWindowS
}

// bigcache_hardMaxCacheSize
func (s *flagCache) setBigcacheHardMaxCacheSize(bigcacheHardMaxCacheSize *int) {
	s.bigcacheHardMaxCacheSize = bigcacheHardMaxCacheSize
}

func (s *flagCache) GetBigcacheHardMaxCacheSize() *int {
	return s.bigcacheHardMaxCacheSize
}

// bigcache_maxEntrySize
func (s *flagCache) setBigcacheMaxEntrySize(bigcacheMaxEntrySize *int) {
	s.bigcacheMaxEntrySize = bigcacheMaxEntrySize
}

func (s *flagCache) GetBigcacheMaxEntrySize() *int {
	return s.bigcacheMaxEntrySize
}

// bigcache_maxEntriesInWindow
func (s *flagCache) setBigcacheMaxEntriesInWindow(bigcacheMaxEntriesInWindow *int) {
	s.bigcacheMaxEntriesInWindow = bigcacheMaxEntriesInWindow
}

func (s *flagCache) GetBigcacheMaxEntriesInWindow() *int {
	return s.bigcacheMaxEntriesInWindow
}

// bigcache_verbose
func (s *flagCache) setBigcacheVerbose(bigcacheVerbose *bool) {
	s.bigcacheVerbose = bigcacheVerbose
}

func (s *flagCache) GetBigcacheVerbose() *bool {
	return s.bigcacheVerbose
}
