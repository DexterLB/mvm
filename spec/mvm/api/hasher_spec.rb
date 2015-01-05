require 'spec_helper'
require 'mvm/hasher'

module Mvm
  describe Hasher do
    subject { Hasher }

    describe '.hash' do
      let(:sample_video) do
        File.dirname(__FILE__) + '/samples/drop.avi'
      end

      it 'hashes sample movie correctly' do
        expect(subject.hash(sample_video)).to eq('450f3f0c98a1f11d')
      end
    end
  end
end

# vim: set shiftwidth=2:
