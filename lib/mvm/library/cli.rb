require 'iso-639'
require 'colorize'
require 'terminfo'
require 'mvm/library/imdb'

class Array
  def stretch(stretch_length)
    (0...stretch_length).map do |index|
      self[((index.to_f / stretch_length.to_f) * size.to_f).to_i]
    end
  end
end

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
                format(' S%02dE%02d '.light_blue,
                       movie.season_number, movie.episode_number) +
                movie.episode_title.cyan
            else
              movie.title.cyan
            end + (" (#{movie.year})".green if movie.year)
          )
          true
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
        end
      end

      def self.print_imdb(movie)
        if movie.imdb_id
          print '[ '.light_blue + Imdb.url(movie.imdb_id) + ' ]'.blue
          true
        end
      end

      def self.print_movies(movies, spoiler_level: nil)
        movies.each { |movie| print_movie(movie, spoiler_level: spoiler_level) }
        nil
      end

      def self.swear
        [
          'fuck you.',
          'you suck.',
          'eat shit.'
        ].sample
      end

      def self.obtain_imdb_id
        loop do
          input = $stdin.gets.chomp
          return nil if input.empty?

          id = Imdb.id(input)
          return id if id

          print(swear + ' try again: ')
        end
      end

      def self.ask_imdb_id_for(movie)
        movie = movie.dup

        print_filename(movie)
        puts
        print_title(movie) && puts
        print 'imdb id: '
        print_imdb(movie) && (print ' ')

        movie.imdb_id = obtain_imdb_id || movie.imdb_id

        movie
      end

      def self.ask_imdb_ids(movies, only_missing: true)
        movies.map do |movie|
          if movie.imdb_id && only_missing
            movie
          else
            ask_imdb_id_for(movie)
          end
        end
      end

      def self.progressbar(finished, all)
        progress = [:finished] * finished
        progress << :processing if finished < all
        progress += [:pending] * (all - progress.size)
        multi_progressbar(progress)
      end

      def self.multi_progressbar(progress)
        length = TermInfo.screen_columns - 1

        left = ' ['
        right = format(
          '] %d/%d',
          progress.select { |item| item == :finished }.size,
          progress.size
        )

        bar = progress.stretch(length - (left + right).size).map do |item|
          { pending: ' ', processing: '-', finished: '#' }[item]
        end.to_a.join

        print left + bar + right + "\r"
        STDOUT.flush
      end
    end
  end
end

# vim: set shiftwidth=2:
