#!/usr/bin/env python
#
# Copyright 2010 Hunter Freyer and Michael Kelly
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import ezt

from google.appengine.ext import webapp

class MainPageHandler(webapp.RequestHandler):
    def get(self):
        template = ezt.Template("templates/vote.html")
        template.generate(self.response.out, dict())

class AboutHandler(webapp.RequestHandler):
    def get(self):
        template = ezt.Template("templates/about.html")
        template.generate(self.response.out, dict())

class MakeHandler(webapp.RequestHandler):
    def get(self):
        template = ezt.Template("templates/make.html")
        template.generate(self.response.out, dict())

