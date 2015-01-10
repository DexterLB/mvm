require 'ostruct'
require 'mvm/files'
require 'mvm/metadata'
require 'mvm/opensubtitles'
require 'mvm/configuration'

module Mvm
  class MovieList   # really, really need a better name for this thing
    attr_accessor :movies

    def initialize
      @movies = []
    end

    include Configuration

    include Files
    include Metadata
    include Opensubtitles
  end
end

# vim: set shiftwidth=2:
