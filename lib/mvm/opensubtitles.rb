require 'mvm/opensubtitles_client'

module Mvm
  class Opensubtitles
    def initialize(**client_options)
      @client = OpensubtitlesClient.new(**client_options)
    end

    def id_by_hash(hash)
      @client.call('CheckMovieHash', [hash])['data'][hash]
    end
  end
end

# vim: set shiftwidth=2:
