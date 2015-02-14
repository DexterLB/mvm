require 'mvm/opensubtitles_client'
require 'mvm/settings'

require 'iso-639'

module Mvm
  class Library
    DEFAULT_SETTINGS.merge!(
      opensubtitles_username: '',
      opensubtitles_password: '',
      opensubtitles_useragent: 'OSTestUserAgent',
      opensubtitles_language: 'en',
      opensubtitles_timeout: 20,
      subtitle_languages: 'en,bg',
      max_subtitles: 5,
      subtitle_filename: '%{filename}.%{subtitle_index}.srt'
    )

    class Opensubtitles
      def initialize(settings: Settings.new)
        @settings = settings
      end

      def id_by_hash_for(movie)
        set_attributes_for(movie, client.lookup_hash(movie.file_hash))
      end

      def id_by_hashes(movies)
        data = client.lookup_hashes(movies.select(&:file_hash).map(&:file_hash))
        movies.map do |movie|
          set_attributes_for(movie, data[movie.file_hash])
        end
      end

      def search_subtitles(movies)
        movies_with_subtitles = movies.each.with_index.map do |movie, current|
          yield [current, movies.size] if block_given?
          search_subtitles_for(movie)
        end

        yield [movies.size, movies.size] if block_given?
        movies_with_subtitles
      end

      def search_subtitles_for(movie)
        movie = movie.dup

        languages = parse_languages(@settings.subtitle_languages)
        subtitles_data = languages.map do |language|
          results = client.search_subtitles([subtitle_query(movie, language)])
          sort_subtitles(results).take(@settings.max_subtitles)
        end.flatten

        movie.subtitles = subtitles_data.map do |data|
          subtitle_attributes(data)
        end

        movie
      end

      def download_subtitles_for(movie)
        movie = movie.dup
        movie.subtitles = movie.subtitles.dup

        movie.subtitles.each_with_index do |subtitle, index|
          subtitle.filename = format(@settings.subtitle_filename,
                                     movie.to_h.merge(subtitle_index: index))
          @client.download_gz(subtitle.url, subtitle.filename)
        end
      end

      def download_subtitles(movies)
        movies_with_sub_files = movies.each.with_index.map do |movie, current|
          yield [current, movies.size] if block_given?
          download_subtitles_for(movie)
        end

        yield [movies.size, movies.size] if block_given?
        movies_with_sub_files
      end

      private

      def subtitle_query(movie, languages)
        {
          sublanguageid: languages,
          moviehash: movie.file_hash,
          moviebytesize: movie.filesize,
          imdbid: movie.imdb_id,
          query: movie.title,
          season: movie.season_number,
          episode: movie.episode_number
        }.select { |_, value| value }.map do |key, value|   # haters gonna hate
          [key.to_s, value]
        end.to_h
      end

      def subtitle_attributes(data)
        begin
          encoding = Encoding.find(data['SubEncoding'])
        rescue ArgumentError
          encoding = nil
        end

        OpenStruct.new(
          raw_info: data,
          language: ISO_639.find_by_code(data['SubLanguageID']),
          release: data['MovieReleaseName'],
          framerate: data['MovieFPS'].to_f,
          rating: data['SubRating'].to_f,
          downloads: data['SubDownloadsCnt'].to_i,
          encoding: encoding,
          url: data['SubDownloadLink']
        )
      end

      def set_attributes_for(movie, attributes)
        movie = movie.dup

        return movie unless attributes
        return movie unless %w(episode movie).include? attributes['MovieKind']

        movie.title = attributes['MovieName']
        movie.type = attributes['MovieKind'].to_sym
        movie.year = attributes['MovieYear'].to_i
        movie.imdb_id = attributes['MovieImdbID']

        if movie.type == :episode
          movie.series_title, movie.episode_title = movie.title
                                                    .match(/\"(.+)\"\s(.+)/)
                                                    .captures

          movie.season_number = attributes['SeriesSeason'].to_i
          movie.episode_number = attributes['SeriesEpisode'].to_i
        end

        movie
      end

      def sort_subtitles(results)
        results.sort_by { |result| sort_rating(result) }
      end

      def sort_rating(result)
        [
          case result['MatchedBy']
          when 'moviehash' then 0
          when 'imdbid'    then 1
          else 2
          end,
          -result['SubDownloadsCnt'].to_i,
          -result['SubRating'].to_f
        ]
      end

      def client
        @client ||= OpensubtitlesClient.new(
          username: @settings.opensubtitles_username,
          password: @settings.opensubtitles_password,
          useragent: @settings.opensubtitles_useragent,
          language: ISO_639.find(@settings.opensubtitles_language),
          timeout: @settings.opensubtitles_timeout
        )
      end

      def parse_languages(comma_separated_languages)
        comma_separated_languages.split(',').map do |code|
          language = ISO_639.find(code)
          fail 'Unknown language: ' + code unless language
          language.alpha3
        end
      end
    end
  end
end

# vim: set shiftwidth=2:
