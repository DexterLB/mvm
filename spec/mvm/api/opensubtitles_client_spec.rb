require 'spec_helper'
require 'vcr'
require 'mvm/api/opensubtitles_client'
require 'mvm/api/opensubtitles_error'

module Mvm
  module Api
    describe OpensubtitlesClient do
      subject { @client }

      before :all do
        @client = OpensubtitlesClient.new useragent: 'OSTestUserAgent'
      end

      describe '#info' do
        let(:info) { subject.info }

        it 'returns non-nil info' do
          VCR.use_cassette('os_client') do
            expect(info).not_to be_nil
          end
        end

        it 'gives correct website url' do
          VCR.use_cassette('os_client') do
            expect(info['website_url']).to eq('http://www.opensubtitles.org')
          end
        end
      end

      describe '#login' do
        it 'changes token to non-nil' do
          VCR.use_cassette('os_login') do
            expect { subject.login }.to change(subject, :token).from(nil)
          end
        end
      end

      describe '#logout' do
        it 'changes token to nil' do
          VCR.use_cassette('os_logout') do
            expect { subject.logout }.to change(subject, :token).to(nil)
          end
        end
      end

      describe '#safe_client_call' do
        it 'throws UnauthorizedError when used with wrong token' do
          VCR.use_cassette('os_safe_call') do
            expect do
              subject.safe_client_call('CheckMovieHash',
                                       '',
                                       ['46e33be00464c12e'])
            end.to raise_error(OpensubtitlesClient::UnauthorizedError)
          end
        end
      end

      describe '#call' do
        it 'logs in automatically and successfully calls function' do
          VCR.use_cassette('os_call') do
            subject.logout
            expect(subject.call('CheckMovieHash',
                                ['46e33be00464c12e'])['data'].keys.first
                              ).to eq('46e33be00464c12e')
          end
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
