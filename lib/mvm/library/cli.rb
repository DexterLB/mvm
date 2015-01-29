require 'iso-639'
require 'colorize'

module Mvm
  class Library
    # Console interface and user interaction
    class Cli
      def self.print_movie(movie, spoiler_level: nil)
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
        puts movie.plot[spoiler_level] if movie.plot and spoiler_level
      end

      def self.print_movies(movies, spoiler_level: nil)
        movies.each { |movie| print_movie(movie, spoiler_level) }
        nil
      end
    end
  end
end

# vim: set shiftwidth=2:
