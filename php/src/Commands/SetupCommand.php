<?php

namespace Encrypt\Commands;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'setup',
    description: 'Set up secrets on a new machine'
)]
class SetupCommand extends BaseCommand
{
    protected function configure(): void
    {
        $this->addArgument('password', InputArgument::REQUIRED, 'Password to unlock secrets');
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $password = $input->getArgument('password');

        if (!$this->vault->exists()) {
            $this->showError('Vault not found. Run \'encrypt init\' first.');
        }

        $this->showSpinner('Setting up vault...', function () use ($password) {
            $this->vault->setup($password);
        });

        $this->showInfo('🔓 Your secrets are now available for use.');

        return Command::SUCCESS;
    }
}
