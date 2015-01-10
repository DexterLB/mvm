require 'ostruct'
require 'streamio-ffmpeg'
require 'mvm/api/hasher'

module Mvm
  module Files
    def add_movies(filenames)
      filenames.each do |filename|
        movies << OpenStruct.new(filename: filename)
      end
    end

    def calculate_hashes
      movies.each do |movie|
        movie.file_hash = Api::Hasher.hash(movie.filename)
      end
    end
  end
end

# vim: set shiftwidth=2:
