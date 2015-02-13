require 'spec_helper'
require 'mvm/opensubtitles_client'

module Mvm
  describe OpensubtitlesClient do
    subject { @opensubtitles }

    before :all do
      @opensubtitles = OpensubtitlesClient.new useragent: 'OSTestUserAgent'
    end

    describe '#lookup_hash' do
      it 'identifies movie correctly' do
        VCR.use_cassette('os_client_lookup_hash_movie') do
          expect(subject.lookup_hash('09a2c497663259cb')
                ).to include(
                  'MovieKind'      => 'movie',
                  'MovieName'      => 'Nochnoy dozor',
                  'MovieYear'      => '2004',
                  'MovieImdbID'    => '0403358'
                )
        end
      end

      it 'identifies episode correctly' do
        VCR.use_cassette('os_client_lookup_hash_ep') do
          expect(subject.lookup_hash('46e33be00464c12e')
                ).to include(
                  'MovieKind'      => 'episode',
                  'MovieName'      => '"Game of Thrones" Two Swords',
                  'MovieYear'      => '2014',
                  'MovieImdbID'    => '2816136',
                  'SeriesSeason'   => '4',
                  'SeriesEpisode'  => '1'
                )
        end
      end

      it 'returns empty hash for non-existant movie' do
        VCR.use_cassette('os_client_lookup_hash_wrong') do
          expect(subject.lookup_hash('450f3f0c98a1f11d')).to eq({})
        end
      end
    end

    describe '#lookup_hashes' do
      let(:hashes) { %w(09a2c497663259cb 46e33be00464c12e) }

      it 'identifies multiple items correctly' do
        VCR.use_cassette('os_client_lookup_hashes') do
          expect(subject.lookup_hashes(hashes).map do |hash, data|
            { hash => data['MovieImdbID'] }
          end.reduce(&:merge)).to eq(Hash[hashes.map do |hash|
            [hash, subject.lookup_hash(hash)['MovieImdbID']]
          end])
        end
      end
    end

    describe '#search_subtitles' do
      it 'finds correct results for query with hash' do
        VCR.use_cassette('os_client_search_subtitles_hash') do
          results = subject.search_subtitles([
            { moviehash: '09a2c497663259cb' }
          ])
          expect(results).not_to be_empty
          valid_results = 0
          results.each do |result|
            valid_results += 1 if result['MovieName'] == 'Nochnoy dozor'
            expect(result['SubDownloadLink']).to match %r{^http://.*\.gz}
          end
          expect(valid_results > 1).to be_truthy
        end
      end

      it 'finds correct results for query with imdb_id' do
        VCR.use_cassette('os_client_search_subtitles_imdb_id') do
          results = subject.search_subtitles([
            { imdbid: '0403358' }
          ])
          expect(results).not_to be_empty
          valid_results = 0
          results.each do |result|
            valid_results += 1 if result['MovieName'] == 'Nochnoy dozor'
            expect(result['SubDownloadLink']).to match %r{^http://.*\.gz}
          end
          expect(valid_results > 1).to be_truthy
        end
      end

      it 'finds subtitles with correct languages' do
        VCR.use_cassette('os_client_search_subtitles_lang') do
          results = subject.search_subtitles([
            { moviehash: '09a2c497663259cb', sublanguageid: 'eng,rus' }
          ])
          expect(results).not_to be_empty
          results.each do |result|
            expect(%w(English Russian)).to include(result['LanguageName'])
          end
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
