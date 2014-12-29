Gem::Specification.new do |s|
  s.name = "mvm"
  s.version = 0.0.1
  s.authors = ["Angel Angelov"]
  s.email = ["hextwoa at gmail.com"]
  s.homepage = %q{http://github.com/DexterLB/mvm}
  s.summary = %q{Movie identifier, renamer and lister}
  s.files = `git ls-files`.split("\n")
  s.test_files = `git ls-files -- {test,spec,features}/*`.split("\n")
  s.executables = `git ls-files -- bin/*`.split("\n").map{ |f| File.basename(f) }
  s.require_paths = ["lib"]
  s.license = 'GPLv3'
  s.add_development_dependency 'rspec'
end
