require 'spec_helper'
require 'mvm/movielist'

module Mvm
  describe MovieList do
    subject { @movielist }

    before :each do
      @movielist = MovieList.new
    end

    describe '#initialize' do
      it 'has an empty movie list when created' do
        expect(subject.movies).to eq([])
      end
    end
  end
end

# vim: set shiftwidth=2:
