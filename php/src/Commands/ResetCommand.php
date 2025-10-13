<?php

namespace Encrypt\Commands;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Question\ConfirmationQuestion;

#[AsCommand(
    name: 'reset',
    description: 'Remove vault (careful!)'
)]
class ResetCommand extends BaseCommand
{
    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $helper = $this->getHelper('question');
        $question = new ConfirmationQuestion(
            'Are you sure you want to reset the vault? This will delete all encrypted secrets. (y/N): ',
            false
        );

        if (!$helper->ask($input, $output, $question)) {
            $this->showWarning('Reset cancelled.');
            return Command::SUCCESS;
        }

        $this->showSpinner('Resetting vault...', function () {
            $this->vault->reset();
        });

        $this->showSuccess('Vault reset successfully!');

        return Command::SUCCESS;
    }
}
