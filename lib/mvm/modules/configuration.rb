require 'ostruct'
require 'yaml'
require 'xdg'

module Mvm
  module Modules
    module Configuration
      Defaults = {}

      class Settings < OpenStruct
        def to_yaml
          YAML.dump(to_h).gsub(/^---\n/, '')
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
end
# vim: set shiftwidth=2:
