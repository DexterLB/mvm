require 'spec_helper'
require 'mvm/configuration'

require 'ostruct'
require 'yaml'

module Mvm
  describe Configuration do
    class Dummy
      include Configuration
    end

    subject { @dummy }

    before :each do
      @dummy = Dummy.new
    end

    describe '#settings' do
      it 'is an instance of Settings' do
        expect(subject.settings).to be_a Configuration::Settings
      end

      it 'stores configuration' do
        subject.settings.foo = 42
        expect(subject.settings.foo).to eq(42)
      end

      it 'uses Defaults' do
        stub_const('Mvm::Configuration::Defaults', {
          bar: 'baz'
        })
        expect(subject.settings.bar).to eq('baz')
      end

      it 'overrides Defaults' do
        stub_const('Mvm::Configuration::Defaults', {
          bar: 'baz'
        })
        subject.settings.bar = 'qux'
        expect(subject.settings.bar).to eq('qux')
      end

      it 'generates correct YAML' do
        subject.settings.foo = 5
        subject.settings.bar = 3.14
        subject.settings.baz = 'qux'
        expect(YAML.load(subject.settings.to_yaml)).to include(
          'foo' => 5, 'bar' => 3.14, 'baz' => 'qux'
        )
      end

      it 'reads YAML correctly' do
        subject.settings.from_yaml({
          'foo' => 5, 'bar' => 3.14, 'baz' => 'qux'
        }.to_yaml)
        expect(subject.settings.foo).to eq(5)
        expect(subject.settings.bar).to eq(3.14)
        expect(subject.settings.baz).to eq('qux')
      end

      it 'keeps old settings when reading YAML' do
        subject.settings.foo = 42
        subject.settings.from_yaml({bar: 3.14}.to_yaml)
        expect(subject.settings.foo).to eq(42)
        expect(subject.settings.bar).to eq(3.14)
      end

      it 'overrides old settings with new ones when reading YAML' do
        subject.settings.foo = 42
        subject.settings.from_yaml({'foo' => 3.14}.to_yaml)
        expect(subject.settings.foo).to eq(3.14)
      end
    end

    describe '#clear_settings' do
      it 'clears settings back to Defaults' do
        stub_const('Mvm::Configuration::Defaults', {
          'bar' => 'baz'
        })

        subject.settings.foo = 5
        subject.settings.bar = 3.14

        subject.clear_settings

        subject.settings.bar = 'qux'
        expect(subject.settings.foo).to be_nil
        expect(subject.settings.bar).to eq('qux')
      end
    end
  end
end

# vim: set shiftwidth=2:
