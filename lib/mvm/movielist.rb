require 'ostruct'

module Mvm
  class MovieList   # really, really need a better name for this thing
    attr_reader :movies

    class Movie < OpenStruct
    end

    def movie_class
      Movie
    end

    def initialize
      @movies = []
    end
  end
end

# vim: set shiftwidth=2:
