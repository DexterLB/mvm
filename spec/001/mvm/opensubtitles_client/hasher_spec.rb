require 'spec_helper'
require 'mvm/opensubtitles_client/hasher'

module Mvm
  class OpensubtitlesClient
    describe Hasher do
      subject { Hasher }

      describe '.hash' do
        let(:sample_video) do
          File.dirname(__FILE__) + '/../../../fixtures/drop.avi'
        end

        let(:empty_file) do
          File.dirname(__FILE__) + '/../../../fixtures/empty.file'
        end

        it 'hashes sample movie correctly' do
          expect(subject.hash(sample_video)).to eq('450f3f0c98a1f11d')
        end

        it 'returns nil for empty file' do
          expect(subject.hash(empty_file)).to be_nil
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
