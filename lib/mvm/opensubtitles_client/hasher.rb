module Mvm
  class OpensubtitlesClient
    class Hasher
      # opensubtitles.org hash algorithm:
      # sum of first 64kb + sum of last 64kb + filesize, truncated to 64bits

      CHUNK_SIZE = 64 * 1024 # in bytes

      def self.hash(filename)
        size = File.size(filename)
        return nil unless size >= CHUNK_SIZE

        hash = size

        File.open(filename, 'rb') do |f|
          hash += hash_fragment(f) # hash first 64kb
          f.seek([0, size - CHUNK_SIZE].max, IO::SEEK_SET)
          hash += hash_fragment(f) # hash last 64kb
        end
        sprintf('%016x', hash & 2**64 - 1)
      end

      def self.hash_fragment(file)
        file.read(CHUNK_SIZE).unpack('Q*').inject(:+)
      end
    end
  end
end

# vim: set shiftwidth=2:
