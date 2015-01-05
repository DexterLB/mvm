require 'spec_helper'
require 'mvm/opensubtitles_client'
require 'mvm/opensubtitles_error'

module Mvm
  describe OpensubtitlesClient do
    subject { @client }

    before :all do
      @client = OpensubtitlesClient.new useragent: 'OSTestUserAgent'
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

    describe '#safe_client_call' do
      it 'throws UnauthorizedError when used with wrong token' do
        expect {
          subject.safe_client_call('CheckMovieHash',
                                   '',
                                   ['46e33be00464c12e'])
        }.to raise_error(OpensubtitlesClient::UnauthorizedError)
      end
    end

    describe '#call' do
      it 'logs in automatically and successfully calls function' do
        subject.logout
        expect(subject.call('CheckMovieHash',
                            ['46e33be00464c12e'])['data'].keys.first
                           ).to eq('46e33be00464c12e')
      end
    end
  end
end

# vim: set shiftwidth=2:
