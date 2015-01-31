require 'iso-639'
require 'colorize'

module Mvm
  class Library
    # Console interface and user interaction
    class Cli
      def self.print_movie(movie, spoiler_level: nil)
        print_title(movie) || print_filename(movie)
        puts
        print_plot(movie, spoiler_level: spoiler_level)
      end

      def self.print_title(movie)
        if movie.title
          print(
            if movie.type == :episode
              movie.series_title.cyan +
                ' S%02dE%02d '.light_blue.format([movie.season_number,
                                                  movie.episode_number]) +
                movie.episode_title.cyan
            else
              movie.title.cyan
            end + (" (#{movie.year})".green if movie.year)
          )
          true
        else
          false
        end
      end

      def self.print_filename(movie)
        print movie.filename.yellow
        true
      end

      def self.print_plot(movie, spoiler_level: 0)
        if movie.plot && spoiler_level
          puts movie.plot[spoiler_level]
          true
        else
          false
        end
      end

      def self.print_imdb(movie)
        if movie.imdb_id
          print "[ ".light_blue + Imdb.url(movie.imdb_id) + " ]".blue
          true
        else
          false
        end
      end

      def self.print_movies(movies, spoiler_level: nil)
        movies.each { |movie| print_movie(movie, spoiler_level: spoiler_level) }
        nil
      end

      def self.swear
        "fuck you."
      end

      def self.get_imdb_id
        begin
          input = gets.chomp
          return nil if input.empty?
          id = Imdb.id(input)
          if not id
            print swear + " try again: "
          end
        end while not id
        id
      end

      def self.ask_imdb_id_for(movie)
        movie = movie.dup

        print_filename(movie)
        puts
        print_title(movie) && puts
        print 'imdb id: '
        print_imdb(movie) && (print ' ')

        movie.imdb_id = get_imdb_id || movie.imdb_id

        movie
      end
    end
  end
end

# vim: set shiftwidth=2:
