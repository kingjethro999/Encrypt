<?php

namespace Encrypt\Commands;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'unlock',
    description: 'Decrypt everything into .env'
)]
class UnlockCommand extends BaseCommand
{
    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        if (!$this->vault->exists()) {
            $this->showError('Vault not found. Run \'encrypt init\' first.');
        }

        if (!$this->vault->unlocked()) {
            $this->showError('Vault is locked. Run \'encrypt setup <password>\' to unlock secrets.');
        }

        $this->showSpinner('Unlocking secrets...', function () {
            $secrets = $this->vault->all();
            $envContent = '';
            foreach ($secrets as $key => $value) {
                $envContent .= "$key=$value\n";
            }
            file_put_contents('.env', $envContent);
        });

        $this->showSuccess('Secrets unlocked and written to .env');

        return Command::SUCCESS;
    }
}
