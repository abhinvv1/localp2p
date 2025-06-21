const axios = require('axios');
const chalk = require('chalk');
const ora = require('ora');
const config = require('../config');

async function discover(options) {
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

        console.log(chalk.green(`\nFound ${peers.length} peer(s):`));
        console.log(chalk.gray('─'.repeat(60)));
        
        peers.forEach((peer, index) => {
            console.log(`${index + 1}. ${chalk.cyan(peer.Name)} (${chalk.gray(peer.ID)})`);
            console.log(`   ${chalk.gray('Address:')} ${peer.Address}:${peer.Port}`);
            console.log(`   ${chalk.gray('Last seen:')} ${new Date(peer.LastSeen).toLocaleString()}`);
            console.log();
        });

    } catch (error) {
        spinner.stop();
        console.error(chalk.red('Failed to discover peers:'), error.message);
        
        if (error.code === 'ECONNREFUSED') {
            console.log(chalk.yellow('\nTip: Make sure the LocalP2P core is running'));
            console.log(chalk.gray('Run: localp2p start'));
        }
    }
}

module.exports = discover;