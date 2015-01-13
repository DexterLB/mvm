require 'ostruct'
require 'find'
require 'streamio-ffmpeg'
require 'mvm/api/hasher'
require 'mvm/settings'

module Mvm
  DEFAULT_SETTINGS.merge!(
    valid_movie_extensions: '.mkv .avi .mp4'
  )

  class Files
    def initialize(settings)
      @settings = settings
    end

    def self.calculate_hashes(movies)
      movies.map { |movie| calculate_hash_for(movie) }
    end

    def self.calculate_hash_for(movie)
      movie = movie.dup

      movie.file_hash = Api::Hasher.hash(movie.filename)

      movie
    end

    def self.movies_from_filenames(filenames)
      filenames.map { |filename| OpenStruct.new filename: filename }
    end

    def valid_movie?(filename)
      return false unless File.file? filename
      @settings.valid_movie_extensions.split.include? File.extname(filename)
    end

    def scan_folder(folder)
      files = Find.find(folder).to_a.select { |file| valid_movie? file }
      self.class.movies_from_filenames files
    end
  end
end

# vim: set shiftwidth=2:
