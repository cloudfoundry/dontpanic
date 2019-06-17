package osreporter

import (
	"context"
	"time"
)

type StreamPlugin func(context.Context) ([]byte, error)

// TODO: type ArchivePlugin func() ([]byte, error)

func (r *Runner) RegisterStream(name, filename string, plugin StreamPlugin, timeout ...time.Duration) *RegisteredPlugin {
	maxDuration := 10 * time.Second
	if len(timeout) > 0 {
		maxDuration = timeout[0]
	}

	registeredPlugin := RegisteredPlugin{
		streamPlugin: plugin,
		name:         name,
		filename:     filename,
		timeout:      maxDuration,
	}
	r.plugins = append(r.plugins, &registeredPlugin)
	return &registeredPlugin
}

func (r *Runner) RegisterEchoStream(name, filename string, plugin StreamPlugin, timeout ...time.Duration) *RegisteredPlugin {
	registeredPlugin := r.RegisterStream(name, filename, plugin, timeout...)
	registeredPlugin.echoOutput = true
	return registeredPlugin
}

// TODO: func (r *Runner) RegisterArchive()
