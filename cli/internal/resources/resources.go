package resources

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
	Memory     string
	CPUs       string
	RootSize   string
	DockerSize string
}

func Load(toolkitRoot, profile string) (*Profile, error) {
	path := filepath.Join(toolkitRoot, "config", "resources-"+profile+".env")
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("resources profile %q: %w", profile, err)
	}
	defer f.Close()

	fileVals := map[string]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		fileVals[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	get := func(envKey, fileKey string) string {
		if v := os.Getenv(envKey); v != "" {
			return v
		}
		return fileVals[fileKey]
	}

	p := &Profile{
		Memory:     get("SBX_MEMORY", "SBX_MEMORY"),
		CPUs:       get("SBX_CPUS", "SBX_CPUS"),
		RootSize:   get("SBX_ROOT_SIZE", "SBX_ROOT_SIZE"),
		DockerSize: get("SBX_DOCKER_SIZE", "SBX_DOCKER_SIZE"),
	}
	if v := os.Getenv("DOCKER_SANDBOXES_ROOT_SIZE"); v != "" {
		p.RootSize = v
	}
	if v := os.Getenv("DOCKER_SANDBOXES_DOCKER_SIZE"); v != "" {
		p.DockerSize = v
	}
	if p.Memory == "" || p.CPUs == "" {
		return nil, fmt.Errorf("incomplete resources in %s", path)
	}
	return p, nil
}
