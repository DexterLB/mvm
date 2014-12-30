require 'spec_helper'
require 'mvm/opensubtitles'

module Mvm
  describe Opensubtitles do
    subject { @opensubtitles }

    before :all do
      @opensubtitles = Opensubtitles.new useragent: 'OSTestUserAgent'
    end

    describe '#info' do
      let(:info) { subject.info }

      it 'returns non-nil info' do
        expect(info).not_to be_nil
      end

      it 'gives correct website url' do
        expect(info['website_url']).to eq('http://www.opensubtitles.org')
      end
    end

    describe '#login' do
      it 'changes token to non-nil' do
        expect { subject.login }.to change(subject, :token).from(nil)
      end
    end

    describe '#logout' do
      it 'changes token to nil' do
        expect { subject.logout }.to change(subject, :token).to(nil)
      end
    end
  end
end

# vim: set shiftwidth=2:
