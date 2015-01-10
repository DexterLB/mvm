require 'streamio-ffmpeg'

module Mvm
  module Metadata
    def read_metadata_for(movie)
      ffmpeg = FFMPEG::Movie.new(movie.filename)
      return unless ffmpeg.valid?

      movie.length = ffmpeg.duration
      movie.bitrate = ffmpeg.bitrate
      movie.filesize = ffmpeg.size

      movie.video_codec = ffmpeg.video_codec
      movie.width = ffmpeg.width
      movie.height = ffmpeg.height
      movie.frame_rate = ffmpeg.frame_rate

      movie.audio_codec = ffmpeg.audio_codec
      movie.audio_sample_rate = ffmpeg.audio_sample_rate
      movie.audio_channels = ffmpeg.audio_channels

      movie
    end

    def read_metadata
      movies.each { |movie| read_metadata_for(movie) }
    end
  end
end

# vim: set shiftwidth=2:
