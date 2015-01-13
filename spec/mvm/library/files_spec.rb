require 'spec_helper'
require 'mvm/library/files'

module Mvm
  class Library
    describe Files do
      subject { Files }

      let(:sample_video) do
        File.dirname(__FILE__) + '/../../fixtures/drop.avi'
      end

      let(:movie) do
        OpenStruct.new(filename: sample_video)
      end

      describe '.calculate_hash_for' do
        it 'calculates correct hash' do
          expect(subject.calculate_hash_for(movie).file_hash)
            .to eq('450f3f0c98a1f11d')
        end

        it 'doesn\'t mutate the movie object' do
          old_movie = movie.dup
          subject.calculate_hash_for(movie)
          expect(movie).to eq(old_movie)
        end
      end

      describe '.calculate_hashes' do
        let(:movies) { [movie] }
        it 'calculates correct hash for list of one movie' do
          result = subject.calculate_hashes(movies)
          expect(result.first.file_hash).to eq('450f3f0c98a1f11d')
          expect(result.size).to eq(1)
        end

        it 'doesn\'t mutate the movies object' do
          old_movies = movies.dup
          subject.calculate_hashes(movies)
          expect(movies).to eq(old_movies)
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
