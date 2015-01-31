require 'mvm/settings'
require 'imdb'
require 'iso-639'

module Mvm
  class Library
    DEFAULT_SETTINGS.merge!(
        {}
    )

    class Imdb
      def initialize(settings)
        @settings = settings
      end

      def get_data_for(movie)
        set_data_for(movie, ::Imdb::Movie.new(movie.imdb_id)) rescue movie
      end

      def get_data(movies)
        movies.each_with_index.map do |movie, current|
          yield [current, movies.size] if block_given?
          get_data_for(movie)
        end
      end

      def set_data_for(movie, data)
        movie = movie.dup

        return movie unless data

        [:director, :trailer_url, :genres, :languages, :company, :rating,
         :tagline, :year
        ].each { |attribute| movie[attribute] = data.send(attribute) }

        movie.poster_url = data.poster
        movie.plot = [data.plot, data.plot_summary, data.plot_synopsis]

        movie.type = data.episode? ? :episode : :movie

        if movie.type == :episode
          movie.series_title = data.title
          movie.episode_title = data.episode_title
          movie.season_number = data.episode_season
          movie.episode_number = data.episode_number
          movie.series_imdb_id = data.episode_serie_id
        else
          movie.title = data.title
        end

        movie
      end

      def self.url(imdb_id)
        "http://www.imdb.com/title/tt#{imdb_id}"
      end

      def self.id(soup)
        match = /(.*\.imdb\.com\/title\/tt)?(?<id>[0-9]+)(\/.*)?/.match(soup)
        if match
          match['id']
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
