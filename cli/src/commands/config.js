const fs = require('fs');
const path = require('path');
const os = require('os');

class Config {
    constructor() {
        this.configDir = this.getConfigDir();
        this.configFile = path.join(this.configDir, 'cli-config.json');
        this.load();
    }

    getConfigDir() {
        switch (process.platform) {
            case 'win32':
                return path.join(os.homedir(), 'AppData', 'Roaming', 'localp2p');
            case 'darwin':
                return path.join(os.homedir(), 'Library', 'Application Support', 'localp2p');
            default:
                return path.join(os.homedir(), '.config', 'localp2p');
        }
    }

    load() {
        try {
            if (fs.existsSync(this.configFile)) {
                const data = fs.readFileSync(this.configFile, 'utf8');
                this.data = JSON.parse(data);
            } else {
                this.data = this.getDefaults();
                this.save();
            }
        } catch (error) {
            this.data = this.getDefaults();
        }
    }

    save() {
        try {
            if (!fs.existsSync(this.configDir)) {
                fs.mkdirSync(this.configDir, { recursive: true });
            }
            fs.writeFileSync(this.configFile, JSON.stringify(this.data, null, 2));
        } catch (error) {
            console.error('Failed to save config:', error.message);
        }
    }

    getDefaults() {
        return {
            rpcPort: 9090,
            corePort: 8080,
            coreBinary: './core/localp2p',
            autoStart: true
        };
    }

    get(key) {
        return this.data[key];
    }

    set(key, value) {
        this.data[key] = value;
        this.save();
    }
}

module.exports = new Config();