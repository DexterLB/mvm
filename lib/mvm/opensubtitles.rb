require 'xmlrpc/client'
require 'iso-639'

module Mvm
  class Opensubtitles
    DEFAULT_CLIENT_OPTIONS = {
      host: 'http://api.opensubtitles.org',
      path: '/xml-rpc',
      timeout: 20,
    }

    def initialize(username: '',
                   password: '',
                   language: ISO_639.find('en'),
                   useragent:,
                   **client_options)
      @username = username
      @password = password
      @language = language
      @useragent = useragent

      client_options = DEFAULT_CLIENT_OPTIONS.merge(client_options)
      @server = XMLRPC::Client.new(**client_options)
    end

  end
end

# vim: set shiftwidth=2:
