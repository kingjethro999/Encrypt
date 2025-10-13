<?php

namespace Encrypt;

use Encrypt\Commands\InitCommand;
use Encrypt\Commands\LockupCommand;
use Encrypt\Commands\SetupCommand;
use Encrypt\Commands\SetCommand;
use Encrypt\Commands\GetCommand;
use Encrypt\Commands\UnlockCommand;
use Encrypt\Commands\StatusCommand;
use Encrypt\Commands\ResetCommand;
use Symfony\Component\Console\Application as BaseApplication;

class Application extends BaseApplication
{
    public function __construct()
    {
        parent::__construct('encrypt', '1.0.0');
        
        $this->setDescription('A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev.');
        
        // Add commands
        $this->add(new InitCommand());
        $this->add(new LockupCommand());
        $this->add(new SetupCommand());
        $this->add(new SetCommand());
        $this->add(new GetCommand());
        $this->add(new UnlockCommand());
        $this->add(new StatusCommand());
        $this->add(new ResetCommand());
    }
}
