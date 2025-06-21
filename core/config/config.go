package config

import (
    "encoding/json"
    "os"
    "path/filepath"
    "runtime"
)

type Config struct {
    NodeID       string `json:"node_id"`
    DisplayName  string `json:"display_name"`
    Port         int    `json:"port"`
    ServiceName  string `json:"service_name"`
    Domain       string `json:"domain"`
    DataDir      string `json:"data_dir"`
}

func getConfigDir() string {
    var configDir string
    switch runtime.GOOS {
    case "windows":
        configDir = filepath.Join(os.Getenv("APPDATA"), "localp2p")
    case "darwin":
        configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "localp2p")
    default:
        configDir = filepath.Join(os.Getenv("HOME"), ".config", "localp2p")
    }
    return configDir
}

func LoadConfig() (*Config, error) {
    configDir := getConfigDir()
    configFile := filepath.Join(configDir, "config.json")
    
    // Create config directory if it doesn't exist
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return nil, err
    }
    
    // Check if config file exists
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        // Create default config
        config := &Config{
            NodeID:      generateNodeID(),
            DisplayName: getHostname(),
            Port:        8080,
            ServiceName: "_localp2p._tcp",
            Domain:      "local.",
            DataDir:     configDir,
        }
        
        if err := config.Save(); err != nil {
            return nil, err
        }
        return config, nil
    }
    
    // Load existing config
    data, err := os.ReadFile(configFile)
    if err != nil {
        return nil, err
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}

func (c *Config) Save() error {
    configDir := getConfigDir()
    configFile := filepath.Join(configDir, "config.json")
    
    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(configFile, data, 0644)
}

func generateNodeID() string {
    // Simple node ID generation - in production, use proper UUID
    hostname, _ := os.Hostname()
    return hostname + "-" + randomString(8)
}

func getHostname() string {
    hostname, err := os.Hostname()
    if err != nil {
        return "unknown"
    }
    return hostname
}

func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[len(charset)-1] // Simple random for demo
    }
    return string(b)
}