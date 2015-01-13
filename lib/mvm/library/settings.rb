require 'ostruct'
require 'yaml'

module Mvm
  class Library
    DEFAULT_SETTINGS = {}

    class Settings < OpenStruct
      def initialize
        super DEFAULT_SETTINGS
      end

      def to_string_hash
        to_h.map { |setting, value| [setting.to_s, value] }.to_h
      end

      def to_yaml
        YAML.dump(to_string_hash).gsub(/^---\n/, '')
      end

      def from_yaml(yaml)
        YAML.load(yaml).map { |setting, value| self[setting] = value }
      end

      def clear
        each_pair { |key, _| delete_field key }
        DEFAULT_SETTINGS.each { |key, value| self[key] = value }
      end
    end
  end
end
# vim: set shiftwidth=2:
