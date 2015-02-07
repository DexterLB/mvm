require 'spec_helper'
require 'fakefs/safe'
require 'fileutils'
require 'mvm/library/rename'

module Mvm
  class Library
    describe Renamer do
      subject { @renamer }

      before :each do
        @settings = Settings.new
        @renamer = Renamer.new @settings
      end

      describe '#rename_movie' do
        it 'only changes filename with strategy: dummy' do
          @settings.rename_strategy = 'dummy'
          @settings.movie_pattern = 'bar'
          movie = OpenStruct.new(filename: 'foo', type: :movie)
          new_movie = subject.rename_movie(movie)
          expect(new_movie.filename).to eq('bar')
        end

        it 'uses episode_pattern for episodes' do
          @settings.rename_strategy = 'dummy'
          @settings.movie_pattern = 'bar'
          @settings.episode_pattern = 'baz'
          movie = OpenStruct.new(filename: 'foo', type: :episode)
          new_movie = subject.rename_movie(movie)
          expect(new_movie.filename).to eq('baz')
        end

        it 'copies file with strategy: copy' do
          @settings.rename_strategy = 'copy'
          @settings.movie_pattern = 'bar'
          movie = OpenStruct.new(filename: 'foo', type: :movie)
          FakeFS do
            File.write('foo', 'qux')
            new_movie = subject.rename_movie(movie)
            expect(new_movie.filename).to eq('bar')
            expect(File.read('bar')).to eq('qux')
            expect(File.read('foo')).to eq('qux')
          end
        end

        it 'moves file with strategy: move' do
          @settings.rename_strategy = 'move'
          @settings.movie_pattern = 'bar'
          movie = OpenStruct.new(filename: 'foo', type: :movie)
          FakeFS do
            File.write('foo', 'qux')
            new_movie = subject.rename_movie(movie)
            expect(new_movie.filename).to eq('bar')
            expect(File.read('bar')).to eq('qux')
            expect(File.exist?('foo')).to be_falsey
          end
        end

        it 'creates symlink to file with strategy: symlink' do
          @settings.rename_strategy = 'symlink'
          @settings.movie_pattern = 'bar'
          movie = OpenStruct.new(filename: 'foo', type: :movie)
          FakeFS do
            File.write('foo', 'qux')
            new_movie = subject.rename_movie(movie)
            expect(new_movie.filename).to eq('bar')
            expect(File.read('bar')).to eq('qux')
            expect(File.readlink('bar')).to eq('foo')
          end
        end

        it 'moves file and replaces it with symlink with strategy: keeplink' do
          @settings.rename_strategy = 'keeplink'
          @settings.movie_pattern = 'bar'
          movie = OpenStruct.new(filename: 'foo', type: :movie)
          FakeFS do
            File.write('foo', 'qux')
            new_movie = subject.rename_movie(movie)
            expect(new_movie.filename).to eq('bar')
            expect(File.read('bar')).to eq('qux')
            expect(File.readlink('foo')).to eq('bar')
          end
        end

        it 'executes command with strategy: exec: ' do
          @settings.rename_strategy = 'exec: this is a command'
          @settings.movie_pattern = 'bar'
          movie = OpenStruct.new(filename: 'foo', type: :movie)
          expect(subject).to receive(:system).with('this is a command')
          new_movie = subject.rename_movie(movie)
          expect(new_movie.filename).to eq('bar')
        end

        it 'expands arguments of command with strategy: exec: ' do
          @settings.rename_strategy = 'exec: qux %{old} -> %{new}'
          @settings.movie_pattern = 'bar'
          movie = OpenStruct.new(filename: 'foo', type: :movie)
          expect(subject).to receive(:system).with('qux foo -> bar')
          new_movie = subject.rename_movie(movie)
          expect(new_movie.filename).to eq('bar')
        end

        it 'expands movie attributes inside target file' do
          @settings.rename_strategy = 'dummy'
          @settings.movie_pattern = 'bar %{title}'
          movie = OpenStruct.new(filename: 'foo', type: :movie, title: '42')
          new_movie = subject.rename_movie(movie)
          expect(new_movie.filename).to eq('bar 42')
        end

        it 'replaces nasty characters with underscores by default' do
          @settings.rename_strategy = 'dummy'
          @settings.movie_pattern = '%{title}'
          movie = OpenStruct.new(filename: 'foo',
                                 title: 'a|<>:b',
                                 type: :movie)
          new_movie = subject.rename_movie(movie)
          expect(new_movie.filename).to eq('a____b')
        end

        it 'takes fs_forbidden_char_exp & fs_forbidden_char_replace in mind' do
          @settings.rename_strategy = 'dummy'
          @settings.movie_pattern = '%{title}'
          @settings.fs_forbidden_char_exp = 'bar'
          @settings.fs_forbidden_char_replace = 'qux'
          movie = OpenStruct.new(filename: 'foo',
                                 title: 'foo bar baz',
                                 type: :movie)
          new_movie = subject.rename_movie(movie)
          expect(new_movie.filename).to eq('foo qux baz')
        end
      end

      describe '#rename_movies' do
        it 'calls #rename_movie for all elements of the array' do
          expect(subject).to receive(:rename_movie).with(:a_movie)
          expect(subject).to receive(:rename_movie).with(:another_movie)
          expect(subject).to receive(:rename_movie).with(:a_third_movie)
          subject.rename_movies([:a_movie, :another_movie, :a_third_movie])
        end

        it 'returns an array of movies with changed filenames' do
          @settings.rename_strategy = 'dummy'
          @settings.movie_pattern = '%{title}'
          movies = [
            OpenStruct.new(filename: 'foo', title: 't_foo', type: :movie),
            OpenStruct.new(filename: 'bar', title: 't_bar', type: :movie),
            OpenStruct.new(filename: 'baz', title: 't_baz', type: :movie)
          ]
          expect(subject.rename_movies(movies).map(&:filename)).to eq(
            %w(t_foo t_bar t_baz)
          )
        end

        it 'doesn\'t mutate the movies object' do
          @settings.rename_strategy = 'dummy'
          @settings.movie_pattern = 'qux'
          movies = [
            OpenStruct.new(filename: 'foo', type: :movie),
            OpenStruct.new(filename: 'bar', type: :movie),
            OpenStruct.new(filename: 'baz', type: :movie)
          ]
          old_movies = movies.dup
          subject.rename_movies(movies)
          expect(old_movies).to eq(movies)
        end

        it 'yields consecutive progress reports' do
          @settings.rename_strategy = 'dummy'
          @settings.movie_pattern = 'qux'
          movies = [
            OpenStruct.new(filename: 'foo', type: :movie),
            OpenStruct.new(filename: 'bar', type: :movie),
            OpenStruct.new(filename: 'baz', type: :movie)
          ]
          progress_reports = []
          subject.rename_movies(movies) { |p| progress_reports << p }
          expect(progress_reports).to eq([
            [0, 3],
            [1, 3],
            [2, 3],
            [3, 3]
          ])
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
