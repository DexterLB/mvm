require 'spec_helper'
require 'mvm/library/metadata'

module Mvm
  class Library
    describe Metadata do
      subject { Metadata }

      let(:sample_video) do
        File.dirname(__FILE__) + '/../../fixtures/drop.avi'
      end

      let(:movie) do
        OpenStruct.new(filename: sample_video)
      end

      describe '.read_metadata_for' do
        it 'reads correct metadata' do
          expect(subject.read_metadata_for(movie).to_h).to include(
            width: 256,
            height: 240,
            frame_rate: 30.0,
            duration: 6.07
          )
        end

        it 'doesn\'t mutate the movie object' do
          old_movie = movie.dup
          subject.read_metadata_for(movie)
          expect(movie).to eq(old_movie)
        end
      end

      describe '.read_metadata' do
        let(:movies) { [movie] }

        it 'reads correct metadata for a list of one movie' do
          result = subject.read_metadata(movies)
          expect(result.first.to_h).to include(
            width: 256,
            height: 240,
            frame_rate: 30.0,
            duration: 6.07
          )
          expect(result.size).to eq(1)
        end

        it 'doesn\'t mutate the movies object' do
          old_movies = movies.dup
          subject.read_metadata(movies)
          expect(movies).to eq(old_movies)
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
