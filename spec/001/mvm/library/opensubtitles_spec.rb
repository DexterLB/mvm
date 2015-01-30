require 'spec_helper'
require 'ostruct'
require 'mvm/library/opensubtitles'
require 'mvm/settings'

module Mvm
  class Library
    describe Opensubtitles do
      subject { @opensubtitles }

      before :all do
        settings = Settings.new
        @opensubtitles = Opensubtitles.new settings
      end

      describe '#id_by_hash_for' do
        it 'identifies movie correctly' do
          VCR.use_cassette('os_id_by_hash_movie') do
            movie = OpenStruct.new(file_hash: '09a2c497663259cb')
            expect(subject.id_by_hash_for(movie).to_h
                  ).to include(
                    type: :movie,
                    title: 'Nochnoy dozor',
                    year: 2004,
                    imdb_id: '0403358'
                  )
          end
        end

        it 'identifies episode correctly' do
          VCR.use_cassette('os_id_by_hash_ep') do
            movie = OpenStruct.new(file_hash: '46e33be00464c12e')
            expect(subject.id_by_hash_for(movie).to_h
                  ).to include(
                    type: :episode,
                    episode_title: 'Two Swords',
                    year: 2014,
                    imdb_id: '2816136',
                    series_title: 'Game of Thrones',
                    season_number: 4,
                    episode_number: 1
                  )
          end
        end

        it 'does nothing unidentified movie' do
          VCR.use_cassette('os_id_by_hash_wrong') do
            movie = OpenStruct.new(file_hash: '450f3f0c98a1f11d')
            unchanged = movie.dup
            expect(subject.id_by_hash_for(movie)).to eq(unchanged)
          end
        end

        it 'doesn\'t mutate the movie object' do
          VCR.use_cassette('os_id_by_hash_wrong') do
            movie = OpenStruct.new(file_hash: '450f3f0c98a1f11d')
            old_movie = movie.dup
            subject.id_by_hash_for(movie)
            expect(movie).to eq(old_movie)
          end
        end
      end

      describe '#id_by_hashes' do
        let(:hashes) { %w(09a2c497663259cb 46e33be00464c12e) }

        it 'identifies multiple items correctly' do
          VCR.use_cassette('os_id_by_hashes') do
            movies = hashes.map { |hash| OpenStruct.new(file_hash: hash) }
            expect(subject.id_by_hashes(movies).map do |movie|
              { movie.file_hash => movie.imdb_id }
            end.reduce(&:merge)).to eq(
              '09a2c497663259cb' => '0403358',
              '46e33be00464c12e' => '2816136'
            )
          end
        end

        it 'doesn\'t mutate the movies object' do
          VCR.use_cassette('os_id_by_hashes') do
            movies = hashes.map { |hash| OpenStruct.new(file_hash: hash) }
            old_movies = movies.dup
            subject.id_by_hashes(movies)
            expect(movies).to eq(old_movies)
          end
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
