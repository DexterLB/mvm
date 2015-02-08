require 'mvm/settings'

module Mvm
  DEFAULT_SETTINGS.merge!(
    movie_pattern:
      '%{library_folder}/movies/%{title} (%{year})/%{title}%{extension}',
    episode_pattern:
      '%{library_folder}/series/%{series_title}/' \
      'S%<season_number>02dE%<episode_number>02d' \
      '- %{episode_title}%{extension}',
    rename_strategy: 'symlink',
    fs_forbidden_char_exp: '[\x00\/\\:\*\?\"<>\|]', # windoze-friendly
    fs_forbidden_char_replace: '_'
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

        yield [movies.size, movies.size] if block_given?
        renamed_movies
      end

      def rename_movie(movie)
        movie = movie.dup

        new_filename = format(
          { episode: @settings.episode_pattern,
            movie:   @settings.movie_pattern
          }[movie.type],
          sanitize_hash(movie.to_h.merge(@settings.to_h))
        )

        mkdirs(new_filename)
        rename_file(movie.filename, new_filename)

        movie.filename = new_filename
        movie
      end

      private

      def mkdirs(filename)
        FileUtils.mkdir_p(File.dirname(filename))
      end

      def rename_file(old_name, new_name)
        strategy = format(@settings.rename_strategy,
                          old: old_name, new: new_name)
        case strategy
        when 'dummy' then nil
        when 'copy' then FileUtils.cp(old_name, new_name)
        when 'move' then FileUtils.mv(old_name, new_name)
        when 'symlink' then FileUtils.ln_sf(old_name, new_name)
        when 'keeplink'
          FileUtils.mv(old_name, new_name)
          FileUtils.ln_s(new_name, old_name)
        else exec_strategy(strategy)
        end
      end

      def exec_strategy(strategy)
        command_match = /exec: (.+)/.match(strategy)
        if command_match
          system(command_match[1])
        else
          fail 'Unknown rename strategy: ' + strategy
        end
      end

      private

      def sanitize_hash(hash)
        hash.map do |key, value|
          [key, (value.is_a? String) ? sanitize(value) : value]
        end.to_h
      end

      def sanitize(string)
        match = Regexp.new(@settings.fs_forbidden_char_exp)
        string.gsub(match, @settings.fs_forbidden_char_replace)
      end
    end
  end
end

# vim: set shiftwidth=2:
