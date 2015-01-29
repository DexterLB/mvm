require 'iso-639'
require 'colorize'

module Mvm
  class Library
    # Console interface and user interaction
    class Cli
      def self.print_movie(movie)
        puts(
          if movie.title
            if movie.type == :episode
              movie.series_title.cyan +
                ' S%02dE%02d '.light_blue.format([movie.season_number,
                                                  movie.episode_number]) +
                movie.episode_title.cyan
            else
              movie.title.cyan
            end + " (#{movie.year})".green
          else
            movie.filename.yellow
          end
        )
      end

      def self.print_movies(movies)
        movies.each { |movie| print_movie(movie) }
        nil
      end
    end
  end
end

# vim: set shiftwidth=2:
