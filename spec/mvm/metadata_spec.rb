require 'spec_helper'
require 'mvm/metadata'

module Mvm
  describe Metadata do
    subject { Metadata }

    let(:sample_video) do
      File.dirname(__FILE__) + '/api/samples/drop.avi'
    end

    describe '.read_metadata_for' do
      it 'reads correct metadata' do
        movie = OpenStruct.new(filename: sample_video)
        subject.read_metadata_for(movie)
        expect(movie.to_h).to include(
          width: 256,
          height: 240,
          frame_rate: 30.0,
          duration: 6.07
        )
      end

      it 'returns the movie' do
        movie = OpenStruct.new(filename: sample_video)
        expect(subject.read_metadata_for(movie)).to equal(movie)
      end
    end

    describe '.read_metadata' do
      it 'reads correct metadata for a list of one movie' do
        movies = [OpenStruct.new(filename: sample_video)]
        subject.read_metadata(movies)
        expect(movies.first.to_h).to include(
          width: 256,
          height: 240,
          frame_rate: 30.0,
          duration: 6.07
        )
        expect(movies.size).to eq(1)
      end

      it 'returns the movies object' do
        movies = [OpenStruct.new(filename: sample_video)]
        expect(subject.read_metadata(movies)).to equal(movies)
      end
    end
  end
end

# vim: set shiftwidth=2:
