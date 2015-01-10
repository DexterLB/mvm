require 'mvm/api/opensubtitles_client'
require 'mvm/configuration'

require 'iso-639'

module Mvm
  Configuration::DEFAULTS.merge!(
    opensubtitles_username: '',
    opensubtitles_password: '',
    opensubtitles_useragent: 'OSTestUserAgent',
    opensubtitles_language: 'en',
    opensubtitles_timeout: 20
  )

  module Opensubtitles
    class OpensubtitlesApi
      def initialize(**client_options)
        @client = Api::OpensubtitlesClient.new(**client_options)
      end

      def id_by_hash(hash)
        id_by_hashes([hash])[hash]
      end

      def id_by_hashes(hashes)
        @client.call('CheckMovieHash', hashes)['data'].map do |hash, data|
          if data.empty? # XMLRPC returns [] instead of {} when it's empty
            [hash, {}]
          else
            [hash, data]
          end
        end.to_h
      end

      def self.set_attributes_for(movie, attributes)
        return movie unless %w(episode movie).include? attributes['MovieKind']

        movie.title = attributes['MovieName']
        movie.type = attributes['MovieKind'].to_sym
        movie.year = attributes['MovieYear'].to_i
        movie.imdb_id = attributes['MovieImdbID']

        if movie.type == :episode
          movie.series_title, movie.episode_title = movie.title
                                                    .match(/\"(.+)\"\s(.+)/)
                                                    .captures

          movie.season_number = attributes['SeriesSeason'].to_i
          movie.episode_number = attributes['SeriesEpisode'].to_i
        end

        movie
      end
    end

    def opensubtitles
      @client ||= OpensubtitlesApi.new(
        username: settings.opensubtitles_username,
        password: settings.opensubtitles_password,
        useragent: settings.opensubtitles_useragent,
        language: ISO_639.find(settings.opensubtitles_language),
        timeout: settings.opensubtitles_timeout
      )
    end

    def id_by_hash_for(movie)
      OpensubtitlesApi.set_attributes_for(movie, opensubtitles
                                                 .id_by_hash(movie.file_hash))
      movie
    end

    def id_by_hashes
      data = opensubtitles.id_by_hashes(movies.map(&:file_hash))
      movies.each do |movie|
        OpensubtitlesApi.set_attributes_for(movie, data[movie.file_hash])
      end
    end
  end
end

# vim: set shiftwidth=2:
