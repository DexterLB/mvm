require 'spec_helper'
require 'mvm/files'

module Mvm
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
        subject.add_movies(%w(foo bar baz))
        expect(subject.movies.map(&:filename)).to eq(%w(foo bar baz))
      end
    end
  end
end

# vim: set shiftwidth=2:
