<?php

namespace Encrypt\Commands;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'set',
    description: 'Add/update a key'
)]
class SetCommand extends BaseCommand
{
    protected function configure(): void
    {
        $this->addArgument('key_value', InputArgument::REQUIRED, 'Key=value pair to set');
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $keyValue = $input->getArgument('key_value');

        if (!str_contains($keyValue, '=')) {
            $this->showError('Invalid format. Use: encrypt set KEY=value');
        }

        [$key, $value] = explode('=', $keyValue, 2);

        if (!$this->vault->exists()) {
            $this->showError('Vault not found. Run \'encrypt init\' first.');
        }

        if (!$this->vault->unlocked()) {
            $this->showError('Vault is locked. Run \'encrypt setup <password>\' to unlock secrets.');
        }

        $this->showSpinner("Setting secret: $key", function () use ($key, $value) {
            $this->vault->set($key, $value);
        });

        $this->showSuccess("Secret '$key' set successfully!");

        return Command::SUCCESS;
    }
}
