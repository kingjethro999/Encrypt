<?php

namespace Encrypt\Commands;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'init',
    description: 'Create .encrypt vault'
)]
class InitCommand extends BaseCommand
{
    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $this->showSpinner('Initializing vault...', function () {
            $this->vault->init();
        });

        $this->showInfo('📁 Created .encrypt directory with secure configuration.');

        return Command::SUCCESS;
    }
}
