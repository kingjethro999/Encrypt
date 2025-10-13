<?php

namespace Encrypt\Commands;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'status',
    description: 'Check if vault is locked, list keys'
)]
class StatusCommand extends BaseCommand
{
    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        if (!$this->vault->exists()) {
            $this->showWarning('Vault not found. Run \'encrypt init\' first.');
            return Command::SUCCESS;
        }

        $status = $this->vault->status();

        $this->io->writeln('📊 <info>Vault Status:</info>');
        $this->io->writeln('┌─────────────────┬──────────────────────────────┐');
        $this->io->writeln('│ Property        │ Value                        │');
        $this->io->writeln('├─────────────────┼──────────────────────────────┤');

        $statusIcon = $status->isLocked ? '🔒' : '🔓';
        $statusText = $status->isLocked ? 'Locked' : 'Unlocked';
        $statusColor = $status->isLocked ? 'red' : 'green';

        $this->io->writeln("│ Status          │ $statusIcon <$statusColor>$statusText</$statusColor>");
        $this->io->writeln('│ Keys            │ ' . str_pad(count($status->keys), 28) . ' │');

        if ($status->lastModified !== null) {
            $lastModified = date('Y-m-d\TH:i:s.v\Z', $status->lastModified);
            $this->io->writeln('│ Last Modified   │ ' . str_pad($lastModified, 28) . ' │');
        }

        $this->io->writeln('└─────────────────┴──────────────────────────────┘');

        if (!empty($status->keys)) {
            $this->io->writeln('');
            $this->io->writeln('🔑 <info>Available keys:</info>');
            foreach ($status->keys as $key) {
                $this->io->writeln("  • $key");
            }
        }

        return Command::SUCCESS;
    }
}
