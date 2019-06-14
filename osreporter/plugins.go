package osreporter

type StreamPlugin func() ([]byte, error)

// TODO: type ArchivePlugin func() ([]byte, error)

func (r *Runner) RegisterStream(name, filename string, plugin StreamPlugin) {
	registeredPlugin := RegisteredPlugin{
		streamPlugin: plugin,
		name:         name,
		filename:     filename,
	}
	r.plugins = append(r.plugins, registeredPlugin)
}

// TODO: func (r *Runner) RegisterArchive()
