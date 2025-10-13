<?php

namespace Encrypt\Commands;

use Encrypt\Exceptions\EncryptException;
use Encrypt\Vault\Vault;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Style\SymfonyStyle;

abstract class BaseCommand extends Command
{
    protected SymfonyStyle $io;
    protected Vault $vault;

    protected function initialize(InputInterface $input, OutputInterface $output): void
    {
        $this->io = new SymfonyStyle($input, $output);
        $this->vault = new Vault();
    }

    protected function showSpinner(string $message, callable $callback): mixed
    {
        $this->io->write("⠋ $message");
        
        try {
            $result = $callback();
            $this->io->write("\r✅ " . str_replace('⠋ ', '', $message) . " successfully!\n");
            return $result;
        } catch (EncryptException $e) {
            $this->io->write("\r❌ Failed to " . strtolower(str_replace('⠋ ', '', $message)) . "\n");
            $this->io->error($e->getMessage());
            exit(1);
        }
    }

    protected function showError(string $message): void
    {
        $this->io->error($message);
        exit(1);
    }

    protected function showSuccess(string $message): void
    {
        $this->io->success($message);
    }

    protected function showInfo(string $message): void
    {
        $this->io->info($message);
    }

    protected function showWarning(string $message): void
    {
        $this->io->warning($message);
    }
}
