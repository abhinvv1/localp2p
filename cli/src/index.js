#!/usr/bin/env node

const { program } = require('commander');
const chalk = require('chalk');
const { spawn } = require('child_process');
const path = require('path');

const config = require('./config');
const discover = require('./commands/discover');
const connect = require('./commands/connect');
const send = require('./commands/send');

program
    .name('localp2p')
    .description('LocalP2P - Secure local network communication')
    .version('1.0.0');

program
    .command('start')
    .description('Start the LocalP2P core daemon')
    .option('-p, --port ', 'Core service port', config.get('corePort'))
    .option('-r, --rpc-port ', 'RPC server port', config.get('rpcPort'))
    .action((options) => {
        console.log(chalk.blue('Starting LocalP2P core...'));
        
        const coreBinary = path.resolve(__dirname, '../../core/localp2p');
        const args = [`--rpc-port=${options.rpcPort}`];
        
        const core = spawn(coreBinary, args, {
            stdio: 'inherit',
            detached: false
        });
        
        core.on('error', (error) => {
            console.error(chalk.red('Failed to start core:'), error.message);
            console.log(chalk.yellow('Make sure the core binary is built:'));
            console.log(chalk.gray('cd core && go build -o localp2p'));
        });

        core.on('exit', (code) => {
            console.log(chalk.yellow(`Core exited with code ${code}`));
        });

        // Handle shutdown
        process.on('SIGINT', () => {
            console.log(chalk.yellow('\nShutting down...'));
            core.kill('SIGTERM');
        });
    });

program
    .command('discover')
    .description('Discover peers on the local network')
    .action(discover);

program
    .command('connect')
    .description('Connect to a peer')
    .option('-a, --address ', 'Peer IP address')
    .option('-p, --port ', 'Peer port', parseInt)
    .action(connect);

program
    .command('send')
    .description('Send a message to a connected peer')
    .argument('[message]', 'Message to send')
    .option('-t, --to ', 'Target peer ID')
    .action(send);

program
    .command('status')
    .description('Show current status and connections')
    .action(async () => {
        try {
            const response = await require('axios').post(`http://localhost:${config.get('rpcPort')}/rpc`, {
                method: 'connections',
                params: {}
            });

            if (response.data.error) {
                console.error(chalk.red('Error:'), response.data.error);
                return;
            }

            const connections = response.data.result || [];
            console.log(chalk.green(`Active connections: ${connections.length}`));
            
            if (connections.length > 0) {
                console.log(chalk.gray('─'.repeat(30)));
                connections.forEach((peer, index) => {
                    console.log(`${index + 1}. ${chalk.cyan(peer)}`);
                });
            }

        } catch (error) {
            console.error(chalk.red('Core is not running'));
            console.log(chalk.gray('Use "localp2p start" to start the core'));
        }
    });

// Parse command line arguments
program.parse();

// Show help if no command provided
if (!process.argv.slice(2).length) {
    program.outputHelp();
    console.log();
    console.log(chalk.gray('Quick start:'));
    console.log(chalk.gray('  1. localp2p start    # Start the core service'));
    console.log(chalk.gray('  2. localp2p discover # Find peers'));
    console.log(chalk.gray('  3. localp2p connect  # Connect to a peer'));
    console.log(chalk.gray('  4. localp2p send     # Send a message'));
}