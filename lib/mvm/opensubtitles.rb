require 'xmlrpc/client'


require 'iso-639'

module Mvm
  class Opensubtitles
    attr_reader :token

    DEFAULT_CLIENT_OPTIONS = {
      host: 'api.opensubtitles.org',
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
      @client = new_client(**client_options)
    end

    def info
      @client.call('ServerInfo')
    end

    def login
      @token = call('LogIn',
                    @username, 
                    @password, 
                    @language.alpha2, 
                    @useragent)['token']
    end

    def call(function, *arguments)
      @client.call(function, *arguments)
      # todo: check the status
    end

    private
    def new_client(**options)
      # maybe make this a monkey patch?
      option_sequence = options.values_at(:host, :path, :port, :proxy_host,
                                          :proxy_port, :http_user,
                                          :http_password, :use_ssl, :timeout)
      XMLRPC::Client.new(*option_sequence)
    end
  end
end

# vim: set shiftwidth=2:
