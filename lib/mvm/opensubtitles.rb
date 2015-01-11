require 'mvm/api/opensubtitles_client'
require 'mvm/settings'

require 'iso-639'

module Mvm
  DEFAULT_SETTINGS.merge!(
    opensubtitles_username: '',
    opensubtitles_password: '',
    opensubtitles_useragent: 'OSTestUserAgent',
    opensubtitles_language: 'en',
    opensubtitles_timeout: 20
  )

  class Opensubtitles
    def initialize(settings)
      @settings = settings
    end

    def id_by_hash_for(movie)
      set_attributes_for(movie, lookup_hash(movie.file_hash))

      movie
    end

    def id_by_hashes(movies)
      data = lookup_hashes(movies.map(&:file_hash))
      movies.each do |movie|
        set_attributes_for(movie, data[movie.file_hash])
      end
    end

    def set_attributes_for(movie, attributes)
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

    def lookup_hash(hash)
      lookup_hashes([hash])[hash]
    end

    def lookup_hashes(hashes)
      client.call('CheckMovieHash', hashes)['data'].map do |hash, data|
        if data.empty? # XMLRPC returns [] instead of {} when it's empty
          [hash, {}]
        else
          [hash, data]
        end
      end.to_h
    end

    def client
      @client ||= Api::OpensubtitlesClient.new(
        username: @settings.opensubtitles_username,
        password: @settings.opensubtitles_password,
        useragent: @settings.opensubtitles_useragent,
        language: ISO_639.find(@settings.opensubtitles_language),
        timeout: @settings.opensubtitles_timeout
      )
    end
  end
end

# vim: set shiftwidth=2:
