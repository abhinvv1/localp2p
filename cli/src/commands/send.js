const axios = require('axios');
const chalk = require('chalk');
const ora = require('ora');
const inquirer = require('inquirer');
const config = require('../config');

async function send(message, options) {
    let targetPeer = options.to;

    // If no target specified, show connected peers
    if (!targetPeer) {
        const spinner = ora('Getting connected peers...').start();
        
        try {
            const response = await axios.post(`http://localhost:${config.get('rpcPort')}/rpc`, {
                method: 'connections',
                params: {}
            });

            spinner.stop();

            if (response.data.error) {
                console.error(chalk.red('Error:'), response.data.error);
                return;
            }

            const connections = response.data.result || [];
            
            if (connections.length === 0) {
                console.log(chalk.yellow('No active connections'));
                console.log(chalk.gray('Use "localp2p connect" to connect to a peer first'));
                return;
            }

            const { selectedPeer } = await inquirer.prompt([{
                type: 'list',
                name: 'selectedPeer',
                message: 'Select a peer to send message to:',
                choices: connections
            }]);

            targetPeer = selectedPeer;

        } catch (error) {
            spinner.stop();
            console.error(chalk.red('Failed to get connections:'), error.message);
            return;
        }
    }

    // If no message provided, prompt for it
    if (!message) {
        const { inputMessage } = await inquirer.prompt([{
            type: 'input',
            name: 'inputMessage',
            message: 'Enter your message:',
            validate: input => input.trim().length > 0 || 'Message cannot be empty'
        }]);
        
        message = inputMessage;
    }

    // Send the message
    const spinner = ora(`Sending message to ${targetPeer}...`).start();
    
    try {
        const response = await axios.post(`http://localhost:${config.get('rpcPort')}/rpc`, {
            method: 'send',
            params: { to: targetPeer, content: message }
        });

        spinner.stop();

        if (response.data.error) {
            console.error(chalk.red('Failed to send message:'), response.data.error);
            return;
        }

        console.log(chalk.green(`✓ Message sent to ${targetPeer}`));
        console.log(chalk.gray(`"${message}"`));

    } catch (error) {
        spinner.stop();
        console.error(chalk.red('Failed to send message:'), error.message);
    }
}

module.exports = send;