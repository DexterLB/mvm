require 'mvm/settings'

module Mvm
  DEFAULT_SETTINGS.merge!(
    movie_pattern:
      '%{library_folder}/movies/%{title} (%{year})/%{title}.?{extension}',
    episode_pattern:
      '%{library_folder}/series/%{series_title}/' \
      'S%<season_number>02dE%<episode_number>02d' \
      '- %{episode_title}%{extension}'
  )

  class Library
    class Renamer
      def initialize(settings = Settings.new)
        @settings = settings
      end

      def rename_movies(movies)
        movies.each { |movie| rename_movie(movie) }
      end

      def rename_movie(movie)
        new_filename = format(
          { episode: @settings.episode_pattern,
            movie:   @settings.movie_pattern
          }[movie.type],
          movie.to_h.merge(@settings.to_h)
        )

        puts movie.filename + ' -> ' + new_filename
      end
    end
  end
end

# vim: set shiftwidth=2:
