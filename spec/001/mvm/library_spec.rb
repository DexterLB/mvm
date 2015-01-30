require 'spec_helper'
require 'mvm/library'

module Mvm
  describe Library do
    subject { @library }

    before :each do
      @library = Library.new
    end

    describe '#initialize' do
      it 'has an empty movie list when created' do
        expect(subject.movies).to eq([])
      end
    end
  end
end

# vim: set shiftwidth=2:
