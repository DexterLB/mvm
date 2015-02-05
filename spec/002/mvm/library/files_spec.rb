require 'spec_helper'
require 'fakefs/safe'
require 'fileutils'
require 'mvm/library/files'

module Mvm
  class Library
    describe Files do
      subject { Files }

      let(:sample_video) do
        File.dirname(__FILE__) + '/../../../fixtures/drop.avi'
      end

      let(:movie) do
        OpenStruct.new(filename: sample_video)
      end

      describe '.calculate_hash_for' do
        it 'calculates correct hash' do
          expect(subject.calculate_hash_for(movie).file_hash)
            .to eq('450f3f0c98a1f11d')
        end

        it 'doesn\'t mutate the movie object' do
          old_movie = movie.dup
          subject.calculate_hash_for(movie)
          expect(movie).to eq(old_movie)
        end
      end

      describe '.calculate_hashes' do
        let(:movies) { [movie] }
        it 'calculates correct hash for list of one movie' do
          result = subject.calculate_hashes(movies)
          expect(result.first.file_hash).to eq('450f3f0c98a1f11d')
          expect(result.size).to eq(1)
        end

        it 'doesn\'t mutate the movies object' do
          old_movies = movies.dup
          subject.calculate_hashes(movies)
          expect(movies).to eq(old_movies)
        end
      end

      describe '.scan_folder' do
        it 'recursively adds mkv movies from folders' do
          FakeFS do
            FileUtils.touch('foo.mkv')
            FileUtils.touch('bar.mkv')
            Dir.mkdir('baz')
            FileUtils.touch('baz/qux.mkv')
            FileUtils.touch('not_a_movie.txt')

            movies = subject.new.scan_folder('.')
            expect(movies.map(&:filename)).to match_array(
              ['./foo.mkv', './bar.mkv', './baz/qux.mkv']
            )
          end
        end

        it 'creates and "added" attribute which is time' do
          FakeFS do
            FileUtils.touch('foo.mkv')
            movies = subject.new.scan_folder('.')
            expect(movies[0].added).to be_instance_of(Time)
            # expect(movies[0].added).to eq(File.mtime('foo.mkv'))
            # bug in FakeFS
          end
        end

        it 'creates and "extension" attribute which is the file extension' do
          FakeFS do
            FileUtils.touch('foo.mkv')
            movies = subject.new.scan_folder('.')
            expect(movies[0].extension).to eq('.mkv')
            # expect(movies[0].added).to eq(File.mtime('foo.mkv'))
            # bug in FakeFS
          end
        end

        it 'uses settings#valid_movie_extensions' do
          FakeFS do
            FileUtils.touch('file1.foo')
            FileUtils.touch('file2.foo')
            FileUtils.touch('file3.bar')
            FileUtils.touch('file4.mkv')

            movies = subject.new(Settings.new(
              valid_movie_extensions: '.foo .bar'
            )).scan_folder('.')

            expect(movies.map(&:filename)).to match_array(
              ['./file1.foo', './file2.foo', './file3.bar']
            )
          end
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
