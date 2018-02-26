package rotation

type RetainByCount struct {
	MaxArchivesCount int
}

func (strategy *RetainByCount) PurgeSet(archives []Archive) []Archive {
	if len(archives) <= strategy.MaxArchivesCount {
		return nil
	}
	purgeCount := len(archives) - strategy.MaxArchivesCount
	return archives[:purgeCount]
}