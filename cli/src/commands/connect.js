const axios = require('axios');
const chalk = require('chalk');
const ora = require('ora');
const inquirer = require('inquirer');
const config = require('../config');

async function connect(options) {
    let address = options.address;
    let port = options.port;

    // If no address provided, discover and let user choose
    if (!address) {
        const spinner = ora('Discovering peers...').start();
        
        try {
            const response = await axios.post(`http://localhost:${config.get('rpcPort')}/rpc`, {
                method: 'discover',
                params: {}
            });

            spinner.stop();

            if (response.data.error) {
                console.error(chalk.red('Error:'), response.data.error);
                return;
            }

            const peers = response.data.result || [];
            
            if (peers.length === 0) {
                console.log(chalk.yellow('No peers discovered'));
                return;
            }

            const choices = peers.map(peer => ({
                name: `${peer.Name} (${peer.Address}:${peer.Port})`,
                value: { address: peer.Address, port: peer.Port }
            }));

            const { selectedPeer } = await inquirer.prompt([{
                type: 'list',
                name: 'selectedPeer',
                message: 'Select a peer to connect to:',
                choices
            }]);

            address = selectedPeer.address;
            port = selectedPeer.port;

        } catch (error) {
            spinner.stop();
            console.error(chalk.red('Failed to discover peers:'), error.message);
            return;
        }
    }

    // Connect to the selected peer
    const spinner = ora(`Connecting to ${address}:${port}...`).start();
    
    try {
        const response = await axios.post(`http://localhost:${config.get('rpcPort')}/rpc`, {
            method: 'connect',
            params: { address, port }
        });

        spinner.stop();

        if (response.data.error) {
            console.error(chalk.red('Connection failed:'), response.data.error);
            return;
        }

        console.log(chalk.green(`✓ Connected to ${address}:${port}`));

    } catch (error) {
        spinner.stop();
        console.error(chalk.red('Failed to connect:'), error.message);
    }
}

module.exports = connect;
