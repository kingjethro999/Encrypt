Gem::Specification.new do |spec|
  spec.name          = "encrypt"
  spec.version       = "1.0.0"
  spec.authors       = ["Encrypt Team"]
  spec.email         = [""]

  spec.summary       = "A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev."
  spec.description   = "Encrypt replaces .env files with an encrypted local secrets vault that provides triple-layer encryption, seamless team onboarding, and production-ready deployment."
  spec.homepage      = ""
  spec.license       = "MIT"

  spec.files         = Dir["lib/**/*", "bin/**/*", "*.md", "*.txt"]
  spec.bindir        = "exe"
  spec.executables   = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }
  spec.require_paths = ["lib"]

  spec.required_ruby_version = ">= 2.7.0"

  spec.add_dependency "thor", "~> 1.2"
  spec.add_dependency "colorize", "~> 0.8"
  spec.add_dependency "tty-spinner", "~> 0.9"

  spec.add_development_dependency "bundler", "~> 2.0"
  spec.add_development_dependency "rake", "~> 13.0"
  spec.add_development_dependency "rspec", "~> 3.0"
  spec.add_development_dependency "pry", "~> 0.14"
end
