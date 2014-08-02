gico as golang cli
====

switch .gitconfig according to your environment

## Description

create and switch your git configuration file (a.k.a. .gitconfig).  
It helps you to switch some git configurations (proxy, authentication etc...) at home, company and more !

## Usage

At first, initialize your gico environment

    $ gico init

and create new git env

   $ gico create macbook
   # edit ~/.gitconfig.macbook to write git configuration for your macbook

switch to use it

    $ gico list
    * local
      macbook
    $ gico use macbook

enjoy !

## Install

download from [releases page](https://github.com/fumitoito/gico/releases)

## License

The MIT License (MIT)

Copyright (c) 2014 [fumitoito](https://github.com/fumitoito)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

## Author

[fumitoito](https://github.com/fumitoito)
