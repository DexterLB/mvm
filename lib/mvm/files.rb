require 'ostruct'
require 'streamio-ffmpeg'
require 'mvm/api/hasher'

module Mvm
  class Files
    def self.calculate_hashes(movies)
      movies.each { |movie| calculate_hash_for(movie) }
    end

    def self.calculate_hash_for(movie)
      movie.file_hash = Api::Hasher.hash(movie.filename)

      movie
    end
  end
end

# vim: set shiftwidth=2:
