require 'ostruct'
require 'mvm/library/files'
require 'mvm/library/metadata'
require 'mvm/library/opensubtitles'
require 'mvm/library/settings'

module Mvm
  class Library   # really, really need a better name for this thing
    attr_accessor :movies

    def initialize
      @movies = []
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

    def identify
      calculate_hashes.id_by_hashes
    end

    def identify!
      calculate_hashes!
      id_by_hashes!
    end

    def scan_folder!(folder)
      @movies += files.scan_folder folder
      self
    end

    def scan_folder(folder)
      dup.scan_folder!(folder)
    end

    def settings
      @settings ||= Settings.new
    end

    private

    def opensubtitles
      @opensubtitles ||= Opensubtitles.new settings
    end

    def files
      @files ||= Files.new settings
    end
  end
end

# vim: set shiftwidth=2:
