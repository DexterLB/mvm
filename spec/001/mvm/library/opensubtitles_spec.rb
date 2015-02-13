require 'spec_helper'
require 'ostruct'
require 'mvm/library/opensubtitles'
require 'mvm/settings'

module Mvm
  class Library
    describe Opensubtitles do
      subject { @opensubtitles }

      before :each do
        @settings = Settings.new
        @opensubtitles = Opensubtitles.new settings: @settings
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

      describe '#search_subtitles_for' do
        it 'finds subtitles for given hash' do
          VCR.use_cassette('os_search_subtitles_for_hash') do
            movie = OpenStruct.new(file_hash: 'cabaa1d4f966107b')
            movie = subject.search_subtitles_for(movie)
            expect(movie.subtitles).not_to be_empty
            movie.subtitles.each do |subtitle|
              expect(
                subtitle.raw_info['MovieName']
              ).to eq('"Arrow" Home Invasion')
            end
          end
        end

        it 'gives url-ish download url' do
          VCR.use_cassette('os_search_subtitles_for_hash') do
            movie = OpenStruct.new(file_hash: 'cabaa1d4f966107b')
            movie = subject.search_subtitles_for(movie)
            expect(movie.subtitles).not_to be_empty
            movie.subtitles.each do |subtitle|
              expect(subtitle.url).to match(%r{^(http|ftp)://.*\.gz$})
            end
          end
        end

        it 'specifies an encoding' do
          VCR.use_cassette('os_search_subtitles_for_hash') do
            movie = OpenStruct.new(file_hash: 'cabaa1d4f966107b')
            movie = subject.search_subtitles_for(movie)
            expect(movie.subtitles).not_to be_empty
            movie.subtitles.each do |subtitle|
              expect(subtitle.encoding).to be_instance_of(Encoding)
            end
          end
        end

        it 'finds subtitles for given imdb id' do
          VCR.use_cassette('os_search_subtitles_for_imdb_id') do
            movie = OpenStruct.new(imdb_id: '2761432')
            movie = subject.search_subtitles_for(movie)
            expect(movie.subtitles).not_to be_empty
            movie.subtitles.each do |subtitle|
              expect(
                subtitle.raw_info['MovieName']
              ).to eq('"Arrow" Home Invasion')
            end
          end
        end

        it 'finds subtitles for given title' do
          VCR.use_cassette('os_search_subtitles_for_title') do
            movie = OpenStruct.new(title: 'Arrow Home Invasion')
            movie = subject.search_subtitles_for(movie)
            expect(movie.subtitles).not_to be_empty
            movie.subtitles.each do |subtitle|
              expect(
                subtitle.raw_info['MovieName']
              ).to eq('"Arrow" Home Invasion')
            end
          end
        end

        it 'finds subtitles with specified languages' do
          VCR.use_cassette('os_search_subtitles_for_imdb_id_lang') do
            @settings.subtitle_languages = 'ru,de'
            movie = OpenStruct.new(imdb_id: '2761432')
            movie = subject.search_subtitles_for(movie)
            expect(movie.subtitles.select do |subtitle|
              subtitle.language.english_name == 'Russian'
            end).not_to be_empty
            expect(movie.subtitles.select do |subtitle|
              subtitle.language.english_name == 'German'
            end).not_to be_empty
          end
        end

        it 'doesn\'t mutate the movie object' do
          VCR.use_cassette('os_search_subtitles_for_hash') do
            movie = OpenStruct.new(file_hash: 'cabaa1d4f966107b')
            old_movie = movie.dup
            subject.search_subtitles_for(movie)
            expect(old_movie).to eq(movie)
          end
        end
      end

      describe '#search_subtitles' do
        it 'calls #search_subtitles_for for all elements of the array' do
          expect(subject).to receive(:search_subtitles_for).with(:a_movie)
          expect(subject).to receive(:search_subtitles_for).with(:another_movie)
          expect(subject).to receive(:search_subtitles_for).with(:a_third_movie)
          subject.search_subtitles([:a_movie, :another_movie, :a_third_movie])
        end

        it 'returns an array of movies with subtitles' do
          movies = [
            OpenStruct.new(file_hash: '09a2c497663259cb'),
            OpenStruct.new(file_hash: '46e33be00464c12e')
          ]
          VCR.use_cassette('os_search_subtitles') do
            subject.search_subtitles(movies).each do |subtitle|
              expect(subtitle).to be_instance_of(OpenStruct)
            end
          end
        end

        it 'doesn\'t mutate the movies object' do
          movies = [
            OpenStruct.new(file_hash: '09a2c497663259cb'),
            OpenStruct.new(file_hash: '46e33be00464c12e')
          ]
          old_movies = movies.dup
          VCR.use_cassette('os_search_subtitles') do
            subject.search_subtitles(movies)
            expect(old_movies).to eq(movies)
          end
        end

        it 'yields consecutive progress reports' do
          movies = [
            OpenStruct.new(file_hash: '09a2c497663259cb'),
            OpenStruct.new(file_hash: '46e33be00464c12e')
          ]
          progress_reports = []
          VCR.use_cassette('os_search_subtitles') do
            subject.search_subtitles(movies) { |p| progress_reports << p }
            expect(progress_reports).to eq([
              [0, 2],
              [1, 2],
              [2, 2]
            ])
          end
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
