require 'spec_helper'
require 'ostruct'
require 'fakefs/safe'
require 'mvm/library/store'
require 'mvm/settings'

module Mvm
  class Library
    describe Store do
      subject { @store }

      after :each do
        FakeFS::FileSystem.clear
      end

      describe '#store_movie' do
        it 'stores movie to file given in settings' do
          @store = Store.new(Settings.new(store_file: 'foo'))
          FakeFS do
            movie = OpenStruct.new(filename: 'bar', baz: 'qux')
            subject.store_movie(movie)
            expect(YAML.load(File.read('foo'))).to eq(movie.to_h)
          end
        end

        it 'stores movie to filename.mvm by default' do
          @store = Store.new(Settings.new)
          FakeFS do
            movie = OpenStruct.new(filename: 'bar', baz: 'qux')
            subject.store_movie(movie)
            expect(YAML.load(File.read('bar.mvm'))).to eq(movie.to_h)
          end
        end

        it 'allows formatting with movie values in filename' do
          @store = Store.new(
            Settings.new(store_file: 'a_file_with_%{baz}_inside')
          )
          FakeFS do
            movie = OpenStruct.new(filename: 'bar', baz: 'qux')
            subject.store_movie(movie)
            expect(
              YAML.load(File.read('a_file_with_qux_inside'))
            ).to eq(movie.to_h)
          end
        end
      end

      describe '#load_movie' do
        before :each do
          @store = Store.new(Settings.new)
        end

        it 'reads data from yaml file' do
          movie = OpenStruct.new(filename: 'bar', baz: 'qux')
          FakeFS do
            File.write('foo', movie.to_h.to_yaml)
            expect(subject.load_movie('foo')).to eq(movie)
          end
        end
      end

      describe '#store_movies' do
        it 'stores multiple movies to filename.mvm by default' do
          @store = Store.new(Settings.new)
          FakeFS do
            movies = [
              OpenStruct.new(filename: 'foo', bar: 'baz'),
              OpenStruct.new(filename: 'qux', omg: 42)
            ]
            subject.store_movies(movies)
            expect(YAML.load(File.read('foo.mvm'))).to eq(movies[0].to_h)
            expect(YAML.load(File.read('qux.mvm'))).to eq(movies[1].to_h)
          end
        end
      end

      describe '#load_movies' do
        it 'loads multiple movies from folder with .mvm\'s by default' do
          @store = Store.new(Settings.new)
          movies = [
            OpenStruct.new(filename: 'foo', bar: 'baz'),
            OpenStruct.new(filename: 'qux', omg: 42)
          ]
          FakeFS do
            Dir.mkdir('folder0')
            movies.each do |movie|
              File.write("folder0/#{movie.filename}.mvm", movie.to_h.to_yaml)
            end
            expect(subject.load_movies('folder0')).to match_array(movies)
          end
        end

        it 'can match another expression' do
          @store = Store.new(Settings.new(store_match: '\-42$'))
          movies = [
            OpenStruct.new(filename: 'foo', bar: 'baz'),
            OpenStruct.new(filename: 'qux', omg: 42)
          ]
          FakeFS do
            Dir.mkdir('folder1')
            movies.each do |movie|
              File.write("folder1/#{movie.filename}-42", movie.to_h.to_yaml)
            end
            expect(subject.load_movies('folder1')).to match_array(movies)
          end
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
