require 'xmlrpc/client'
require 'iso-639'

require 'mvm/api/opensubtitles_error'

module Mvm
  module Api
    class OpensubtitlesClient
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
        @token = safe_client_call('LogIn',
                                  @username, 
                                  @password, 
                                  @language.alpha2, 
                                  @useragent)['token']
      end
      
      def logout
        call('LogOut') if @token
        @token = nil
      end

      def safe_client_call(function, *arguments)
        data = @client.call(function, *arguments)
        if not data['status']
          raise NoStatusError
        end
        code = data['status'].split.first.to_i
        unless code == 200
          error = ERRORS.fetch(code.to_i, UnknownError)
          raise error, data['status']
        end
        data
      end

      def call(function, *arguments)
        if not @token
          login
        end

        begin
          safe_client_call(function, @token, *arguments)
        rescue NoSessionError
          login
          safe_client_call(function, @token, *arguments)
        end
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
end

# vim: set shiftwidth=2:
