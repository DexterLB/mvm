require 'spec_helper'
require 'fakefs/safe'
require 'fileutils'
require 'mvm/library/files'

module Mvm
  class Library
    describe Files do
      subject { @files }

      before :each do
        @settings = Settings.new
        @files = Files.new(@settings)
      end

      let(:sample_video) do
        File.dirname(__FILE__) + '/../../../fixtures/drop.avi'
      end

      let(:movie) do
        OpenStruct.new(filename: sample_video)
      end

      describe '#calculate_hash_for' do
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

      describe '#calculate_hashes' do
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

      describe '#load' do
        after :each do
          FakeFS::FileSystem.clear
        end

        it 'recognizes a single file' do
          FakeFS do
            FileUtils.touch('foo.mkv')
            movies = subject.load('foo.mkv')
            expect(movies.size).to eq(1)
            expect(movies[0].filename).to eq('foo.mkv')
          end
        end

        it 'recursively adds mkv movies from folders' do
          FakeFS do
            FileUtils.touch('foo.mkv')
            FileUtils.touch('bar.mkv')
            Dir.mkdir('baz')
            FileUtils.touch('baz/qux.mkv')
            FileUtils.touch('not_a_movie.txt')

            movies = subject.load('.')
            expect(movies.map(&:filename)).to match_array(
              ['./foo.mkv', './bar.mkv', './baz/qux.mkv']
            )
          end
        end

        it 'fails for non-existant file' do
          FakeFS do
            expect { subject.load('foo.mkv') }.to raise_error
          end
        end

        it 'creates an "added" attribute which is time' do
          FakeFS do
            FileUtils.touch('foo.mkv')
            expect(subject.load('foo.mkv')[0].added).to be_instance_of(Time)
            # expect(subject.load('foo.mkv').added).to eq(File.mtime('foo.mkv'))
            # bug in FakeFS
          end
        end

        it 'creates an "extension" attribute which is the file extension' do
          FakeFS do
            FileUtils.touch('foo.mkv')
            expect(subject.load('foo.mkv')[0].extension).to eq('.mkv')
          end
        end

        it 'creates a "filesize" attribute which is the file size' do
          FakeFS do
            File.write('foo.mkv', 'Ruby Ruby Ruby Ruby Ruby')
            expect(
              subject.load('foo.mkv')[0].filesize
            ).to eq('Ruby Ruby Ruby Ruby Ruby'.size)
          end
        end

        it 'uses settings#valid_movie_extensions' do
          FakeFS do
            FileUtils.touch('file1.foo')
            FileUtils.touch('file2.foo')
            FileUtils.touch('file3.bar')
            FileUtils.touch('file4.mkv')

            @settings.valid_movie_extensions = '.foo .bar'
            movies = subject.load('.')

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
