#!/usr/bin/env python
#
# Copyright 2010 Hunter Freyer and Michael Kelly
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
from google.appengine.ext import webapp
from google.appengine.ext.webapp import util

import html_handlers
import models

def main():
  application = webapp.WSGIApplication([
      ("/", html_handlers.MainPageHandler),
      ("/about", html_handlers.AboutHandler),
      ("/make", html_handlers.MakeHandler),
      ("/make.do", html_handlers.MakeSubmitHandler),
      ("/mymfks", html_handlers.MyMfksHandler),
      ("/vote/(.*)", html_handlers.VoteHandler),
      ("/vote.do", html_handlers.VoteSubmitHandler),
      ("/i/(.*)", html_handlers.EntityImageHandler),
      ("/search", html_handlers.ImageSearchHandler),
      ("/.*", html_handlers.CatchAllHandler),
    ])
  util.run_wsgi_app(application)

if __name__ == '__main__':
  main()
