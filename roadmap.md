mvm
===
mvm ще е инструмент за идентифициране и сортиране на филми и сериали.
Ще използва данни от [IMDB](http://imdb.com) и ще съхранява информацията в прости текстови файлове.

Функции
-------

- [c] Команден интерфейс:
  `mvm import foo.mp4 --moviepath="movies/%name% (%year%).%ext%"`
- [f] Конфигурационен файл - зареждане на настройките от файл вместо от командата
- [c] Разпознаване на файл по хеш от [opensubtitles](http://opensubtitles.org)
- [c] Намиране на [IMDB](http://imdb.com) id
    - [c] по хеша от [opensubtitles](http://opensubtitles.org)
    - [f] чрез търсене по метаданните от файла
    - [f] чрез търсене на думите от името на файла
    - [c] чрез интерактивен избор на id от потребителя
    - [f] може би комбинация от предните три, и предлагане на възможности за избор
- [c] Взимане на информация за филм от [IMDB](http://imdb.com): име на филм,
  година, сезон/епизод (за сериали)..
- [c] Преименуване на файл спрямо зададено форматиране, например:
    - `"big.buck.bunny.webdl_1080p_x264.mkv"`   
       -> `"movies/Big Buck Bunny (2008)/Big Buck Bunny.mkv"`
    - `"BigFmiTheory.TheRubyMutation[MitioThePirate].avi"`  
       -> `"series/Big FMI Theory/season 2/08 - The Ruby Mutation.avi"`
- [f] Намиране на субтитри от [opensubtitles](http://opensubtitles.org): автоматично
  сваляне на първите n най-близки субтитри на езици x, y, z
- [f] Записване на данните за филм (opensubtitles id, IMDB id, име, година, сезон,
  епизод) в текстов файл
- [f] обработване на много филми наведнъж, намиране на всички филми в папка
- [m] извеждане на списък с вече обработените филми, търсене на филм
- [m] създаване на плейлисти по сериал/сезон
- [m] директно пускане на плеър върху намерен филм: `mvm play "two towers"`
- [m] директно пускане на плейлисти от епизоди: `mvm play "fmi theory" s2e3-5`

Нещата, маркирани с [c] са основни, тези с [f] са за ако ми остане време,
а тези с [m] са далечни мечти.

Ruby stuff
----------
Разни библиотеки, които може да се ползват в проекта:

- [XMLRPC::Client](http://www.ruby-doc.org/stdlib-2.1.5/libdoc/xmlrpc/rdoc/XMLRPC/Client.html):
  за ползване на функциите на [opensubtitles API](http://trac.opensubtitles.org/projects/opensubtitles/wiki/XMLRPC#CheckMovieHash2)
- [imdb](https://github.com/ariejan/imdb): за намиране на данни за вече
  разпознат филм
- [metadata](https://github.com/kig/metadata): за изваждане на метаданните на
  видео файл

Подобни неща
------------

- [filebot](http://www.filebot.net/): Прави същото, но има много повече функции
  и е много по-тежко. И е написано на java.  Много по-добро от това,
  което се опитвам да направя.
- [lm](https://github.com/RedRise/lm/blob/master/lm.py): Прави само разпознаване
  на филми и сваляне на субтитри. Не работи за сериали.
- [beets](http://beets.radbox.org/): Същото нещо, но за музика.
