project(
  'frontend',
  'c',
  default_options: [
    'werror=false',
    'warning_level=3',
    'buildtype=debug'
  ]
)

add_global_arguments('-g', language : 'c')

src = [
  'main.c',
  'qrcodegen.c',
  'external/QRCode/src/qrcode.c'
]

deps = [
  dependency('gtk+-3.0'),
  dependency('libcurl'),
  dependency('threads'),
]

incl = [
  include_directories('external/QRCode/src')
]

executable(
  'frontend',
  sources : src,
  dependencies : deps,
  include_directories : incl,
)

