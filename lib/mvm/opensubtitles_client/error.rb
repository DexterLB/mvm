module Mvm
  class OpensubtitlesClient
    class Error < RuntimeError; end

    class NoStatusError             < Error; end
    class UnknownError              < Error; end

    class PartialContentError       < Error; end
    class MovedError                < Error; end
    class UnauthorizedError         < Error; end
    class InvalidFormatError        < Error; end
    class HashMismatchError         < Error; end
    class InvalidLanguageError      < Error; end
    class NotEnoughArgumentsError   < Error; end
    class NoSessionError            < Error; end
    class DownloadLimitError        < Error; end
    class InvalidArgumentsError     < Error; end
    class InvalidMethodError        < Error; end
    class UserAgentError            < Error; end
    class InvalidStringError        < Error; end
    class InvalidImdbidError        < Error; end
    class SubtitleError             < Error; end
    class ServiceUnavailableError   < Error; end

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
      410 => Error,
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
