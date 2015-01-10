require 'spec_helper'
require 'mvm/metadata'

module Mvm
  describe Files do
    class DummyMovieList
      attr_accessor :movies

      def initialize
        @movies = []
      end

      include Metadata
    end

    subject { @dummy }

    before :each do
      @dummy = DummyMovieList.new
    end

    describe '#read_metadata' do
      let(:sample_video) do
        File.dirname(__FILE__) + '/api/samples/drop.avi'
      end

      it 'reads correct metadata' do
        subject.movies = [OpenStruct.new(filename: sample_video)]
        subject.read_metadata
        expect(subject.movies.first.to_h).to include(
          width: 256,
          height: 240,
          frame_rate: 30.0,
          duration: 6.07
        )
      end

      it 'returns the movies object' do
        subject.movies = [OpenStruct.new(filename: sample_video)]
        expect(subject.read_metadata).to equal(subject.movies)
      end
    end
  end
end

# vim: set shiftwidth=2:
