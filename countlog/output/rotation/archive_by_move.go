package rotation

import "os"

type ArchiveByMove struct {
	NamingStrategy NamingStrategy
}

func (strategy *ArchiveByMove) Archive(path string) ([]Archive, error) {
	newPath, err := strategy.NamingStrategy.NextFile()
	if err != nil {
		return nil, err
	}
	err = os.Rename(path, newPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return strategy.NamingStrategy.ListFiles()
}
