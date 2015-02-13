require 'ostruct'
require 'find'
require 'streamio-ffmpeg'
require 'mvm/opensubtitles_client/hasher'
require 'mvm/settings'

module Mvm
  class Library
    DEFAULT_SETTINGS.merge!(
      valid_movie_extensions: '.mkv .avi .mp4'
    )

    class Files
      def initialize(settings = Settings.new)
        @settings = settings
      end

      def calculate_hashes(movies)
        movies.map { |movie| calculate_hash_for(movie) }
      end

      def calculate_hash_for(movie)
        movie = movie.dup

        movie.file_hash = OpensubtitlesClient::Hasher.hash(movie.filename)

        movie
      end

      def load(path)
        if File.directory? path
          scan_folder path
        elsif File.file? path
          movies_from_filenames [path]
        else
          fail 'Not a file/folder: ' + path
        end
      end

      private

      def movies_from_filenames(filenames)
        filenames.map do |filename|
          OpenStruct.new(
            filename: filename,
            filesize: File.size(filename),
            basename: File.basename(filename),
            added: File.mtime(filename),
            extension: File.extname(filename)
          )
        end
      end

      def valid_movie?(filename)
        return false unless File.file? filename
        @settings.valid_movie_extensions.split.include? File.extname(filename)
      end

      def scan_folder(folder)
        files = Find.find(folder).to_a.select { |file| valid_movie? file }
        movies_from_filenames files
      end
    end
  end
end

# vim: set shiftwidth=2:
