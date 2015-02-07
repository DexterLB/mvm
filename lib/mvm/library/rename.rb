require 'mvm/settings'

module Mvm
  DEFAULT_SETTINGS.merge!(
    movie_pattern:
      '%{library_folder}/movies/%{title} (%{year})/%{title}.?{extension}',
    episode_pattern:
      '%{library_folder}/series/%{series_title}/' \
      'S%<season_number>02dE%<episode_number>02d' \
      '- %{episode_title}?{extension}'
  )

  class Library
    class Renamer
      def initialize(settings = Settings.new)
        @settings = settings
      end

      def rename_movies(movies)
        renamed_movies = movies.each.with_index.map do |movie, current|
          yield [current, movies.size] if block_given?

          rename_movie(movie)
        end

        yield [movies.size, movies.size]
        renamed_movies
      end

      def rename_movie(movie)
        movie = movie.dup

        new_filename = format(
          { episode: @settings.episode_pattern,
            movie:   @settings.movie_pattern
          }[movie.type],
          movie.to_h.merge(@settings.to_h)
        )

        sleep 0.1
        movie.filename = new_filename
        movie
      end
    end
  end
end

# vim: set shiftwidth=2:
