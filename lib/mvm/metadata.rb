require 'streamio-ffmpeg'

module Mvm
  module Metadata
    def read_metadata_for(movie)
      ffmpeg = FFMPEG::Movie.new(movie.filename)
      return unless ffmpeg.valid?

      [:duration, :bitrate, :size, :video_codec, :width, :height,
       :frame_rate, :audio_codec, :audio_sample_rate, :audio_channels
      ].each do |attribute|
        movie[attribute] = ffmpeg.send(attribute)
      end
    end

    def read_metadata
      movies.each { |movie| read_metadata_for(movie) }
    end
  end
end

# vim: set shiftwidth=2:
