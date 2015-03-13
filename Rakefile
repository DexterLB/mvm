require 'rubygems'
require 'rake'
require 'rspec/core/rake_task'
require 'rubocop/rake_task'

RSpec::Core::RakeTask.new(:spec)
RuboCop::RakeTask.new(:rubocop)

namespace :version do
  VERSION_FILE = './lib/mvm/version.rb'
  require VERSION_FILE

  def update_version(new_version)
    File.write(VERSION_FILE, File.read(VERSION_FILE).gsub(
      /(^\s*VERSION\s*=\s*')(.*)('\s*)/,
      '\1' + new_version + '\3'
    ))
    system("git add #{VERSION_FILE}")
    system("git commit -m 'bump version to #{new_version}'")
  end

  def bump_version(version, id)
    patches = version.split('.').map(&:to_i)
    patches = patches.each.with_index.map do |patch, index|
      if index == id
        patch + 1
      elsif index > id
        0
      else
        patch
      end
    end
    update_version(patches.join('.'))
  end

  task :show do
    puts Mvm::VERSION
  end

  task :set do
    new_version = ENV['version']
    fail 'invalid version' unless /\d+\.\d+\.\d+/.match(new_version)

    update_version(new_version)
  end

  task :bump do
    bump_version(Mvm::VERSION, 2)
  end

  task :bump_minor do
    bump_version(Mvm::VERSION, 1)
  end

  task :bump_major do
    bump_version(Mvm::VERSION, 0)
  end
end

task default: [:spec, :rubocop]

# vim: set shiftwidth=2:
