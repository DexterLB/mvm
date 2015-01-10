require 'ostruct'
require 'mvm/files'
require 'mvm/metadata'

module Mvm
  class MovieList   # really, really need a better name for this thing
    attr_accessor :movies

    def initialize
      @movies = []
    end

    include Files
    include Metadata
  end
end

# vim: set shiftwidth=2:
