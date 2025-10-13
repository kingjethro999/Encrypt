<?php

namespace Encrypt\Commands;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'lockup',
    description: 'Encrypt and secure secrets with password'
)]
class LockupCommand extends BaseCommand
{
    protected function configure(): void
    {
        $this->addArgument('password', InputArgument::REQUIRED, 'Password to encrypt secrets');
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $password = $input->getArgument('password');

        if (!$this->vault->exists()) {
            $this->showError('Vault not found. Run \'encrypt init\' first.');
        }

        $this->showSpinner('Locking up secrets...', function () use ($password) {
            $this->vault->lockup($password);
        });

        $this->showInfo('🔒 Your secrets are now encrypted and secure.');

        return Command::SUCCESS;
    }
}
