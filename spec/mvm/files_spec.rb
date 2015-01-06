require 'spec_helper'
require 'mvm/modules/files'

module Mvm
  module Modules
    describe Files do
      class DummyMovieList
        class Movie < OpenStruct
        end

        def movie_class
          Movie
        end

        def movies
          @movies ||= []
        end

        include Files
      end

      subject { @dummy }

      before :each do
        @dummy = DummyMovieList.new
      end

      describe '#add_movies' do
        it 'adds movies with correct filenames' do
          subject.add_movies(['foo', 'bar', 'baz'])
          expect(subject.movies.map { |movie| movie.filename }).to eq(
            ['foo', 'bar', 'baz']
          )
        end
      end

    end
  end
end

# vim: set shiftwidth=2:
