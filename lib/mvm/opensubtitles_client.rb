require 'mvm/opensubtitles_client/api'
require 'zlib'
require 'open-uri'

module Mvm
  class OpensubtitlesClient
    def initialize(**client_options)
      @api = Api.new(**client_options)
    end

    def lookup_hash(hash)
      lookup_hashes([hash])[hash]
    end

    def lookup_hashes(hashes)
      hashes.each_slice(199).map do |hash_list|
        lookup_hashes_under_200 hash_list
      end.inject(&:merge)
    end

    def search_subtitles(queries)
      @api.call('SearchSubtitles', queries)['data'] || []
    end

    def download_gz(url, target_filename, encoding: nil)
      open(url, 'rb') do |reader|
        gz_reader = Zlib::GzipReader.new(reader, encoding: encoding)
        open(target_filename, 'w') do |writer|
          writer.write(gz_reader.read)
        end
        gz_reader.close
      end
    end

    private

    def lookup_hashes_under_200(hashes)
      @api.call('CheckMovieHash', hashes)['data'].map do |hash, data|
        if data.empty? # XMLRPC returns [] instead of {} when it's empty
          [hash, {}]
        else
          [hash, data]
        end
      end.to_h
    end
  end
end

# vim: set shiftwidth=2:
