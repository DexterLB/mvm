require 'spec_helper'
require 'mvm/opensubtitles'

module Mvm
  describe Opensubtitles do
    subject { @opensubtitles }

    before :all do
      @opensubtitles = Opensubtitles.new useragent: 'OSTestUserAgent'
    end

    describe '#info' do
      let(:info) { @opensubtitles.info }

      it 'returns non-nil info' do
        expect(info).not_to be_nil
      end

      it 'gives correct website url' do
        expect(info['website_url']).to eq('http://www.opensubtitles.org')
      end
    end
  end
end

# vim: set shiftwidth=2:
