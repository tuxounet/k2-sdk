package bases

func (b *BaseControllerPlaybook) Init() error {

	provider := b.getComputePlaybooksProvider()

	definition := b.GetDefinition()
	err := provider.RegisterDefinition(*definition)

	if err != nil {
		b.GetLogger().ErrorF("Failed to Init container defintion inside provider: %s", err)
		return err
	}

	return nil
}
