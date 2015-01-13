require 'streamio-ffmpeg'

module Mvm
  class Library
    class Metadata
      def self.read_metadata_for(movie)
        movie = movie.dup

        ffmpeg = FFMPEG::Movie.new(movie.filename)
        return unless ffmpeg.valid?

        [:duration, :bitrate, :size, :video_codec, :width, :height,
         :frame_rate, :audio_codec, :audio_sample_rate, :audio_channels
        ].each do |attribute|
          movie[attribute] = ffmpeg.send(attribute)
        end

        movie
      end

      def self.read_metadata(movies)
        movies.map { |movie| read_metadata_for(movie) }
      end
    end
  end
end

# vim: set shiftwidth=2:
