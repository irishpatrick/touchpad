project(
  'driver',
  'c',
  default_options: [
    'werror=true',
    'warning_level=3',
    'buildtype=debugoptimized'
  ]
)

src = [
  'driver.c'
]

deps = [
  dependency('libevdev')
]

add_project_arguments(['-Werror'], language : 'c')

library(
  'driver',
  sources : src,
  dependencies : deps
)
