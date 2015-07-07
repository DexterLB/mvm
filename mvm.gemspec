# -*- encoding: utf-8 -*-
# stub: mvm 0.13.7 ruby lib

Gem::Specification.new do |s|
  s.name = "mvm"
  s.version = "0.13.7"

  s.required_rubygems_version = Gem::Requirement.new(">= 0") if s.respond_to? :required_rubygems_version=
  s.require_paths = ["lib"]
  s.authors = ["Angel Angelov"]
  s.date = "2015-07-07"
  s.description = "Provides an interface for identifying and renaming movies."
  s.email = ["hextwoa at gmail.com"]
  s.executables = ["mvm"]
  s.files = [".gitignore", ".rspec", ".rubocop.yml", "Gemfile", "README.md", "Rakefile", "VERSION", "bin/mvm", "lib/mvm.rb", "lib/mvm/library.rb", "lib/mvm/library/cli.rb", "lib/mvm/library/files.rb", "lib/mvm/library/imdb.rb", "lib/mvm/library/metadata.rb", "lib/mvm/library/opensubtitles.rb", "lib/mvm/library/rename.rb", "lib/mvm/library/store.rb", "lib/mvm/opensubtitles_client.rb", "lib/mvm/opensubtitles_client/api.rb", "lib/mvm/opensubtitles_client/error.rb", "lib/mvm/opensubtitles_client/hasher.rb", "lib/mvm/settings.rb", "lib/mvm/version.rb", "mvm.gemspec", "roadmap.md", "spec/001/mvm/library/imdb_spec.rb", "spec/001/mvm/library/metadata_spec.rb", "spec/001/mvm/library/opensubtitles_spec.rb", "spec/001/mvm/library_spec.rb", "spec/001/mvm/opensubtitles_client/api_spec.rb", "spec/001/mvm/opensubtitles_client/hasher_spec.rb", "spec/001/mvm/opensubtitles_client_spec.rb", "spec/001/mvm/settings_spec.rb", "spec/002/mvm/library/files_spec.rb", "spec/002/mvm/library/rename_spec.rb", "spec/002/mvm/library/store_spec.rb", "spec/fixtures/drop.avi", "spec/fixtures/empty.file", "spec/fixtures/http/imdb_get_data.yml", "spec/fixtures/http/imdb_get_data_for_episode.yml", "spec/fixtures/http/imdb_get_data_for_movie.yml", "spec/fixtures/http/imdb_get_data_for_wrong.yml", "spec/fixtures/http/os_call.yml", "spec/fixtures/http/os_client.yml", "spec/fixtures/http/os_client_lookup_hash_ep.yml", "spec/fixtures/http/os_client_lookup_hash_movie.yml", "spec/fixtures/http/os_client_lookup_hash_wrong.yml", "spec/fixtures/http/os_client_lookup_hashes.yml", "spec/fixtures/http/os_client_search_subtitles_hash.yml", "spec/fixtures/http/os_client_search_subtitles_imdb_id.yml", "spec/fixtures/http/os_client_search_subtitles_lang.yml", "spec/fixtures/http/os_id_by_hash_ep.yml", "spec/fixtures/http/os_id_by_hash_movie.yml", "spec/fixtures/http/os_id_by_hash_wrong.yml", "spec/fixtures/http/os_id_by_hashes.yml", "spec/fixtures/http/os_login.yml", "spec/fixtures/http/os_logout.yml", "spec/fixtures/http/os_lookup_hash_ep.yml", "spec/fixtures/http/os_lookup_hash_movie.yml", "spec/fixtures/http/os_lookup_hash_wrong.yml", "spec/fixtures/http/os_lookup_hashes.yml", "spec/fixtures/http/os_safe_call.yml", "spec/fixtures/http/os_search_subtitles.yml", "spec/fixtures/http/os_search_subtitles_for_hash.yml", "spec/fixtures/http/os_search_subtitles_for_imdb_id.yml", "spec/fixtures/http/os_search_subtitles_for_imdb_id_lang.yml", "spec/fixtures/http/os_search_subtitles_for_title.yml", "spec/spec_helper.rb"]
  s.homepage = "http://github.com/DexterLB/mvm"
  s.licenses = ["GPLv3"]
  s.rubygems_version = "2.4.5"
  s.summary = "Movie identifier, renamer and lister"
  s.test_files = ["spec/001/mvm/library/imdb_spec.rb", "spec/001/mvm/library/metadata_spec.rb", "spec/001/mvm/library/opensubtitles_spec.rb", "spec/001/mvm/library_spec.rb", "spec/001/mvm/opensubtitles_client/api_spec.rb", "spec/001/mvm/opensubtitles_client/hasher_spec.rb", "spec/001/mvm/opensubtitles_client_spec.rb", "spec/001/mvm/settings_spec.rb", "spec/002/mvm/library/files_spec.rb", "spec/002/mvm/library/rename_spec.rb", "spec/002/mvm/library/store_spec.rb", "spec/fixtures/drop.avi", "spec/fixtures/empty.file", "spec/fixtures/http/imdb_get_data.yml", "spec/fixtures/http/imdb_get_data_for_episode.yml", "spec/fixtures/http/imdb_get_data_for_movie.yml", "spec/fixtures/http/imdb_get_data_for_wrong.yml", "spec/fixtures/http/os_call.yml", "spec/fixtures/http/os_client.yml", "spec/fixtures/http/os_client_lookup_hash_ep.yml", "spec/fixtures/http/os_client_lookup_hash_movie.yml", "spec/fixtures/http/os_client_lookup_hash_wrong.yml", "spec/fixtures/http/os_client_lookup_hashes.yml", "spec/fixtures/http/os_client_search_subtitles_hash.yml", "spec/fixtures/http/os_client_search_subtitles_imdb_id.yml", "spec/fixtures/http/os_client_search_subtitles_lang.yml", "spec/fixtures/http/os_id_by_hash_ep.yml", "spec/fixtures/http/os_id_by_hash_movie.yml", "spec/fixtures/http/os_id_by_hash_wrong.yml", "spec/fixtures/http/os_id_by_hashes.yml", "spec/fixtures/http/os_login.yml", "spec/fixtures/http/os_logout.yml", "spec/fixtures/http/os_lookup_hash_ep.yml", "spec/fixtures/http/os_lookup_hash_movie.yml", "spec/fixtures/http/os_lookup_hash_wrong.yml", "spec/fixtures/http/os_lookup_hashes.yml", "spec/fixtures/http/os_safe_call.yml", "spec/fixtures/http/os_search_subtitles.yml", "spec/fixtures/http/os_search_subtitles_for_hash.yml", "spec/fixtures/http/os_search_subtitles_for_imdb_id.yml", "spec/fixtures/http/os_search_subtitles_for_imdb_id_lang.yml", "spec/fixtures/http/os_search_subtitles_for_title.yml", "spec/spec_helper.rb"]

  if s.respond_to? :specification_version then
    s.specification_version = 4

    if Gem::Version.new(Gem::VERSION) >= Gem::Version.new('1.2.0') then
      s.add_runtime_dependency(%q<xdg>, ["~> 2"])
      s.add_runtime_dependency(%q<version>, ["~> 1"])
      s.add_runtime_dependency(%q<iso-639>, ["~> 0"])
      s.add_runtime_dependency(%q<streamio-ffmpeg>, ["~> 1"])
      s.add_runtime_dependency(%q<colorize>, ["~> 0"])
      s.add_runtime_dependency(%q<parallel>, ["~> 1"])
      s.add_runtime_dependency(%q<ruby-terminfo>, ["~> 0"])
      s.add_development_dependency(%q<rspec>, ["~> 3"])
      s.add_development_dependency(%q<rubocop>, ["~> 0"])
      s.add_development_dependency(%q<vcr>, ["~> 2"])
      s.add_development_dependency(%q<webmock>, ["~> 1"])
      s.add_development_dependency(%q<fakefs>, ["~> 0"])
    else
      s.add_dependency(%q<xdg>, ["~> 2"])
      s.add_dependency(%q<version>, ["~> 1"])
      s.add_dependency(%q<iso-639>, ["~> 0"])
      s.add_dependency(%q<streamio-ffmpeg>, ["~> 1"])
      s.add_dependency(%q<colorize>, ["~> 0"])
      s.add_dependency(%q<parallel>, ["~> 1"])
      s.add_dependency(%q<ruby-terminfo>, ["~> 0"])
      s.add_dependency(%q<rspec>, ["~> 3"])
      s.add_dependency(%q<rubocop>, ["~> 0"])
      s.add_dependency(%q<vcr>, ["~> 2"])
      s.add_dependency(%q<webmock>, ["~> 1"])
      s.add_dependency(%q<fakefs>, ["~> 0"])
    end
  else
    s.add_dependency(%q<xdg>, ["~> 2"])
    s.add_dependency(%q<version>, ["~> 1"])
    s.add_dependency(%q<iso-639>, ["~> 0"])
    s.add_dependency(%q<streamio-ffmpeg>, ["~> 1"])
    s.add_dependency(%q<colorize>, ["~> 0"])
    s.add_dependency(%q<parallel>, ["~> 1"])
    s.add_dependency(%q<ruby-terminfo>, ["~> 0"])
    s.add_dependency(%q<rspec>, ["~> 3"])
    s.add_dependency(%q<rubocop>, ["~> 0"])
    s.add_dependency(%q<vcr>, ["~> 2"])
    s.add_dependency(%q<webmock>, ["~> 1"])
    s.add_dependency(%q<fakefs>, ["~> 0"])
  end
end
