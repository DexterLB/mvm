require 'spec_helper'
require 'ostruct'
require 'mvm/library/imdb'
require 'mvm/settings'

module Mvm
  class Library
    describe Imdb do
      subject { @imdb }

      before :all do
        settings = Settings.new
        @imdb = Imdb.new settings
      end

      describe '#get_data_for' do
        it 'gives correct movie data' do
          VCR.use_cassette('imdb_get_data_for_movie') do
            movie = OpenStruct.new(imdb_id: '0403358')
            expect(subject.get_data_for(movie).to_h
                  ).to include(
                    type: :movie,
                    title: 'Nochnoy dozor',
                    year: 2004
                  )
          end
        end

        it 'gets episode data correctly' do
          VCR.use_cassette('imdb_get_data_for_episode') do
            movie = OpenStruct.new(imdb_id: '2816136')
            expect(subject.get_data_for(movie).to_h
                  ).to include(
                    type: :episode,
                    episode_title: 'Two Swords',
                    year: 2014,
                    series_title: 'Game of Thrones',
                    season_number: 4,
                    episode_number: 1
                  )
          end
        end

        it 'does nothing unidentified movie' do
          VCR.use_cassette('imdb_get_data_for_wrong') do
            movie = OpenStruct.new(imdb_id: '9999999')
            unchanged = movie.dup
            expect(subject.get_data_for(movie)).to eq(unchanged)
          end
        end

        it 'doesn\'t mutate the movie object' do
          VCR.use_cassette('imdb_get_data_for_episode') do
            movie = OpenStruct.new(imdb_id: '2816136')
            old_movie = movie.dup
            subject.get_data_for(movie)
            expect(movie).to eq(old_movie)
          end
        end
      end

      describe '#get_data' do
        let(:ids) { %w(0403358 2816136) }

        it 'gives the same result as get_data_for for multiple movies' do
          VCR.use_cassette('imdb_get_data') do
            movies = ids.map { |id| OpenStruct.new(imdb_id: id) }
            expect(subject.get_data(movies)).to eq(movies.map do |movie|
              subject.get_data_for(movie)
            end)
          end
        end

        it 'doesn\'t mutate the movies object' do
          VCR.use_cassette('imdb_get_data') do
            movies = ids.map { |id| OpenStruct.new(imdb_id: id) }
            old_movies = movies.dup
            subject.get_data(movies)
            expect(movies).to eq(old_movies)
          end
        end
      end

      describe '.id' do
        subject { Imdb }

        it 'returns the id from an imdb url' do
          expect(
            subject.id('http://www.imdb.com/title/tt1063669')
          ).to eq('1063669')
        end

        it 'returns the id from a long imdb url' do
          expect(
            subject.id('http://www.imdb.com/title/tt1617661/?ref_=ttmd_md_nm')
          ).to eq('1617661')
        end

        it 'returns the id from.. the id' do
          expect(subject.id('1617661')).to eq('1617661')
        end

        it 'returns nil for wrong input' do
          expect(subject.id('foo bar baz')).to be_nil
        end
      end

      describe '.url' do
        subject { Imdb }

        it 'returns the url from an imdb id' do
          expect(
            subject.url('1063669')
          ).to eq('http://www.imdb.com/title/tt1063669')
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
