package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/nkapatos/sbx-kit/cli/internal/binding"
	"github.com/nkapatos/sbx-kit/cli/internal/catalog"
	"github.com/nkapatos/sbx-kit/cli/internal/resources"
	"github.com/nkapatos/sbx-kit/cli/internal/sbxname"
	"github.com/nkapatos/sbx-kit/cli/internal/xdg"
)

type resolvedSandbox struct {
	AgentName   string
	SbxAgent    string
	ProjectDir  string
	SandboxName string
	ProfileID   string
	KitPaths    []string
	ImageName   string
	TemplateFB  string
	Resources   *resources.Profile
	ResProfile  string
	Root        string
}

func resolveFromAgent(agentName, projectDir string) (*resolvedSandbox, error) {
	if err := xdg.Ensure(); err != nil {
		return nil, err
	}
	root, err := requireToolkitRoot()
	if err != nil {
		return nil, err
	}
	cat, err := catalog.Load(filepath.Join(root, "config", "agents.yaml"))
	if err != nil {
		return nil, err
	}
	agent, ok := cat.Agents[agentName]
	if !ok {
		return nil, fmt.Errorf("unknown recipe %q (try: sbx-kit recipes)", agentName)
	}
	if agent.Stub {
		return nil, fmt.Errorf("recipe %q is still a stub in config/agents.yaml", agentName)
	}
	abs, err := filepath.Abs(projectDir)
	if err != nil {
		return nil, err
	}

	name := ""
	profileID := sbxname.NewProfileID(agentName, abs)
	if rec, err := binding.Get(abs, agentName); err == nil && rec != nil {
		name = rec.SandboxName
		profileID = rec.ProfileID
	}
	if name == "" {
		var err error
		name, err = sbxname.FromDir(abs)
		if err != nil {
			return nil, err
		}
	}

	kits := catalog.ResolveKits(agent.Kits, cat.Defaults.Kits)
	kitPaths := catalog.KitPaths(root, kits)

	resProfile := cat.Defaults.Resources
	if resProfile == "" {
		resProfile = "remote-llm"
	}
	res, err := resources.Load(root, resProfile)
	if err != nil {
		return nil, err
	}

	return &resolvedSandbox{
		AgentName:   agentName,
		SbxAgent:    agent.SbxAgent,
		ProjectDir:  abs,
		SandboxName: name,
		ProfileID:   profileID,
		KitPaths:    kitPaths,
		ImageName:   agent.ImageName,
		TemplateFB:  agent.TemplateFallback,
		Resources:   res,
		ResProfile:  resProfile,
		Root:        root,
	}, nil
}

func resolveSandboxArg(arg, projectDir string) (*resolvedSandbox, error) {
	// Prefer catalog recipe id when it matches.
	root, err := requireToolkitRoot()
	if err != nil {
		return nil, err
	}
	cat, err := catalog.Load(filepath.Join(root, "config", "agents.yaml"))
	if err != nil {
		return nil, err
	}
	if _, ok := cat.Agents[arg]; ok {
		return resolveFromAgent(arg, projectDir)
	}
	// Treat as sandbox name.
	if err := xdg.Ensure(); err != nil {
		return nil, err
	}
	rec, err := binding.GetBySandbox(arg)
	if err != nil {
		return nil, err
	}
	if rec != nil {
		full, err := resolveFromAgent(rec.Agent, rec.ProjectDir)
		if err != nil {
			// Binding exists but recipe missing from catalog — still useful for rm/check.
			return &resolvedSandbox{
				AgentName:   rec.Agent,
				ProjectDir:  rec.ProjectDir,
				SandboxName: rec.SandboxName,
				ProfileID:   rec.ProfileID,
				Root:        root,
			}, nil
		}
		full.SandboxName = rec.SandboxName
		full.ProfileID = rec.ProfileID
		return full, nil
	}
	return &resolvedSandbox{
		SandboxName: arg,
		ProfileID:   arg,
		ProjectDir:  projectDir,
		Root:        root,
	}, nil
}
