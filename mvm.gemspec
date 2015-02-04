require File.dirname(__FILE__) + '/git-version-bump/lib/git-version-bump'

Gem::Specification.new do |s|
  s.name = 'mvm'
  s.version = GVB.version
  s.date = GVB.date
  s.authors = ['Angel Angelov']
  s.email = ['hextwoa at gmail.com']
  s.homepage = 'http://github.com/DexterLB/mvm'
  s.summary = 'Movie identifier, renamer and lister'
  s.description = 'Provides an interface for identifying and renaming movies.'
  s.files = `git ls-files`.split("\n")
  s.test_files = `git ls-files -- {test,spec,features}/*`.split("\n")
  s.executables = `git ls-files -- bin/*`.split("\n")
    .map { |f| File.basename(f) }
  s.require_paths = ['lib']
  s.license = 'GPLv3'
  s.add_dependency 'git-version-bump', '~> 0.10'
  s.add_runtime_dependency 'xdg', '~> 2'
  s.add_runtime_dependency 'iso-639', '~> 0'
  s.add_runtime_dependency 'streamio-ffmpeg', '~> 1'
  s.add_runtime_dependency 'colorize', '~> 0'
  s.add_runtime_dependency 'parallel', '~> 1'
  s.add_development_dependency 'rspec', '~> 3'
  s.add_development_dependency 'rubocop', '~> 0'
  s.add_development_dependency 'vcr', '~> 2'
  s.add_development_dependency 'webmock', '~> 1'
  s.add_development_dependency 'fakefs', '~> 0'
end

# vim: set shiftwidth=2:
