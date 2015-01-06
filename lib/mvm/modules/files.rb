require 'ostruct'

module Mvm
  module Modules
    module Files
      def add_movies(filenames)
        filenames.each do |filename|
          movies << movie_class.new(filename: filename)
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
