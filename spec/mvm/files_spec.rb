require 'spec_helper'
require 'mvm/files'

module Mvm
  describe Files do
    subject { Files }

    let(:sample_video) do
      File.dirname(__FILE__) + '/api/samples/drop.avi'
    end

    describe '.calculate_hash_for' do
      it 'calculates correct hash' do
        movie = OpenStruct.new(filename: sample_video)
        subject.calculate_hash_for(movie)
        expect(movie.file_hash).to eq('450f3f0c98a1f11d')
      end

      it 'returns the movies object' do
        movie = OpenStruct.new(filename: sample_video)
        expect(subject.calculate_hash_for(movie)).to equal(movie)
      end
    end

    describe '.calculate_hashes' do
      it 'calculates correct hash for list of one movie' do
        movies = [OpenStruct.new(filename: sample_video)]
        subject.calculate_hashes(movies)
        expect(movies.first.file_hash).to eq('450f3f0c98a1f11d')
        expect(movies.size).to eq(1)
      end

      it 'returns the movies object' do
        movies = [OpenStruct.new(filename: sample_video)]
        expect(subject.calculate_hashes(movies)).to equal(movies)
      end
    end
  end
end

# vim: set shiftwidth=2:
