require 'spec_helper'
require 'mvm/opensubtitles'

module Mvm
  describe Opensubtitles do
    subject { @opensubtitles }

    before :all do
      @opensubtitles = Opensubtitles.new useragent: 'OSTestUserAgent'
    end

    describe '#id_by_hash' do
      it 'identifies movie correctly' do
        expect(subject.id_by_hash('09a2c497663259cb')
              ).to include('MovieKind'      => 'movie',
                           'MovieName'      => 'Nochnoy dozor',
                           'MovieYear'      => '2004',
                           'MovieImdbID'    => '0403358',
                          )
      end

      it 'identifies episode correctly' do
        expect(subject.id_by_hash('46e33be00464c12e')
              ).to include(
                           'MovieKind'      => 'episode',
                           'MovieName'      => '"Game of Thrones" Two Swords',
                           'MovieYear'      => '2014',
                           'MovieImdbID'    => '2816136',
                           'SeriesSeason'   => '4',
                           'SeriesEpisode'  => '1',
                          )
      end

      it 'returns empty hash for non-existant movie' do
        expect(subject.id_by_hash('450f3f0c98a1f11d')).to eq({})
      end
    end

    describe '#id_by_hashes' do
      let(:hashes) { ['09a2c497663259cb', '46e33be00464c12e'] }

      it 'identifies multiple items correctly' do
        expect(subject.id_by_hashes(hashes).map do |hash, data|
          { hash => data['MovieImdbID'] } 
        end.reduce(&:merge)).to eq(Hash[ hashes.map do |hash|
          [hash, subject.id_by_hash(hash)['MovieImdbID'] ]
        end ])
      end
    end
  end
end

# vim: set shiftwidth=2:
