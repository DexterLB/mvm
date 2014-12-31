require 'mvm/opensubtitles_client'

module Mvm
  class Opensubtitles
    def initialize(**client_options)
      @client = OpensubtitlesClient.new(**client_options)
    end

    def id_by_hash(hash)
      id_by_hashes([hash])[hash]
    end

    def id_by_hashes(hashes)
      @client.call('CheckMovieHash', hashes)['data']
    end
  end
end

# vim: set shiftwidth=2:
