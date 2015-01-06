require 'ostruct'
require 'yaml'

module Mvm
  module Configuration
    Defaults = {}

    class Settings < OpenStruct
      def to_string_hash
        to_h.map { |setting, value| [setting.to_s, value] }.to_h
      end

      def to_yaml
        YAML.dump(to_string_hash).gsub(/^---\n/, '')
      end

      def from_yaml(yaml)
        YAML.load(yaml).map { |setting, value| self[setting] = value }
      end
    end

    def settings
      @settings || @settings = Settings.new(Defaults)
    end

    def clear_settings
      @settings = nil
    end
  end
end
# vim: set shiftwidth=2:
