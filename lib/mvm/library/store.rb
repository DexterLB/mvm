require 'iso-639'
require 'yaml'

module Mvm
  DEFAULT_SETTINGS.merge!(
    store_file: '%{filename}.mvm',
    store_match: '\.mvm$'
  )

  class Library
    # Stores movie data into yaml files
    class Store
      def initialize(settings)
        @settings = settings
      end

      def store_movie(movie)
        File.write(@settings.store_file % movie.to_h, movie.to_h.to_yaml)
      end

      def store_movies(movies)
        movies.each { |movie| store_movie(movie) }
      end

      def load_movie(yaml_file)
        OpenStruct.new(YAML.load(File.read(yaml_file)))
      end

      def load_movies(folder)
        valid = Regexp.new(@settings.store_match)
        Find.find(folder).to_a
          .select { |file| valid.match file }
          .map    { |file| load_movie(file) }
      end
    end
  end
end

# vim: set shiftwidth=2:
