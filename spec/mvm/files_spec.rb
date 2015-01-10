require 'spec_helper'
require 'mvm/files'

module Mvm
  describe Files do
    class DummyMovieList
      attr_accessor :movies

      def initialize
        @movies = []
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

      it 'returns the movies object' do
        expect(subject.add_movies(%w(foo bar baz))).to equal(subject.movies)
      end
    end

    describe '#calculate_hashes' do
      let(:sample_video) do
        File.dirname(__FILE__) + '/api/samples/drop.avi'
      end

      it 'calculates correct hash' do
        subject.movies = [OpenStruct.new(filename: sample_video)]
        subject.calculate_hashes
        expect(subject.movies.first.file_hash).to eq('450f3f0c98a1f11d')
      end

      it 'returns the movies object' do
        subject.movies = [OpenStruct.new(filename: sample_video)]
        expect(subject.calculate_hashes).to equal(subject.movies)
      end
    end
  end
end

# vim: set shiftwidth=2:
