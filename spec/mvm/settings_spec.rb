require 'spec_helper'
require 'mvm/settings'

require 'ostruct'
require 'yaml'

module Mvm
  describe Settings do
    subject { @settings }

    before :each do
      @settings = Settings.new
    end

    it 'stores configuration' do
      subject.foo = 42
      expect(subject.foo).to eq(42)
    end

    it 'uses DEFAULT_SETTINGS' do
      stub_const('Mvm::DEFAULT_SETTINGS', bar: 'baz')
      @settings = Settings.new
      expect(subject.bar).to eq('baz')
    end

    it 'overrides DEFAULT_SETTINGS' do
      stub_const('Mvm::DEFAULT_SETTINGS', bar: 'baz')
      subject.bar = 'qux'
      expect(subject.bar).to eq('qux')
    end

    it 'reads settings from hash given to constructor' do
      @settings = Settings.new(foo: 'bar')
      subject.bar = 'qux'
      expect(subject.foo).to eq('bar')
      expect(subject.bar).to eq('qux')
    end

    it 'overrides DEFAULT_SETTINGS with hash given to constructor' do
      stub_const('Mvm::DEFAULT_SETTINGS', foo: 42)
      @settings = Settings.new(foo: 'bar')
      subject.bar = 'qux'
      expect(subject.foo).to eq('bar')
      expect(subject.bar).to eq('qux')
    end

    describe '#to_yaml' do
      it 'generates correct YAML' do
        subject.foo = 5
        subject.bar = 3.14
        subject.baz = 'qux'
        expect(YAML.load(subject.to_yaml)).to include(
          'foo' => 5, 'bar' => 3.14, 'baz' => 'qux'
        )
      end
    end

    describe '#from_yaml' do
      it 'reads YAML correctly' do
        subject.from_yaml({
          'foo' => 5, 'bar' => 3.14, 'baz' => 'qux'
        }.to_yaml)
        expect(subject.foo).to eq(5)
        expect(subject.bar).to eq(3.14)
        expect(subject.baz).to eq('qux')
      end

      it 'keeps old settings when reading YAML' do
        subject.foo = 42
        subject.from_yaml({ bar: 3.14 }.to_yaml)
        expect(subject.foo).to eq(42)
        expect(subject.bar).to eq(3.14)
      end

      it 'overrides old settings with new ones when reading YAML' do
        subject.foo = 42
        subject.from_yaml({ 'foo' => 3.14 }.to_yaml)
        expect(subject.foo).to eq(3.14)
      end
    end

    describe '#clear' do
      it 'clears settings back to DEFAULT_SETTINGS' do
        stub_const('Mvm::DEFAULT_SETTINGS', 'bar' => 'baz')
        @settings = Settings.new

        subject.foo = 5
        subject.bar = 3.14

        subject.clear

        expect(subject.foo).to be_nil
        expect(subject.bar).to eq('baz')
      end
    end
  end
end

# vim: set shiftwidth=2:
