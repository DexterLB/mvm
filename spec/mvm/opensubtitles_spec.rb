require 'spec_helper'
require 'mvm/opensubtitles'

module Mvm
  describe Opensubtitles do
    subject { @opensubtitles }

    before :all do
      @opensubtitles = Opensubtitles.new useragent: 'OSTestUserAgent'
    end

    describe '#info' do
      it 'returns non-nil info' do
        @opensubtitles.info.should_not be_nil
      end
    end
  end
end

# vim: set shiftwidth=2:
