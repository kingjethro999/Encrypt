<?php

namespace Encrypt\Commands;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'get',
    description: 'Fetch decrypted value'
)]
class GetCommand extends BaseCommand
{
    protected function configure(): void
    {
        $this->addArgument('key', InputArgument::REQUIRED, 'Key to fetch');
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $key = $input->getArgument('key');

        if (!$this->vault->exists()) {
            $this->showError('Vault not found. Run \'encrypt init\' first.');
        }

        if (!$this->vault->unlocked()) {
            $this->showError('Vault is locked. Run \'encrypt setup <password>\' to unlock secrets.');
        }

        try {
            $value = $this->vault->get($key);
            $output->writeln($value);
        } catch (\Exception $e) {
            $this->showError($e->getMessage());
        }

        return Command::SUCCESS;
    }
}
