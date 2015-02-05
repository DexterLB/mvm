require 'ostruct'
require 'mvm/settings'
require 'mvm/library/files'
require 'mvm/library/metadata'
require 'mvm/library/opensubtitles'
require 'mvm/library/imdb'
require 'mvm/library/cli'
require 'mvm/library/store'
require 'mvm/library/rename'

module Mvm
  class Library   # really, really need a better name for this thing
    attr_accessor :movies

    def initialize(settings: nil)
      @movies = []
      @settings = settings
    end

    def calculate_hashes!
      @movies = Files.calculate_hashes(@movies)
      self
    end

    def calculate_hashes
      dup.calculate_hashes!
    end

    def id_by_hashes!
      @movies = opensubtitles.id_by_hashes(@movies)
      self
    end

    def id_by_hashes
      dup.id_by_hashes!
    end

    def read_metadata!
      @movies = Metadata.read_metadata(movies)
      self
    end

    def read_metadata
      dup.read_metadata!
    end

    def ask_imdb_ids!
      @movies = Cli.ask_imdb_ids(movies)
      self
    end

    def ask_imdb_ids
      dup.ask_imdb_ids!
    end

    def get_data!(&block)
      @movies = imdb.get_data(movies, &block)
      self
    end

    def get_data(&block)
      dup.get_data!(&block)
    end

    def print(spoiler_level: nil)
      Cli.print_movies(@movies, spoiler_level: spoiler_level)
      self
    end

    def scan_folder!(folder)
      @movies += files.scan_folder folder
      self
    end

    def scan_folder(folder)
      dup.scan_folder!(folder)
    end

    def rename_movies!
      @movies = renamer.rename_movies(@movies)
      self
    end

    def rename_movies
      dup.rename_movies!
    end

    def store_movies
      store.store_movies(@movies)
    end

    def load_movies!(folder)
      @movies += store.load_movies(folder)
      self
    end

    def load_movies(folder)
      dup.load_movies!(folder)
    end

    def settings
      @settings ||= Settings.new
    end

    private

    def opensubtitles
      @opensubtitles ||= Opensubtitles.new settings
    end

    def imdb
      @imdb ||= Imdb.new settings
    end

    def files
      @files ||= Files.new settings
    end

    def store
      @store ||= Store.new settings
    end

    def renamer
      @renamer ||= Renamer.new settings
    end
  end
end

# vim: set shiftwidth=2:
