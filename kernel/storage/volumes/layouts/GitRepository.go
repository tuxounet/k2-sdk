package layouts

import "github.com/tuxounet/k2-sdk/kernel/storage/volumes/bases"

type GitRepository struct {
	bases.BaseLayout
}

func NewGitRepositoryLayout() *GitRepository {
	base := bases.NewBaseLayout()
	return &GitRepository{base}
}
