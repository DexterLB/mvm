module Mvm
  class OpensubtitlesClient
    class OpensubtitlesError < RuntimeError; end

    class NoStatusError             < OpensubtitlesError; end
    class UnknownError              < OpensubtitlesError; end

    class PartialContentError       < OpensubtitlesError; end
    class MovedError                < OpensubtitlesError; end
    class UnauthorizedError         < OpensubtitlesError; end
    class InvalidFormatError        < OpensubtitlesError; end
    class HashMismatchError         < OpensubtitlesError; end
    class InvalidLanguageError      < OpensubtitlesError; end
    class NotEnoughArgumentsError   < OpensubtitlesError; end
    class NoSessionError            < OpensubtitlesError; end
    class DownloadLimitError        < OpensubtitlesError; end
    class InvalidArgumentsError     < OpensubtitlesError; end
    class InvalidMethodError        < OpensubtitlesError; end
    class UserAgentError            < OpensubtitlesError; end
    class InvalidStringError        < OpensubtitlesError; end
    class InvalidImdbidError        < OpensubtitlesError; end
    class SubtitleError             < OpensubtitlesError; end
    class ServiceUnavailableError   < OpensubtitlesError; end

    ERRORS = {
      206 => PartialContentError,
      301 => MovedError,
      401 => UnauthorizedError,
      402 => InvalidFormatError,
      403 => HashMismatchError,
      404 => InvalidLanguageError,
      405 => NotEnoughArgumentsError,
      406 => NoSessionError,
      407 => DownloadLimitError,
      408 => InvalidArgumentsError,
      409 => InvalidMethodError,
      410 => OpensubtitlesError,
      411 => UserAgentError,
      412 => InvalidStringError,
      413 => InvalidImdbidError,
      414 => UserAgentError,
      415 => UserAgentError,
      416 => SubtitleError,
      503 => ServiceUnavailableError
    }
  end
end

# vim: set shiftwidth=2:
